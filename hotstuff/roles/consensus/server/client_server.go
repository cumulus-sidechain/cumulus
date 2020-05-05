package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	service_pb "github.com/cumulus-sidechain/cumulus/hotstuff/gen/grpc/services/consensus"
	shared_pb "github.com/cumulus-sidechain/cumulus/hotstuff/gen/grpc/shared"
	"github.com/cumulus-sidechain/cumulus/hotstuff/identity"
	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/bft"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type Server struct {
	*EventBus
	MainChain  *bft.ChainStatus
	CacheChain *bft.ChainStatus
}

type EventBus struct {
	BlockProposalEvent chan *bft.BlockProposal
	ViewChangeEvent    chan *bft.ViewChange
	VoteEvent          chan *bft.Vote
}

type messenger func(client service_pb.ConsensusServiceClient, ctx context.Context) (interface{}, error)

func sendMsg(target string, mg messenger) {
	log.Printf("start sending message...")
	// Set up a connection to replicas
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	log.Printf("dialing")
	defer conn.Close()
	client := service_pb.NewConsensusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := mg(client, ctx)
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}
	log.Printf("message %v sent to: %v", r, target)
}

func serializeQC(qc *bft.QC) *shared_pb.QuorumCertificate {
	return &shared_pb.QuorumCertificate{
		BlockHash: qc.BlockHash,
		AggregatedSignatures: &shared_pb.AggregatedSignatures{
			RawAggregatedSignature: qc.AggregatedSignature.RawSignature,
			Signers:                qc.AggregatedSignature.Signers,
		},
	}
}

func serializeSignature(signature *bft.Signature) *shared_pb.Signature {
	return &shared_pb.Signature{
		RawSignature: signature.RawSignature,
		Signer:       signature.Signer,
	}
}

func serializeBlock(block *bft.Block) *shared_pb.Block {
	// log.Println("serialized block height is ", block.Height)
	// NOTE serializedQC height is default 0 as it does not have a Height field in protobuff
	serializedQC := serializeQC(block.QC)
	return &shared_pb.Block{
		Payload:           block.GetPayload(),
		Height:            block.Height,
		PreviousBlockHash: block.PrevHash,
		Qc:                serializedQC,
	}
}

func deserializeQC(serializedQC *shared_pb.QuorumCertificate) *bft.QC {
	blockHash := serializedQC.GetBlockHash()

	return &bft.QC{
		AggregatedSignature: &bft.AggregatedSignature{
			RawSignature: serializedQC.GetAggregatedSignatures().GetRawAggregatedSignature(),
			Signers:      serializedQC.GetAggregatedSignatures().GetSigners(),
		},
		BlockHash: blockHash,
	}
}

func deserializeSignature(serializedSig *shared_pb.Signature) *bft.Signature {
	return &bft.Signature{
		RawSignature: serializedSig.GetRawSignature(),
		Signer:       serializedSig.GetSigner(),
	}

}

func deserializeBlock(serializedBlock *shared_pb.Block) *bft.Block {
	qc := deserializeQC(serializedBlock.GetQc())
	// as qc in protobuff does not have a height, need to recover qc height from block height
	qc.Height = serializedBlock.Height - 1

	block := bft.NewBlock(serializedBlock.GetPayload(), qc, serializedBlock.GetTimestamp().GetSeconds())
	// log.Println("deserialized block height is ", block.Height)
	return block
}

func SendVote(target string, vote *bft.Vote) {
	mg := func(client service_pb.ConsensusServiceClient, ctx context.Context) (interface{}, error) {
		// construct serialized voteRequest
		serializedSig := serializeSignature(vote.Signature)
		blockHash := vote.BlockHash

		voteRequest := &service_pb.ProcessVoteRequest{
			BlockHash: blockHash,
			Signature: serializedSig,
		}

		return client.ProcessVote(ctx, voteRequest)
	}
	sendMsg(target, mg)
}

func ProposeBlock(selfId uint, table identity.Table, blockProposal *bft.BlockProposal) {
	// time.Sleep(3 * time.Second)

	serializedBlock := serializeBlock(blockProposal.Block)
	// log.Println("serialize block", serializedBlock)
	serializedSig := serializeSignature(blockProposal.Signature)
	// log.Println("serialize sig", serializedSig)
	proposalRequest := &service_pb.ProcessBlockProposalRequest{
		Block:     serializedBlock,
		Signature: serializedSig,
	}

	mg := func(client service_pb.ConsensusServiceClient, ctx context.Context) (interface{}, error) {
		// construct serialized block Proposal
		log.Println("send block through network at height", blockProposal.Block.Height)
		log.Println("blockhash ", blockProposal.Block.GetBlockHash())

		return client.ProcessBlockProposal(ctx, proposalRequest)
	}

	for _, node := range table.Nodes() {
		go sendMsg(node.Address(), mg)
	}

}

func RequestBlock(target string, blockHash []byte) *bft.Block {
	// Set up a connection to the server.
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := service_pb.NewConsensusServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.QueryBlock(ctx, &service_pb.QueryBlockRequest{BlockHash: blockHash})
	if err != nil {
		log.Fatalf("could not query block: %v", err)
	}
	serializedBlock := r.Block
	return deserializeBlock(serializedBlock)
}

func (n *Server) QueryBlock(ctx context.Context, blockQueryRequest *service_pb.QueryBlockRequest) (*service_pb.QueryBlockReply, error) {
	block := n.MainChain.MapHashToBlock[string(blockQueryRequest.GetBlockHash())]
	log.Println("someone asked for block ", block)
	serializedBlock := serializeBlock(block)
	return &service_pb.QueryBlockReply{
		Block: serializedBlock,
	}, nil
}

func ProposeViewChange(target string, vc *bft.ViewChange) {
	mg := func(client service_pb.ConsensusServiceClient, ctx context.Context) (interface{}, error) {
		serializedQC := serializeQC(vc.QC)
		serializedSig := serializeSignature(vc.Signature)

		viewChangeRequest := &service_pb.ProcessViewChangeRequest{
			HighestQc: serializedQC,
			Signature: serializedSig,
		}
		return client.ProcessViewChange(ctx, viewChangeRequest)
	}
	sendMsg(target, mg)
}

func (n *Server) ProcessBlockProposal(ctx context.Context, newProposal *service_pb.ProcessBlockProposalRequest) (*empty.Empty, error) {
	// log.Println("A new block received")

	serializedBlock := newProposal.GetBlock()
	block := deserializeBlock(serializedBlock)
	sig := deserializeSignature(newProposal.GetSignature())

	blockProposal := bft.BlockProposal{
		Block:     block,
		Signature: sig,
	}

	n.EventBus.BlockProposalEvent <- &blockProposal

	// log.Println("process block finished")
	return &empty.Empty{}, nil
}

func (n *Server) ProcessViewChange(ctx context.Context, newViewChange *service_pb.ProcessViewChangeRequest) (*empty.Empty, error) {
	log.Println("View change received")

	qc := deserializeQC(newViewChange.GetHighestQc())
	sig := deserializeSignature(newViewChange.GetSignature())

	viewChange := bft.ViewChange{
		QC:        qc,
		Signature: sig,
	}

	n.EventBus.ViewChangeEvent <- &viewChange

	log.Println("process view change finished")
	return &empty.Empty{}, nil

}

func (n *Server) ProcessVote(ctx context.Context, newVote *service_pb.ProcessVoteRequest) (*empty.Empty, error) {
	// log.Println("A new vote received")

	blockHash := newVote.GetBlockHash()
	sig := deserializeSignature(newVote.GetSignature())

	vote := bft.Vote{
		BlockHash: blockHash,
		Signature: sig,
	}

	n.EventBus.VoteEvent <- &vote

	// log.Println("process vote finished")
	return &empty.Empty{}, nil
}

func (n *Server) SetupServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service_pb.RegisterConsensusServiceServer(s, n)
	fmt.Println("Start listening on", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
