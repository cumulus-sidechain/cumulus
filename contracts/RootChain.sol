pragma solidity ^0.4.0;

contract RootChain {

    uint8 constant public VALIDATOR_NUMBER = 5;

    Validator[5] public validators;
    Client[] public clients;

    uint256 public hight = 0;

    mapping (uint256 => Checkpoint) public checkpoints;
    
    struct Validator {
        bytes32 pubaddress;
    }

    struct Client {
        bytes32 pubaddress;
    }

    struct Checkpoint {
        bytes32 root;
        bytes aggresig;
        uint8 bitindex;
    }

    event CheckpointSubmitted(
        uint256 cpNumber
    );

    event DepositCreated(
        address indexed depositor,
        uint256 indexed epochNumber,
        uint32 amount
    );

    event ExitFinalized(
        uint192 indexed exitId
    );

    function submitCheckpoint(bytes32 _blockhash, bytes _asig, uint8 _inum) public {

        // Create the checkpoint.
        checkpoints[hight] = Checkpoint({
            root: _blockhash,
            aggresig: _asig,
            bitindex: _inum
        });

        hight += 1;

        // emit CheckpointSubmitted(hight);
    }

    function submitDeposit(uint32 amount) public payable{
        // console.log("deposit succeed");
        // emit DepositCreated(msg.sender, hight, amount);
    }
}