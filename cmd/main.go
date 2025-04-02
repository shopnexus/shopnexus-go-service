package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/logger"
	grpc "shopnexus-go-service/internal/server"
	"time"

	"github.com/getsentry/sentry-go"
)

const defaultConfigFile = "config/config.dev.yml"
const productionConfigFile = "config/config.production.yml"

var configFile string

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	setUpConfig()
	setupLogger()
	setupSentry()
	setupGrpcServer()
}

func setUpConfig() {
	fmt.Println("APP_STAGE", os.Getenv("APP_STAGE"))

	if os.Getenv("APP_STAGE") == "production" {
		configFile = productionConfigFile
	} else {
		configFile = defaultConfigFile
	}

	log.Default().Printf("Using config file: %s", configFile)
	config.SetConfig(configFile)
}

func setupLogger() {
	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
}

func setupGrpcServer() {
	logger.Log.Info("Starting gRPC server at port " + fmt.Sprintf("%d", *port))
	server, err := grpc.NewServer(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	server.Start(*port)
}

func setupSentry() {
	logger.Log.Info("Setting up Sentry...")
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.GetConfig().Sentry.Dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for tracing.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	sentry.CaptureMessage("It works!")
}
