package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {

	// // Register servers
	// // pb.RegisterPaymentServer(s, &PaymentServer{})
	// pb.RegisterRefundServer(s, &RefundServer{})

	// log.Printf("server listening at %v", lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
}

func setUpConfig() {
	fmt.Println("APP_STAGE", os.Getenv("APP_STAGE"))

	// if os.Getenv("APP_STAGE") == "production" {
	// 	configFile = productionConfigFile
	// } else {
	// 	configFile = defaultConfigFile
	// }

	// log.Default().Printf("Using config file: %s", configFile)
	// config.SetConfig(configFile)
}

func setupGrpcServer() (*grpc.Server, net.Listener) {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	return s, lis
}
