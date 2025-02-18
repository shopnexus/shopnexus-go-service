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
	Scheduler     Scheduler     `yaml:"scheduler"`
	Schedules     []Schedule    `yaml:"schedules"`
	Postgres      Postgres      `yaml:"postgres"`
	Redis         []Redis       `yaml:"redis"`
	Sentry        Sentry        `yaml:"sentry"`
	SensitiveKeys SensitiveKeys `yaml:"sensitiveKeys"`
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

type Label struct {
	En string `json:"en"`
	Th string `json:"th"`
}

type App struct {
	Name     string `yaml:"name"`
	NameSlug string `yaml:"nameSlug"`
}

type Postgres struct {
	Url             string `yaml:"url" mapstructure:"url"`
	Host            string `yaml:"host" mapstructure:"host"`
	Port            int    `yaml:"port" mapstructure:"port"`
	Username        string `yaml:"username" mapstructure:"username"`
	Password        string `yaml:"password" mapstructure:"password"`
	Database        string `yaml:"database" mapstructure:"database"`
	Schema          string `yaml:"schema" mapstructure:"schema"`
	MaxConnections  int32  `yaml:"maxConnections" mapstructure:"maxConnections"`
	MaxConnIdleTime int32  `yaml:"maxConnIdleTime" mapstructure:"maxConnIdleTime"`
}

type Redis struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	Database int    `yaml:"database" mapstructure:"database"`
}

type Sentry struct {
	Dsn         string `yaml:"dsn" mapstructure:"dsn"`
	Environment string `yaml:"environment" mapstructure:"environment"`
	Release     string `yaml:"release" mapstructure:"release"`
	Debug       bool   `yaml:"debug" mapstructure:"debug"`
}

type SensitiveKeys struct {
	TelegramBotToken string `yaml:"telegramBotToken" mapstructure:"telegramBotToken"`
	JWTSecret        string `yaml:"jwtSecret" mapstructure:"jwtSecret"`
	Pepper           string `yaml:"pepper" mapstructure:"pepper"`
	EncryptionKey    string `yaml:"encryptionKey" mapstructure:"encryptionKey"`
	Iterations       int    `yaml:"iterations" mapstructure:"iterations"`
}

type Scheduler struct {
	Timezone string `yaml:"timezone" mapstructure:"timezone"`
}

type Schedule struct {
	Job       string `yaml:"job" mapstructure:"job"`
	Cron      string `yaml:"cron" mapstructure:"cron"`
	IsEnabled bool   `yaml:"isEnabled" mapstructure:"isEnabled"`
}

type Authentication struct {
	Endpoint string `yaml:"endpoint" mapstructure:"endpoint"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
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
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}
