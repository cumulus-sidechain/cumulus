package node

type Byzantine struct {
	*Node
}

func (b *Byzantine) HandleBlockProposalEvent() {

}

func (b *Byzantine) HandleViewChangeEvent() {

}

func (b *Byzantine) HandleVoteEvent() {

}

func (b *Byzantine) PretendOffLine(duration int64) {

}
