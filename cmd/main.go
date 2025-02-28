package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"shopnexus-go-service/config"
	grpcServer "shopnexus-go-service/internal/grpc/server"
	"shopnexus-go-service/internal/http"
	"shopnexus-go-service/internal/logger"
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

	go func() {
		err := grpcServer.NewServer(fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
	}()

	httpServer, err := http.NewServer(fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	err = httpServer.Start(fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatal(err)
	}
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
