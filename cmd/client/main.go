package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"shopnexus-go-service/gen/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()

	// Set up the credentials for the connection.
	// perRPC := oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchToken())}
	// creds, err := credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	// if err != nil {
	// 	log.Fatalf("failed to load credentials: %v", err)
	// }

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	rgc := pb.NewRefundClient(conn)

	res, err := rgc.Get(context.Background(), &pb.GetRefundRequest{
		RefundId: "123e4567-e89b-12d3-a456-426614174000",
	})

	fmt.Println("HALLO FROM GRPC SERVER :d ", res)

	// res, err := rgc.Create(context.Background(), &pb.CreatePaymentRequest{
	// 	ProductIds: [][]byte{
	// 		uuidToBytes("123e4567-e89b-12d3-a456-426614174000"),
	// 	},
	// 	PaymentMethod: pb.CreatePaymentRequest_MOMO.Enum(),
	// })

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	// log.Printf("PaymentId: %s %s", res.GetPaymentId(), res.GetPaymentUrl())
}

func uuidToBytes(uuid string) []byte {
	return []byte(uuid)
}
