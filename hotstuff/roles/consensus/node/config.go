package node

import (
	"bufio"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/dapperlabs/flow-consensus-research/hotstuff/identity"
)

const (
	filePath = "./node/peers.conf"
)

func ConfigPeerTable() *identity.InMemoryIdentityTable {
	nodes := identity.NewInMemoryIdentityTable(SetupNodes())

	return nodes
}

func SetupNodes() identity.NodeRecords {
	nodeString := GetPeersFromFile()
	nodes := identity.NodeRecords{}

	for i, v := range nodeString {
		nodeInfo := strings.Split(v, " ")
		stake := new(big.Int)
		nodeStake, ok := stake.SetString(nodeInfo[1], 10)
		if !ok {
			log.Fatalf("cannot convert string to bigInt for stake")
		}
		nodes = append(nodes, identity.NewNodeRecord(uint(i), nodeInfo[0], nodeStake))
	}

	return nodes
}

func GetPeersFromFile() []string {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
