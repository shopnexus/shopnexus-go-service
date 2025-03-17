package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var config *Config
var m sync.Mutex

type Config struct {
	Env           string        `yaml:"env"`
	App           App           `yaml:"app"`
	HttpServer    HttpServer    `yaml:"httpServer"`
	Log           Log           `yaml:"log"`
	Postgres      Postgres      `yaml:"postgres"`
	S3            S3            `yaml:"s3"`
	Sentry        Sentry        `yaml:"sentry"`
	SensitiveKeys SensitiveKeys `yaml:"sensitiveKeys"`
	Vnpay         Vnpay         `yaml:"vnpay"`
}

type App struct {
	AccessTokenDuration int64 `yaml:"accessTokenDuration"`
}

type HttpServer struct {
	Port int `yaml:"port"`
}

type Log struct {
	Level           string `yaml:"level"`
	StacktraceLevel string `yaml:"stacktraceLevel"`
	FileEnabled     bool   `yaml:"fileEnabled"`
	FileSize        int    `yaml:"fileSize"`
	FilePath        string `yaml:"filePath"`
	FileCompress    bool   `yaml:"fileCompress"`
	MaxAge          int    `yaml:"maxAge"`
	MaxBackups      int    `yaml:"maxBackups"`
}

type Postgres struct {
	Url             string `yaml:"url"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxConnections  int32  `yaml:"maxConnections"`
	MaxConnIdleTime int32  `yaml:"maxConnIdleTime"`
}

type S3 struct {
	Region          string `yaml:"region"`
	Bucket          string `yaml:"bucket"`
	AccessKeyID     string `yaml:"accessKeyID"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	CloudfrontURL   string `yaml:"cloudfrontURL"`
}

type Sentry struct {
	Dsn         string `yaml:"dsn"`
	Environment string `yaml:"environment"`
	Release     string `yaml:"release"`
	Debug       bool   `yaml:"debug"`
}

type SensitiveKeys struct {
	JWTSecret string `yaml:"jwtSecret" mapstructure:"jwtSecret"`
}

type Vnpay struct {
	TmnCode    string `yaml:"tmnCode"`
	HashSecret string `yaml:"hashSecret"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	/** Because GitHub Actions doesn't have .env, and it will load ENV variables from GitHub Secrets */
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
