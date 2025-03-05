package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/kenneth-wang/go-demo/grpc/grpc-demo/arithpb"
	"google.golang.org/grpc"
)

// ArithServer 实现 gRPC 服务器
type ArithServer struct {
	arithpb.UnimplementedArithServer
}

// Multiply 实现 gRPC 方法
func (s *ArithServer) Multiply(ctx context.Context, req *arithpb.MultiplyRequest) (*arithpb.MultiplyResponse, error) {
	result := req.A * req.B
	return &arithpb.MultiplyResponse{Result: result}, nil
}

func main() {
	// 启动 gRPC 服务器
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	arithpb.RegisterArithServer(grpcServer, &ArithServer{})

	fmt.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
