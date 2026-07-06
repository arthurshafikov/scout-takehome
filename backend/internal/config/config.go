package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	AppProductionEnv = "production"
	AppLocalEnv      = "local"
	AppTestingEnv    = "testing"
)

type Config struct {
	App             `mapstructure:",squash"`
	SQLiteConfig    `mapstructure:",squash"`
	MinIOConfig     `mapstructure:",squash"`
	APIKeyConfig    `mapstructure:",squash"`
	ThumbnailConfig `mapstructure:",squash"`
}

type App struct {
	Env   string `mapstructure:"APP_ENV"`
	Debug bool   `mapstructure:"APP_DEBUG"`
	Port  string `mapstructure:"APP_PORT"`
}

type SQLiteConfig struct {
	DBPath string `mapstructure:"SQLITE_DB_PATH"`
}

type MinIOConfig struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket    string `mapstructure:"MINIO_BUCKET"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

type APIKeyConfig struct {
	APIKey string `mapstructure:"API_KEY"`
}

type ThumbnailConfig struct {
	CacheTTLHours int `mapstructure:"THUMBNAIL_CACHE_TTL_HOURS"`
}

func NewConfig(envFolderPath, configFolder string) *Config {
	config := &Config{}

	if envFolderPath == "" {
		config.readEnvVarsFromSystem()
	} else {
		config.readEnvVarsFromFile(envFolderPath)
	}

	// Validate required fields
	if config.DBPath == "" {
		log.Fatalln("SQLITE_DB_PATH environment variable is required")
	}
	if config.APIKey == "" {
		log.Fatalln("API_KEY environment variable is required")
	}
	if config.Endpoint == "" {
		log.Fatalln("MINIO_ENDPOINT environment variable is required")
	}

	// Default values
	if config.Env == "" {
		config.Env = AppProductionEnv
	}
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.Bucket == "" {
		config.Bucket = "scout"
	}
	if config.CacheTTLHours == 0 {
		config.CacheTTLHours = 24
	}

	switch config.Env {
	case AppLocalEnv:
	case AppProductionEnv:
	case AppTestingEnv:
	default:
		log.Fatalln("unknown APP_ENV: " + config.Env)
	}

	return config
}

func (c *Config) readEnvVarsFromFile(envFolderPath string) {
	viper.AddConfigPath(envFolderPath)
	viper.SetConfigName("main")
	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(c); err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) readEnvVarsFromSystem() {
	c.Debug = os.Getenv("APP_DEBUG") == "true"
	c.Env = os.Getenv("APP_ENV")
	c.Port = os.Getenv("APP_PORT")
	c.DBPath = os.Getenv("SQLITE_DB_PATH")
	c.APIKey = os.Getenv("API_KEY")
	c.Endpoint = os.Getenv("MINIO_ENDPOINT")
	c.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	c.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	c.Bucket = os.Getenv("MINIO_BUCKET")
	c.UseSSL = os.Getenv("MINIO_USE_SSL") == "true"
}
