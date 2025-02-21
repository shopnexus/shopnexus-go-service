package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"shopnexus-go-service/config"
	pgxutil "shopnexus-go-service/internal/db/pgx"
	"shopnexus-go-service/internal/logger"
	"shopnexus-go-service/internal/repository"
)

const defaultConfigFile = "config/config.dev.yml"
const productionConfigFile = "config/config.production.yml"

var configFile string

func main() {
	setUpConfig()

	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
	pool := pgxutil.GetPgxPool()

	repo := repository.NewRepository(pool)

	ctx := context.Background()
	repo.GetBrand(ctx, 1)
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
