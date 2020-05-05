package bft

import "log"

const mapSize = 1000
const commitThreshold = 3

type ChainStatus struct {
	MapHashToBlock                  map[string]*Block
	MapHeightToMapBlockToBool       map[uint64]map[*Block]bool
	MapHashToMapChildrenBlockToBool map[string]map[*Block]bool
	Highest                         uint64
}

func NewChainStatus() *ChainStatus {
	return &ChainStatus{
		MapHashToBlock: make(map[string]*Block, mapSize),
		// still need to initialize inner map
		MapHashToMapChildrenBlockToBool: make(map[string]map[*Block]bool, mapSize),
		MapHeightToMapBlockToBool:       make(map[uint64]map[*Block]bool, mapSize),
	}
}

func (chain *ChainStatus) AddBlock(block *Block) bool {
	ConsensusLogger.Infof("Adding the block ", string(block.Hash))
	if block.QC != nil {
		chain.MapHashToMapChildrenBlockToBool[string(block.QC.BlockHash)] = make(map[*Block]bool)
		chain.MapHashToMapChildrenBlockToBool[string(block.QC.BlockHash)][block] = true
	}

	if block.Height > chain.Highest {
		chain.Highest = block.Height
		log.Println("add block and update highest block to ", chain.Highest)
	}

	chain.MapHeightToMapBlockToBool[block.Height] = make(map[*Block]bool)
	chain.MapHeightToMapBlockToBool[block.Height][block] = true
	// chain.MapHeightToMapBlockToBool[block.Height][block] = true
	// We want this line last as the checks for if a block exist in the chain
	// use this mapping, because AddBlock might not have completed for some reason
	chain.MapHashToBlock[string(block.Hash)] = block

	log.Println("Adding the block ", block)
	return true
}

func (chain *ChainStatus) MergeMissingChainFromBlock(cache *ChainStatus, block *Block, inputBlockIsInMainChain bool) bool {
	// ConsensusLogger.Infof("Merging block ", string(block.Hash))

	// Get the children and remove the parent, Trump-style
	childrenBlockMap := cache.MapHashToMapChildrenBlockToBool[string(block.Hash)]
	if !inputBlockIsInMainChain {
		cache.RemoveBlock(block, true)
		chain.AddBlock(block)
	}
	// Go over each child of block
	for childBlock, _ := range childrenBlockMap {
		// Recursively get the children and delete its parents
		chain.MergeMissingChainFromBlock(cache, childBlock, false)
	}

	return true
}

// isRootOfChain should generally always be called as true because the block that's
// being passed in is the root (even if it's a standalone block) - it is false
// for the rest of the chain so the extra delete/lookup isn't needed
func (chain *ChainStatus) RemoveBlock(block *Block, isRootOfChain bool) bool {
	// ConsensusLogger.Infof("Removing block ", string(block.Hash))

	delete(chain.MapHeightToMapBlockToBool[block.Height], block)
	delete(chain.MapHashToBlock, string(block.Hash))
	delete(chain.MapHashToMapChildrenBlockToBool, string(block.Hash))
	if isRootOfChain {
		delete(chain.MapHashToMapChildrenBlockToBool[string(block.QC.BlockHash)], block)
	}

	return true
}

func (chain *ChainStatus) PruneCompetingChains(block *Block) bool {
	ConsensusLogger.Infof("Pruning competing chains ", string(block.Hash))

	if SiblingBlocks, exists := chain.MapHashToMapChildrenBlockToBool[string(block.QC.BlockHash)]; exists {
		for siblingBlock, _ := range SiblingBlocks {
			if string(siblingBlock.Hash) != string(block.Hash) {
				chain.PruneChainFromBlock(siblingBlock, true)
			}
		}
	}

	return true
}

func (chain *ChainStatus) PruneChainFromBlock(block *Block, isRootOfChain bool) bool {
	ConsensusLogger.Infof("Pruning chain from block ", string(block.Hash))

	// Get the children and remove the parent, Trump-style
	childrenBlockMap := chain.MapHashToMapChildrenBlockToBool[string(block.Hash)]
	chain.RemoveBlock(block, isRootOfChain)
	// Go over each child of block
	for childBlock, _ := range childrenBlockMap {
		// Recursively get the children and delete its parents
		chain.PruneChainFromBlock(childBlock, false)
	}

	return true
}
