package node

import (
	"log"
	"math/big"
	"reflect"
	"sync"
	"time"

	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/mempool"
	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/timer"

	"github.com/cumulus-sidechain/cumulus/hotstuff/identity"
	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/bft"
	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/crypto"
	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/server"
)

const (
	blockSize           = 10
	viewchangethreshold = 3
)

type Node struct {
	identity           *identity.NodeRecord
	IdentityTable      *identity.InMemoryIdentityTable
	CurrentLeaderIndex uint

	MainChain  *bft.ChainStatus
	CacheChain *bft.ChainStatus

	voteThreshold         *big.Int
	HighestQC             *bft.QC
	BlockPayloadGenerator mempool.BlockPayloadGenerator

	ProposedBlockVotes       map[string][]*bft.Vote
	ProposedBlockVoteCounter map[string]*big.Int

	ViewChangeMessageCounter uint

	EventBus *server.EventBus
	timer    *timer.Timer
	mux      sync.Mutex
}

func InitNode(port string, id uint, stake *big.Int, table *identity.InMemoryIdentityTable) *Node {
	return &Node{
		identity:           identity.NewNodeRecord(id, port, stake),
		IdentityTable:      table,
		CurrentLeaderIndex: 0,

		MainChain:  bft.NewChainStatus(),
		CacheChain: bft.NewChainStatus(),

		BlockPayloadGenerator:    mempool.NewDummyBlockPayloadGenerator(blockSize),
		ProposedBlockVoteCounter: map[string]*big.Int{},
		ProposedBlockVotes:       map[string][]*bft.Vote{},

		// voteThreshold: uint((table.Count() - 1) * 1),
		voteThreshold: new(big.Int).Sub(table.TotalStake(), stake),
		// voteThreshold: table.TotalStake(),
		EventBus: &server.EventBus{
			BlockProposalEvent: make(chan *bft.BlockProposal, 1000),
			ViewChangeEvent:    make(chan *bft.ViewChange, 1000),
			VoteEvent:          make(chan *bft.Vote, 1000),
		},
		timer: &timer.Timer{
			Timeout: 10 * time.Second,
		},
	}
}

func (n *Node) HandleBlockProposalEvent() {
	for {
		blockProposal := <-n.EventBus.BlockProposalEvent
		block := blockProposal.Block

		// log.Println("last timer stopped")
		n.timer.ResetTimer()

		// Verify if block height is indeed parent block height (kept locally) + 1
		log.Println("received block height is ", block.Height)
		log.Println("current highest is ", n.MainChain.Highest)
		if n.MainChain.CheckValidHeightAndPath(block, n.MainChain.Highest, true) {
			log.Println("block is ok")
			n.MainChain.AddBlock(block)

			blockHash := block.GetBlockHash()
			qc := block.QC

			// log.Println("receive block at height ", block.Height)
			// log.Println("block hash is", block.GetBlockHash())
			// log.Println("received qc height is ", qc.Height)

			signature := crypto.SignMsg(blockHash, n.identity.ID())

			// update HighestQC
			if n.HighestQC == nil || block.Height >= n.HighestQC.Height {
				n.HighestQC = qc
			}

			// make a vote
			vote := &bft.Vote{
				BlockHash: blockHash,
				Signature: signature,
			}
			// send vote to the next leader
			targetNode := n.getNextPrimary()
			log.Printf("received block, about to send vote, the next leader is %v", targetNode.ID())
			// if block receiver is next leader, do not send vote to self
			// // if targetNode.ID() != n.identity.ID() {
			// // 	log.Println("not next leader, so send vote for received block ", block)

			go server.SendVote(targetNode.Address(), vote)
			// }

			go func() {
				// timeout and start new-view
				<-n.timer.Timer.C
				log.Println("YOU SHOULD NOT SEE THIS")
				nextPrimary := n.getNextPrimary()
				sig := crypto.SignMsg(n.HighestQC, n.identity.ID())
				vc := bft.ViewChange{
					QC:        n.HighestQC,
					Signature: sig,
				}
				server.ProposeViewChange(nextPrimary.Address(), &vc)
				n.timer.ResetTimer()
			}()
		}
	}
}

func (n *Node) HandleViewChangeEvent() {
	for {
		newViewChange := <-n.EventBus.ViewChangeEvent

		receivedQc := newViewChange.QC
		// what if node do not have the block QC is pointing to locally?
		blockHeight := n.MainChain.MapHashToBlock[string(receivedQc.BlockHash)].Height
		receivedQc.Height = blockHeight + 1

		n.ViewChangeMessageCounter++

		log.Println("received viewchange at height ", receivedQc.Height)
		if receivedQc.Height > n.HighestQC.Height {
			n.HighestQC = receivedQc
		}

		// propose a block after receiving enough view change messages
		if n.ViewChangeMessageCounter >= viewchangethreshold {
			newBlock := bft.NewBlock(n.BlockPayloadGenerator.Payload(), n.HighestQC, time.Now().Unix())
			log.Println("MAKE A BLOCK AT HEIGHT ", newBlock.Height)
			// log.Println("block hash ", newBlock.GetBlockHash())

			// Make a block
			// n.MainChain.AddBlock(newBlock)

			sig := crypto.SignMsg(newBlock.GetBlockHash(), n.identity.ID())
			blockProposal := &bft.BlockProposal{
				Block:     newBlock,
				Signature: sig,
			}

			// broadcast block to all peers including itself
			time.Sleep(3 * time.Second)
			go server.ProposeBlock(n.identity.ID(), n.IdentityTable, blockProposal)

			n.ViewChangeMessageCounter = 0
		}
	}
}

func (n *Node) HandleVoteEvent() {
	for {
		vote := <-n.EventBus.VoteEvent
		blockHash := string(vote.BlockHash)

		if voteSender, err := n.IdentityTable.GetByID(uint(vote.Signature.Signer)); err != nil {
			log.Panic("cannot find vote sender identity for stake")
		} else {
			// add stake for votes and store vote for block
			if voteStake, ok := n.ProposedBlockVoteCounter[blockHash]; ok {
				voteStake.Add(voteStake, voteSender.Stake())
				// newVoteStake := big.NewInt(0)
				// newVoteStake.Add(voteStake, voteSender.Stake())
				// n.ProposedBlockVoteCounter[blockHash] = newVoteStake

				// n.printPeerStake()

				log.Println("received vote for a block again from ", voteSender.Address())
				// log.Println("vote stake after accumulating is ", newVoteStake)
			} else {
				// n.ProposedBlockVoteCounter[blockHash] = voteSender.Stake()
				blockStake := big.NewInt(0)
				blockStake.Add(blockStake, voteSender.Stake())
				n.ProposedBlockVoteCounter[blockHash] = blockStake
				// log.Println("vote sender stake is ", voteSender.Stake())
				log.Println("received vote first time for a  block from ", voteSender.Address())
				// log.Println("total stake for this block is now ", n.ProposedBlockVoteCounter[blockHash])
			}
			n.ProposedBlockVotes[blockHash] = append(n.ProposedBlockVotes[blockHash], vote)
		}

		// upon receipt of sufficient votes (in terms of stake)
		if n.ProposedBlockVoteCounter[blockHash].Cmp(n.voteThreshold) >= 0 {
			log.Println("th ", n.voteThreshold)
			log.Println(n.ProposedBlockVoteCounter[blockHash])
			log.Println("collected enough votes for block ", n.MainChain.MapHashToBlock[blockHash])
			if _, ok := n.MainChain.MapHashToBlock[blockHash]; ok != true {
				log.Println("received enough votes for a missing block, query block...")
				// query vote sender for the missing block
				for _, receivedVote := range n.ProposedBlockVotes[blockHash] {
					voteSenderId := receivedVote.Signature.Signer
					if voteSender, err := n.IdentityTable.GetByID(uint(voteSenderId)); err != nil {
						log.Panic("cannot find vote sender identity")
					} else {

						requestedBlock := server.RequestBlock(voteSender.Address(), vote.BlockHash)
						log.Println("queried block ", requestedBlock)
						log.Println("from node ", voteSender.Address())
						// if received requested block has the hash we want, done
						if reflect.DeepEqual(requestedBlock.GetBlockHash(), vote.BlockHash) {
							n.MainChain.AddBlock(requestedBlock)
							break
						}
					}
				}

			}

			// make a qc using collected votes
			qc, err := bft.NewQC(n.MainChain.MapHashToBlock[blockHash].Height, vote.BlockHash, n.ProposedBlockVotes[blockHash])
			if err != nil {
				log.Fatalf("cannot build QC")
			}
			// garbage collection
			delete(n.ProposedBlockVotes, blockHash)
			n.ProposedBlockVoteCounter[blockHash] = big.NewInt(0)

			// build a block on top of newly created qc
			newBlock := bft.NewBlock(n.BlockPayloadGenerator.Payload(), qc, time.Now().Unix())
			log.Println("make a block at height ", newBlock.Height)
			// log.Println("block hash ", newBlock.GetBlockHash())

			// n.MainChain.AddBlock(newBlock)

			sig := crypto.SignMsg(newBlock.GetBlockHash(), n.identity.ID())
			blockProposal := &bft.BlockProposal{
				Block:     newBlock,
				Signature: sig,
			}

			// broadcast block to all peers except self
			log.Println("braodcast block ", newBlock)

			time.Sleep(3 * time.Second)
			go server.ProposeBlock(n.identity.ID(), n.IdentityTable, blockProposal)

			// // make a vote
			// vote := &bft.Vote{
			// 	BlockHash: newBlock.Hash,
			// 	Signature: sig,
			// }

			// // send vote to the next leader
			// targetNode := n.getNextPrimary()

			// log.Println("sending vote to the next leader for newly built block %v", targetNode.ID())
			// go server.SendVote(targetNode.Address(), vote)
		}
	}
}

// Will be replaced by pacemaker
// currently just called by current primary, return next primary (round-robin)
func (n *Node) getNextPrimary() identity.NodeIdentity {
	total := n.IdentityTable.Count()
	n.mux.Lock()
	n.CurrentLeaderIndex++
	n.mux.Unlock()
	nextID := (n.CurrentLeaderIndex) % total
	nextNode, err := n.IdentityTable.GetByID(nextID)
	if err != nil {
		log.Fatalf("cannot get next primary")
	}
	return nextNode
}

func (n *Node) printPeerStake() {
	for _, peer := range n.IdentityTable.Nodes() {
		log.Println("node id ", peer.ID())
		log.Println("node stake ", peer.Stake())
	}
}
