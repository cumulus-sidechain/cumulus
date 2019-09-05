from web3 import Web3
import pprint

def test1(blse_contract):
    iteration_times = 100
    data = Web3.sha3(b'\x74\x78\x74')
    for i in range(0, iteration_times):
        hash_data = data
        s1 = data
        s2 = data

        gas_used = blse_contract.functions.verifyBGLS3(i+1).estimateGas()
        print(gas_used)

def test2(blse_contract):
    iteration_times = 100
    for i in range(0, iteration_times):
        s1 = 385846518441062319503502284295243290270560187383398932887791670182362540842
        s2 = 385846518441062319503502284295243290270560187383398932887791670182362540842

        gas_used = blse_contract.functions.verifyBGLS3(s1, s2, i).estimateGas()
        print(gas_used)

if __name__ == "__main__":
    provider_ip_address = "http://127.0.0.1:7545"
    w3 = Web3(Web3.HTTPProvider(provider_ip_address))

    contract_abi = '''
[
	{
		"constant": false,
		"inputs": [
			{
				"name": "i",
				"type": "int256"
			}
		],
		"name": "verifyBGLS3",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"type": "function",
		"stateMutability": "nonpayable"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "number",
				"type": "uint256"
			}
		],
		"name": "testAdd",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"type": "function",
		"stateMutability": "nonpayable"
	}
]
    '''
    contract_address = w3.toChecksumAddress("0xd0634799919f397b1de9e4e3c68a8d23e9c42c10")

    private_key = "2c62db2f0607295ceb568a3a04a779dbecaf60416a3a12f2905e839eb0d2757e"
    wallet_address = "0xd3d9A80b37222399A512098Ed14662107d6c2725"

    w3.eth.defaultAccount = w3.eth.accounts[1]

    blse_contract = w3.eth.contract(address = contract_address, abi = contract_abi)
    test1(blse_contract)
    # insert_nodes(patricia_contract)

    # index = 8
    # key = bytes("k" + str(index), 'utf-8')
    # value = bytes("v" + str(index), 'utf-8')
    # verify_proof(patricia_contract, key, value)