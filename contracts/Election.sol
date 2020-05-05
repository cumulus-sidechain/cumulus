pragma solidity >=0.5.3 <0.7.0;

import "./VRF.sol";


contract Election {

	enum CPStatus {Unsettled, Revoked, Settled, Empty}

	struct ECert {
		uint256 output; // VRF hash (output).
		uint256[4] proof; // VRF proof.
		uint256[2] publicKey; // public key of the submitter.
		uint256 mh; // minimum mainchain height that this certificate is acceptable to the contract.
	}

	struct Checkpoint {
		bool isExistent; // true if the checkpoint exists.
		CPStatus status; // status of the checkpoint.
		bytes ctx; // the latest state of SC.
		ECert eCert; // the VRF certificate of the checkpoint.
    }

	mapping (address => uint256[2]) private publicKeys; // mapping address to public key.
	mapping (address => uint256) private weights; // mapping address to weight.
	mapping (address => bool) private isValidator; // determines if an address is a validator.
	mapping (uint256 => Checkpoint) private checkpoints; // mapping from epoch number to checkpoint.
	mapping (uint256 => uint256) private mHeights; // mapping from epoch number to mainchain height.
	mapping (uint256 => uint256) private seeds; // mapping from epoch number to seed.
	uint256 private totalWeight; // the total weights of all the validators.
	uint256 private currentEpoch; // the current epoch number.
	uint256 private epochLength; // the number of blocks that the epoch goes on (no submission is accepted after).
	uint256 private challengePeriod; // the number of blocks for the challenge period (no challenge is accepted after).
	uint256 private difficultyParam; // the difficulty parameter.
	bool private protocolStarted; // determines whether the protocol is started or not.

	// auxiliary function to add accounts as validators before the start of the protocol.
	// the implementation is basic and subject to change.
	function addAsValidator(uint256[2] memory _publicKey, uint256 _weight) public {
		require(!protocolStarted);
		require(!isValidator[msg.sender]);
		publicKeys[msg.sender] = _publicKey;
		weights[msg.sender] = _weight;
		totalWeight += _weight;
		isValidator[msg.sender] = true;
	}

	// auxiliary function to set the parameters and start the protocol.
	// the implementation is basic and subject to change.
	function startProtocol(uint256 _seed, uint256 _difficultyParam, uint256 _epochLength, uint256 _challengePeriod) public {
		seeds[0] = _seed;
		mHeights[0] = block.number + 1;
		currentEpoch = 0;
		difficultyParam = _difficultyParam;
		epochLength = _epochLength;
		challengePeriod = _challengePeriod;
		protocolStarted = true;
	}

	// auxiliary function to be called by a user, checks if the given proof is valid.
	function auxVerifyVRF(uint256[4] memory _proof) public view returns (bool) {
		return verifyVRF(publicKeys[msg.sender], _proof, seeds[currentEpoch]);
	}

	// auxiliary function to check if the given output and mainchain-height provided satisfy difficulty requirement.
	function auxDifficultyCheck(uint256 _output, uint256 _mh) public view returns (bool) {
		return (hashTwoNumbers(_output, _mh) % (difficultyParam * totalWeight) <  weights[msg.sender]);
	}

	// gas-efficient conversion of uint256 to bytes.
	function toBytes(uint256 x) public pure returns (bytes memory b) {
		b = new bytes(32);
		assembly { mstore(add(b, 32), x) }
	}

	// VRF verification given the public key, VRF proof, and seed.
	function verifyVRF(uint256[2] memory _publicKey, uint256[4] memory _proof, uint256 _seed) public pure returns (bool) {
		return VRF.verify(_publicKey, _proof, toBytes(_seed));
	}

	// VRF hash (output) computation given the VRF proof.
	function proofToHashVRF(uint256[4] memory _proof) public pure returns (uint256) {
		return uint256(VRF.gammaToHash(_proof[0], _proof[1]));
	}

	// given two numbers, the function outputs a hash of them using sha256.
	function hashTwoNumbers(uint256 _n1, uint256 _n2) public pure returns (uint256) {
		return uint256(sha256(abi.encodePacked(_n1, _n2)));
	}

	// sets the seed and height and enters the next epoch.
	function enterNextEpoch(uint256 _seed) private {
		seeds[currentEpoch + 1] = _seed;
		mHeights[currentEpoch + 1] = block.number + 1;
		currentEpoch += 1;
	}

	// checkpoint submission and verification.
	function submitCheckpoint(bytes memory _ctx, uint256 _epochNum, uint256[4] memory _proof, uint256 _mh) public {
		// the verification of the checkpoint.
		require(_epochNum == currentEpoch); // this line should be first to minimize gas usage of honest failed attempts.
		require(protocolStarted);
		require(isValidator[msg.sender]);
		require(!checkpoints[currentEpoch].isExistent);
		require(block.number < mHeights[currentEpoch] + epochLength);
		require(_mh <= block.number);
		require(mHeights[currentEpoch] + challengePeriod <= _mh);
//		require(verifyVRF(publicKeys[msg.sender], _proof, seeds[currentEpoch]));
		uint256 _output = proofToHashVRF(_proof);
		// changed the requirement formula: now on average, every difficulty block in a row, one weight wins.
		require(hashTwoNumbers(_output, _mh) % (difficultyParam * totalWeight) <  weights[msg.sender]);

		// the checkpoint is accepted.
		if(checkpoints[currentEpoch - 1].status == CPStatus.Unsettled) {
			checkpoints[currentEpoch - 1].status = CPStatus.Settled;
		}
		ECert memory _eCert;
		_eCert.output = _output;
		_eCert.proof = _proof;
		_eCert.publicKey = publicKeys[msg.sender];
		_eCert.mh = _mh;
		Checkpoint memory _checkpoint;
		_checkpoint.isExistent = true;
		_checkpoint.status = CPStatus.Unsettled;
		_checkpoint.ctx = _ctx;
		_checkpoint.eCert = _eCert;
		checkpoints[currentEpoch] = _checkpoint;
		enterNextEpoch(_output);
	}

	// empty checkpoint after timeout.
	function innervate() public {
		require(protocolStarted);
		require(block.number >= mHeights[currentEpoch] + epochLength);
		// empty checkpoint.
		checkpoints[currentEpoch].status = CPStatus.Empty;
		checkpoints[currentEpoch].isExistent = true;
		enterNextEpoch(hashTwoNumbers(seeds[currentEpoch], currentEpoch + 1));
	}

	// returns the total weight, difficulty parameter, length of each epoch, and the challenge period.
	function getProtocolParams() public view returns (uint256 _totalWeight, uint256 _difficultyParam,
		uint256 _epochLength, uint256 _challengePeriod, bool _protocolStarted) {
		return (totalWeight, difficultyParam, epochLength, challengePeriod, protocolStarted);
	}

	// returns the current epoch number.
	function getCurrentEpoch() public view returns (uint256) {
		return currentEpoch;
	}

	// returns the current mainchain height.
	function getCurrentMHeight() public view returns (uint256) {
		return mHeights[currentEpoch];
	}

	// returns the seed of the current epoch.
	function getCurrentSeed() public view returns (uint256) {
		return seeds[currentEpoch];
	}

	// returns ctx of a previous epoch.
	function getPrevCtx(uint256 _epochNum) public view returns (uint256 _status, bytes memory _ctx) {
		require(_epochNum < currentEpoch);
		return (uint256(checkpoints[_epochNum].status), checkpoints[_epochNum].ctx);
	}

}
