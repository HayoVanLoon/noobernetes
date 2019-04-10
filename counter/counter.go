package main

import (
	"fmt"
	pb "github.com/HayoVanLoon/go-generated/noobernetes/v1"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = "8080"
)

type server struct {
	id    string
	ticks int
	tocks int
}

func (s *server) Inc(t bool) {
	if t {
		s.ticks += 1
	} else {
		s.tocks += 1
	}
}

func (s server) String() string {
	return fmt.Sprintf("%v: %v / %v", s.id, s.ticks, s.tocks)
}

func (s *server) PutTick(ctx context.Context, r *pb.PutTickRequest) (*pb.PutTickResponse, error) {
	resp := &pb.PutTickResponse{Request: r}
	s.Inc(r.Message == "tick")
	return resp, nil
}

func (s server) GetTicks(ctx context.Context, r *pb.GetTicksRequest) (*pb.GetTicksResponse, error) {
	resp := &pb.GetTicksResponse{Request: r, Ticks: int64(s.ticks), Tocks: int64(s.tocks)}
	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	id, _ := uuid.NewUUID()
	pb.RegisterCounterServer(s, &server{id: id.String()})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
