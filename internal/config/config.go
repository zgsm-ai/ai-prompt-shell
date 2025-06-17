package config

import (
	"time"

	"github.com/spf13/viper"
)

/**
 * Main application configuration struct
 * Includes all subsystem configurations
 */
type Config struct {
	AppName string        `mapstructure:"APP_NAME"`
	Server  ServerConfig  `mapstructure:"SERVER"`
	Logger  LoggerConfig  `mapstructure:"LOGGER"`
	Redis   RedisConfig   `mapstructure:"REDIS"`
	Refresh RefreshConfig `mapstructure:"REFRESH"`
	LLM     LLMConfig     `mapstructure:"LLM"`
}

type LoggerConfig struct {
	LogLevel    string `mapstructure:"LEVEL"`
	LogFormat   string `mapstructure:"FORMAT"`
	LogOutput   string `mapstructure:"OUTPUT"`
	LogFileName string `mapstructure:"FILE_NAME"`
}

/**
 * Server related configuration
 */
type ServerConfig struct {
	ListenAddr string `mapstructure:"LISTEN_ADDR"`
	Debug      bool   `mapstructure:"DEBUG"`
}

/**
 * Redis connection configuration
 */
type RedisConfig struct {
	Addr     string `mapstructure:"ADDR"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

/**
 * Auto-refresh intervals configuration
 */
type RefreshConfig struct {
	Tool      time.Duration `mapstructure:"TOOL"`
	Extension time.Duration `mapstructure:"EXTENSION"`
	Prompt    time.Duration `mapstructure:"PROMPT"`
	Environ   time.Duration `mapstructure:"ENVIRON"`
}

/**
 * LLM API configuration
 */
type LLMConfig struct {
	ApiKey  string `mapstructure:"API_KEY"`
	ApiBase string `mapstructure:"API_BASE"`
}

var cfg *Config

/**
 * Load and parse configuration from file/environment
 * @return Pointer to config instance
 */
func Load() *Config {
	if cfg != nil {
		return cfg
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
