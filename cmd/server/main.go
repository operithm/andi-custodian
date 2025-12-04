// cmd/server/main.go
package main

import (
	"andi-custodian/internal/chain"
	"andi-custodian/internal/wallet"
	"context"
	"github.com/tyler-smith/go-bip39"
	"log"
	"net"

	pb "andi-custodian/api/custody/v1"
	"andi-custodian/internal/custody"
	"andi-custodian/internal/store"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCustodyServiceServer
	service *custody.Service
}

func (s *server) Transfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	result, err := s.service.Transfer(ctx, &custody.TransferRequest{
		ID:    req.Id,
		Chain: chain.Chain(req.Chain),
		From:  req.From,
		To:    req.To,
		Value: req.Value,
	})
	if err != nil {
		return nil, err
	}
	return &pb.TransferResponse{
		TxId:   result.TxID,
		Status: result.Status,
	}, nil
}
func main() {
	// Use a fixed mnemonic for deterministic demo behavior
	testMnemonic := "slab lonely fish push bomb festival open oval empower federal slot hotel"
	testSeed := bip39.NewSeed(testMnemonic, "")

	// Initialize dependencies
	store := store.NewInMemoryStore()
	signer := wallet.NewSimulatedMPCSigner(testSeed)
	service := custody.NewService(signer, store)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCustodyServiceServer(s, &server{service: service})
	log.Println("Starting gRPC server on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
