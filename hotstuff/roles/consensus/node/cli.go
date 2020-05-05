package node

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"time"

	"github.com/dapperlabs/flow-consensus-research/hotstuff/roles/consensus/bft"
	"github.com/dapperlabs/flow-consensus-research/hotstuff/roles/consensus/crypto"
	"github.com/dapperlabs/flow-consensus-research/hotstuff/roles/consensus/server"
)

type CommandLine struct{}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage: first launch 2 replica, then launch a leader")
	fmt.Println(" leader -port PORT -target TARGET - start a node as first leader")
	fmt.Println(" replica -port PORT -target TARGET - start a node as a repica")
	fmt.Println(" replica2 -port PORT -target TARGET - start a node as a second repica2")
	fmt.Println(" byzantine -port PORT -target TARGET - start a node as a byzantine node")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func InitializeNode(nodeRole, port string, id uint, stake *big.Int) {
	fmt.Printf("Setting up a %s node \n", nodeRole)

	table := ConfigPeerTable()
	for _, v := range table.Nodes() {
		log.Println(v.Address())
	}

	// set up a node
	newNode := InitNode(port, id, stake, table)
	newServer := server.Server{}
	newServer.EventBus = newNode.EventBus
	newServer.MainChain = newNode.MainChain

	// generate the genesis block and add it to main chain
	genesisBlock := bft.NewGenesisBlock()
	newNode.MainChain.AddBlock(genesisBlock)

	// kick-off let node 0 be the first leader
	if nodeRole == "leader" {
		// bootstrap
		// build a genesis block at height 0, also a genesis QC pointing to genesis block at height 1
		genesisQC := bft.NewGenesisQC(genesisBlock.GetBlockHash())
		newNode.HighestQC = genesisQC

		go FirstProposal(newNode)
	}

	go newNode.HandleBlockProposalEvent()
	go newNode.HandleViewChangeEvent()
	go newNode.HandleVoteEvent()

	newServer.SetupServer(port)
}

func InitializeByzantineNode(port string, id uint, stake *big.Int) {
	log.Printf("Setting up a byzantine node \n")

	table := ConfigPeerTable()
	for _, v := range table.Nodes() {
		log.Println(v.Address())
	}

	// set up a node
	newNode := InitNode(port, id, stake, table)
	newServer := server.Server{}
	newServer.EventBus = newNode.EventBus
	newServer.MainChain = newNode.MainChain

	// generate the genesis block and add it to main chain
	genesisBlock := bft.NewGenesisBlock()
	newNode.MainChain.AddBlock(genesisBlock)

	newByzantine := Byzantine{
		Node: newNode,
	}

	go newByzantine.HandleBlockProposalEvent()
	go newByzantine.HandleViewChangeEvent()
	go newByzantine.HandleVoteEvent()

	newServer.SetupServer(port)
}

func (cli *CommandLine) Run() {
	cli.validateArgs()

	leaderCmd := flag.NewFlagSet("leader", flag.ExitOnError)
	replicaCmd := flag.NewFlagSet("replica", flag.ExitOnError)
	replica2Cmd := flag.NewFlagSet("replica2", flag.ExitOnError)
	byzantineCmd := flag.NewFlagSet("byzantine", flag.ExitOnError)

	leaderPort := leaderCmd.String("port", ":50001", "The port to set the server")
	replicaPort := replicaCmd.String("port", ":50002", "The port to set the server")
	replica2Port := replica2Cmd.String("port", ":50003", "The port to set the server")
	byzantinePort := byzantineCmd.String("port", ":50004", "The port to set the server")

	switch os.Args[1] {
	case "leader":
		err := leaderCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "replica":
		err := replicaCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "replica2":
		err := replica2Cmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "byzantine":
		err := byzantineCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if leaderCmd.Parsed() {
		InitializeNode("leader", *leaderPort, 0, big.NewInt(10))
	}

	if replicaCmd.Parsed() {
		InitializeNode("replica", *replicaPort, 1, big.NewInt(10))
	}

	if replica2Cmd.Parsed() {
		InitializeNode("replica2", *replica2Port, 2, big.NewInt(10))
	}

	if byzantineCmd.Parsed() {
		InitializeByzantineNode(*byzantinePort, 2, big.NewInt(10))
	}
}

func FirstProposal(n *Node) {
	// make first block using genesis QC
	firstBlock := bft.NewBlock(n.BlockPayloadGenerator.Payload(), n.HighestQC, time.Now().Unix())
	// log.Println("make first block at height", firstBlock.Height)
	// log.Println("block hash", firstBlock.GetBlockHash())
	// add first block to main chain
	// n.MainChain.AddBlock(firstBlock)

	fmt.Println("first proposal qc height is ", n.HighestQC.Height)
	sig := crypto.SignMsg(firstBlock, n.identity.ID())
	genesisBlockProposal := &bft.BlockProposal{
		Block:     firstBlock,
		Signature: sig,
	}

	// nextLeader := n.getNextPrimary()

	// voteSig := crypto.SignMsg(firstBlock, n.identity.ID())
	// vote := &bft.Vote{
	// 	BlockHash: firstBlock.GetBlockHash(),
	// 	Signature: sig,
	// }
	log.Println("broadcast first block")
	server.ProposeBlock(n.identity.ID(), n.IdentityTable, genesisBlockProposal)
	// log.Println("send vote for first block to next leader ", nextLeader.ID())
	// server.SendVote(nextLeader.Address(), vote)
}
