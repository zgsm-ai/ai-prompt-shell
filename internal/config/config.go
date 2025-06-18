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
	AppName string        `mapstructure:"app_name"`
	Env     string        `mapstructure:"env"`
	Server  ServerConfig  `mapstructure:"server"`
	Logger  LoggerConfig  `mapstructure:"logger"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Refresh RefreshConfig `mapstructure:"refresh"`
	LLM     LLMConfig     `mapstructure:"llm"`
}

type LoggerConfig struct {
	LogLevel    string `mapstructure:"level"`
	LogFormat   string `mapstructure:"format"`
	LogOutput   string `mapstructure:"output"`
	LogFileName string `mapstructure:"file_name"`
}

/**
 * Server related configuration
 */
type ServerConfig struct {
	ListenAddr string `mapstructure:"listen_addr"`
	Debug      bool   `mapstructure:"debug"`
}

/**
 * Redis connection configuration
 */
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

/**
 * Auto-refresh intervals configuration
 */
type RefreshConfig struct {
	Tool      time.Duration `mapstructure:"tool"`
	Extension time.Duration `mapstructure:"extension"`
	Prompt    time.Duration `mapstructure:"prompt"`
	Environ   time.Duration `mapstructure:"environ"`
}

/**
 * LLM API configuration
 */
type LLMConfig struct {
	ApiKey  string `mapstructure:"api_key"`
	ApiBase string `mapstructure:"api_base"`
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

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
