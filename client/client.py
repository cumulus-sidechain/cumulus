from web3 import Web3
import pprint

def insert_nodes(patricia_contract):
    iteration_times = 10
    for i in range(iteration_times):
        key = bytes("k" + str(i), 'utf-8')
        value = bytes("v" + str(i), 'utf-8')
        print('inserting key-value', i)
        
        tx_hash = patricia_contract.functions.insert(key, value).transact()
        w3.eth.waitForTransactionReceipt(tx_hash)

        print('Updated root hash: {}'.format(
            patricia_contract.functions.getRootHash().call().hex()
        ))

def get_proof(patricia_contract, key):
    output = patricia_contract.functions.getProof(key).call()
    branch_mask = output[0]
    siblings = output[1]
    
    return branch_mask, siblings

def wait_for_receipt(w3, tx_hash, poll_interval):
   while True:
       tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
       if tx_receipt:
         return tx_receipt
       time.sleep(poll_interval)

def verify_proof(patricia_contract, key, value):
    root_hash = patricia_contract.functions.getRootHash().call()
    proof = get_proof(patricia_contract, key)
    branch_mask = proof[0]
    siblings = proof[1]
    tx_receipt = patricia_contract.functions.verifyProof(root_hash, key, value, branch_mask, siblings).transact()
    return len(siblings), w3.eth.getTransactionReceipt(tx_receipt)['gasUsed']

def test1(patricia_contract):
    iteration_times = 1000
    for i in range(9707, iteration_times * 10):
        key = bytes("k" + str(i), 'utf-8')
        value = bytes("v" + str(i), 'utf-8')

        tx_hash = patricia_contract.functions.insert(key, value).transact()

        proofnum, gas_used = verify_proof(patricia_contract, key, value)
        
        print(i, proofnum, gas_used)

if __name__ == "__main__":
    provider_ip_address = "http://127.0.0.1:7545"
    w3 = Web3(Web3.HTTPProvider(provider_ip_address))

    contract_abi = '''
[
	{
		"constant": false,
		"inputs": [
			{
				"name": "key",
				"type": "bytes"
			},
			{
				"name": "value",
				"type": "bytes"
			}
		],
		"name": "insert",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "key",
				"type": "bytes"
			}
		],
		"name": "getProof",
		"outputs": [
			{
				"name": "branchMask",
				"type": "uint256"
			},
			{
				"name": "_siblings",
				"type": "bytes32[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "getRootHash",
		"outputs": [
			{
				"name": "",
				"type": "bytes32"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "rootHash",
				"type": "bytes32"
			},
			{
				"name": "key",
				"type": "bytes"
			},
			{
				"name": "value",
				"type": "bytes"
			},
			{
				"name": "branchMask",
				"type": "uint256"
			},
			{
				"name": "siblings",
				"type": "bytes32[]"
			}
		],
		"name": "verifyProof",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]
    '''
    contract_address = w3.toChecksumAddress("0x750012C750459980869F8cB52455561a247AEB21")

    private_key = "87c3d0dae18393fd0a41b9be631dfb9f2fb8e17cea217c59a9cf3e255d0722f6"
    wallet_address = "0x01Fe864288c0792230f70B74a569A0b75bBbD2B0"

    w3.eth.defaultAccount = w3.eth.accounts[2]

    patricia_contract = w3.eth.contract(address = contract_address, abi = contract_abi)
    test1(patricia_contract)
    # insert_nodes(patricia_contract)

    # index = 8
    # key = bytes("k" + str(index), 'utf-8')
    # value = bytes("v" + str(index), 'utf-8')
    # verify_proof(patricia_contract, key, value)