const { randomBytes } = require('crypto')
const secp256k1 = require('secp256k1')
var utils = require('ethers').utils;
var Election = artifacts.require("Election");

var sleepPeriod = 500; // millisecond.

async function sleep(msec) {
    return new Promise(resolve => setTimeout(resolve, msec));
}

async function waitUntilYouSeeBlock(b) {
    let currentBlock = await web3.eth.getBlockNumber();
    while (currentBlock < b) {
        await sleep(sleepPeriod);
        currentBlock = await web3.eth.getBlockNumber();
    }
}

// n could be String|Number|BN|BigNumber.
function toUint8ArrayL32(n) {
    return utils.padZeros(utils.arrayify(web3.utils.toHex(n)), 32);
}

// calculate Secp256k1 VRF proof given a seed and private key.
function calcProof(_seed, _privateKey) {
    _proof = [];
    for(let tmp = 0; tmp < 4; tmp++) {
        _proof[tmp] = randomBytes(32);
    }
    return _proof;
}

contract('Election0', (accounts) => {
    it("benchmark", async () => {
        let election = await Election.deployed();

        var initialSeed = 23;
        var difficultyParam = 5;
        var epochLength = 35;
        var challengePeriod = 10;

        // each validator generates their own keys and send them to contract.
        var honestValidatorNum = 18; // the first "honestValidatorNum" validators are honest. TODO
        var silentValidatorNum = 20 - honestValidatorNum; // the last "silentValidatorNum" validators are silent.
        var weights = [];
        var privateKeys = []; // privateKeys[i] = i + 1;
        var publicKeys = [];

        var totalTestEpochs = 200;
        var ctx = new Uint8Array(100); // only the size matters for us now, not the content.

        var wins = []; // number of epochs that validator i had the lowest mh,
        // if k validators win at the same block, each will get 1/k points.
        var totalWins = 0; // total number of times anyone had the lowest mh in any epoch.
        var totalSubmit = 0;
        var totalSubmitGas = 0;
        var averageSubmitGas = 0;

        for (let i = 0; i < honestValidatorNum + silentValidatorNum; i++) {
//            weights[i] = Math.floor(i/4 + 1); // to set TODO
            weights[i] = 1; // to set TODO
            privateKeys[i] = toUint8ArrayL32(i + 1);
            publicKeys[i] = secp256k1.publicKeyCreate(privateKeys[i], false);
            let pk = [];
            pk[0] = publicKeys[i].subarray(1, 33);
            pk[1] = publicKeys[i].subarray(33, 65);
            await election.addAsValidator(pk, weights[i], {from: accounts[i]});
            wins[i] = 0;
        }

        await election.startProtocol(initialSeed, difficultyParam, epochLength, challengePeriod);

        // testing if the protocol parameters are set as expected.
/*        let params = await election.getProtocolParams();
        assert.equal(params[0], 60, "total weight is not right");
        assert.equal(params[1], 5, "difficulty is not right"); */

        // everything until here surely works, among functions at top all work and calcProof should be written!

        for (let currentEpoch = 0; currentEpoch < totalTestEpochs; currentEpoch++) {
            console.log("________________");
            console.log("epoch", currentEpoch);

            let cS = await election.getCurrentSeed();
            let currentSeed = toUint8ArrayL32(cS);
            console.log("Seed ", currentSeed);

            let proofs = [];
            let outputs = [];
            for (let i = 0; i < honestValidatorNum; i++) {
                proofs[i] = calcProof(currentSeed, privateKeys[i]);
                outputs[i] = await election.proofToHashVRF(proofs[i]);
            }

            let cmh = await election.getCurrentMHeight();
            let currentMHeight = cmh.toNumber();
            console.log("cuurent mheight", currentMHeight);

            let winner = -1; // just one of all winners in an epoch.
            let mh = -1;

            for (let j = currentMHeight + challengePeriod; j < currentMHeight + epochLength; j++) {
                totalEpochWins = 0;
                epochWins = [];
                for (let i = 0; i < honestValidatorNum; i++) {
                    epochWins[i] = 0;
                    let dcheck = await election.auxDifficultyCheck(outputs[i], j, {from: accounts[i]});
                    if (dcheck) {
                        totalWins++;
                        totalEpochWins++;
                        epochWins[i] = 1;
                        winner = i;
                        mh = j;
                    }
                }
                if (mh != -1) {
                    for (let i = 0; i < honestValidatorNum; i++) {
                        if (epochWins[i] == 1) {
                            wins[i] += (1.0 / totalEpochWins);
                            console.log("epoch winner:", i);
                        }
                    }
                    break;
                }
            }
            console.log("mh", mh);

            // one winner submits checkpoint.
            if (mh != -1) {
                await waitUntilYouSeeBlock(mh - 1);
                let errorHappened1 = false;
                let tx;
                try {
                    tx = await election.submitCheckpoint(ctx, currentEpoch, proofs[winner], mh, {from: accounts[winner]});
                } catch (e) {
                    errorHappened1 = true;
                }
                if (errorHappened1) {
                    console.log("submitCheckpoint ERROR");
                }
                else {
                    totalSubmitGas += tx.receipt.gasUsed;
                    totalSubmit++;
                }
            }
            // innervate.
            else {
                await waitUntilYouSeeBlock(currentMHeight + epochLength - 1);
                let errorHappened2 = false;
                try {
                    await election.innervate();
                } catch (e) {
                    errorHappened2 = true;
                }
                if (errorHappened2) {
                    console.log("innervate ERROR");
                }
            }
        }

        console.log("________________");

        if (totalSubmit > 0) {
            averageSubmitGas = totalSubmitGas / totalSubmit;
        }
        console.log("avg gas", averageSubmitGas);

        console.log("success rate", totalSubmit / totalTestEpochs);

        console.log("failed attempts rate", (totalWins - totalSubmit) / totalTestEpochs);

        console.log("weights", weights);

        console.log("wins", wins);

    });
});
