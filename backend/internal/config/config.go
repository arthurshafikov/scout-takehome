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
	App              `mapstructure:",squash"`
	SQLiteConfig     `mapstructure:",squash"`
	MinIOConfig      `mapstructure:",squash"`
	APIKeyConfig     `mapstructure:",squash"`
	ThumbnailConfig  `mapstructure:",squash"`
	MaxEventHandlerGoroutines int
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
	Endpoint   string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey  string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey  string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket     string `mapstructure:"MINIO_BUCKET"`
	UseSSL     bool   `mapstructure:"MINIO_USE_SSL"`
}

type APIKeyConfig struct {
	APIKey string `mapstructure:"API_KEY"`
}

type ThumbnailConfig struct {
	CacheTTLHours int `mapstructure:"THUMBNAIL_CACHE_TTL_HOURS"`
}

func NewConfig(envFolderPath, configFolder string) *Config {
	config := &Config{
		MaxEventHandlerGoroutines: 10,
	}

	if envFolderPath == "" {
		config.readEnvVarsFromSystem()
	} else {
		config.readEnvVarsFromFile(envFolderPath)
	}

	// Validate required fields
	if config.SQLiteConfig.DBPath == "" {
		log.Fatalln("SQLITE_DB_PATH environment variable is required")
	}
	if config.APIKeyConfig.APIKey == "" {
		log.Fatalln("API_KEY environment variable is required")
	}
	if config.MinIOConfig.Endpoint == "" {
		log.Fatalln("MINIO_ENDPOINT environment variable is required")
	}

	// Default values
	if config.App.Env == "" {
		config.App.Env = AppProductionEnv
	}
	if config.App.Port == "" {
		config.App.Port = "8080"
	}
	if config.MinIOConfig.Bucket == "" {
		config.MinIOConfig.Bucket = "scout"
	}
	if config.ThumbnailConfig.CacheTTLHours == 0 {
		config.ThumbnailConfig.CacheTTLHours = 24
	}

	switch config.App.Env {
	case AppLocalEnv:
	case AppProductionEnv:
	case AppTestingEnv:
	default:
		log.Fatalln("unknown APP_ENV: " + config.App.Env)
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
	c.App.Debug = os.Getenv("APP_DEBUG") == "true"
	c.App.Env = os.Getenv("APP_ENV")
	c.App.Port = os.Getenv("APP_PORT")
	c.SQLiteConfig.DBPath = os.Getenv("SQLITE_DB_PATH")
	c.APIKeyConfig.APIKey = os.Getenv("API_KEY")
	c.MinIOConfig.Endpoint = os.Getenv("MINIO_ENDPOINT")
	c.MinIOConfig.AccessKey = os.Getenv("MINIO_ACCESS_KEY")
	c.MinIOConfig.SecretKey = os.Getenv("MINIO_SECRET_KEY")
	c.MinIOConfig.Bucket = os.Getenv("MINIO_BUCKET")
	c.MinIOConfig.UseSSL = os.Getenv("MINIO_USE_SSL") == "true"
}
