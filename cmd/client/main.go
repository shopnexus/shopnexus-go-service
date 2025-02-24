package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func TokenAuth(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Logging request
	fmt.Println("📡 Client Request:", method)

	// Attach metadata with Bearer Token
	md := metadata.Pairs("authorization", "Bearer "+"nigga")
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Call next interceptor or actual request
	return invoker(ctx, method, req, reply, cc, opts...)
}

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(TokenAuth),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rgc := pb.NewProductClient(conn)

	res, err := rgc.GetProduct(context.Background(), &pb.GetProductRequest{
		Id: util.ToPtr(int64(1)),
	})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	fmt.Println(res)
}
