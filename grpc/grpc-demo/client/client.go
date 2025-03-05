package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kenneth-wang/go-demo/grpc/grpc-demo/arithpb"
	"google.golang.org/grpc"
)

func main() {
	// 连接 gRPC 服务器
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := arithpb.NewArithClient(conn)

	// 创建请求
	req := &arithpb.MultiplyRequest{A: 6, B: 7}

	// 调用 gRPC 方法
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Multiply(ctx, req)
	if err != nil {
		log.Fatalf("Could not call Multiply: %v", err)
	}

	// 输出结果
	fmt.Printf("gRPC Result: %d * %d = %d\n", req.A, req.B, res.Result)
}
