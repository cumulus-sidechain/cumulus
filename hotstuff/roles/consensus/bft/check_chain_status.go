package bft

import (
	"log"

	"github.com/juju/loggo"
)

var ConsensusLogger = loggo.GetLogger("Consensus")

func (chain *ChainStatus) CheckValidHeightAndPath(block *Block, mainChainHighest uint64, isMainChain bool) bool {
	if chain.CheckIfExists(block) {
		ConsensusLogger.Warningf("Block already exists in the chain")
		log.Println("Block already exists in the chain")
		return false
	}
	// be careful of unsigned subtractin overflow!
	if mainChainHighest > commitThreshold && block.Height < mainChainHighest-commitThreshold {
		ConsensusLogger.Warningf("Block height lower than latest committed block height in the main chain")
		log.Println("Block height lower than latest committed block height in the main chain")
		return false
	}
	if isMainChain {
		if !chain.CheckIfExists(chain.MapHashToBlock[string(block.QC.BlockHash)]) {
			ConsensusLogger.Warningf("Blocks parent is not in the main chain")
			log.Println("Blocks parent is not in the main chain")
			return false
		}
		if block.Height != block.QC.Height+1 {
			ConsensusLogger.Warningf("Parent height is invalid for the main chain")
			log.Println("Parent height is invalid for the main chain")
			return false
		}
		// add to cache chain
		if block.Height > chain.Highest+1 {
			ConsensusLogger.Warningf("Block height not in the valid range for the main chain")
			log.Println("Block height not in the valid range for the main chain")
			return false
		}
	}

	// add to the main
	return true
}

// checkIfExists
func (chain *ChainStatus) CheckIfExists(block *Block) bool {
	if _, ok := chain.MapHashToBlock[string(block.Hash)]; ok { // block already included
		return true
	}
	return false
}
