package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

const (
	AppProductionEnv = "production"
	AppLocalEnv      = "local"
	AppTestingEnv    = "testing"
	SSLModePrefer    = "prefer"
	TimeZoneUTC      = "UTC"
)

type Config struct {
	App      `mapstructure:",squash"`
	DBConfig `mapstructure:",squash"`
}

type App struct {
	Env                       string `mapstructure:"APP_ENV"`
	Debug                     bool   `mapstructure:"APP_DEBUG"`
	Port                      string `mapstructure:"APP_PORT"`
	MaxEventHandlerGoroutines int
}

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DBName   string `mapstructure:"DB_NAME"`
	Port     int64  `mapstructure:"DB_PORT"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
	TimeZone string `mapstructure:"DB_TIMEZONE"`
}

func NewConfig(envFolderPath, configFolder string) *Config {
	var config Config

	viper.AddConfigPath(configFolder)
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalln(err)
	}

	if envFolderPath == "" {
		config.readEnvVarsFromSystem()
	} else {
		config.readEnvVarsFromFile(envFolderPath)
	}

	// Default values
	if config.DBConfig.Port == 0 {
		config.DBConfig.Port = 3306
	}
	if config.DBConfig.SSLMode == "" {
		config.DBConfig.SSLMode = SSLModePrefer
	}
	if config.DBConfig.TimeZone == "" {
		config.DBConfig.TimeZone = TimeZoneUTC
	}

	switch config.App.Env {
	case AppLocalEnv:
	case AppProductionEnv:
	case AppTestingEnv:
		break
	default:
		panic("unknown APP_ENV type " + config.App.Env)
	}

	return &config
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
	if c.App.Env == "" {
		c.App.Env = AppProductionEnv
	}

	if os.Getenv("DB_PORT") != "" {
		portInt64, err := strconv.ParseInt(os.Getenv("DB_PORT"), 10, 64)
		if err != nil {
			log.Fatalln(err)
		}
		c.DBConfig.Port = portInt64
	}
	c.DBConfig.Host = os.Getenv("DB_HOST")
	c.DBConfig.User = os.Getenv("DB_USER")
	c.DBConfig.Password = os.Getenv("DB_PASSWORD")
	c.DBConfig.DBName = os.Getenv("DB_NAME")
	c.DBConfig.SSLMode = os.Getenv("DB_SSL_MODE")
	c.DBConfig.TimeZone = os.Getenv("DB_TIMEZONE")
}
