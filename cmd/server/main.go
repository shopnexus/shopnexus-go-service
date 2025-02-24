package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/server"
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

	err := server.NewServer(fmt.Sprintf(":%d", *port))
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
