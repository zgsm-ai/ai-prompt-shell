package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config 应用全局配置
type Config struct {
	AppName  string        `mapstructure:"APP_NAME"`
	LogLevel string        `mapstructure:"LOG_LEVEL"`
	Server   ServerConfig  `mapstructure:"SERVER"`
	Redis    RedisConfig   `mapstructure:"REDIS"`
	Refresh  RefreshConfig `mapstructure:"REFRESH"`
}

type ServerConfig struct {
	ListenAddr string `mapstructure:"LISTEN_ADDR"`
	Debug      bool   `mapstructure:"DEBUG"`
}

// RedisConfig Redis连接配置
type RedisConfig struct {
	Addr     string `mapstructure:"ADDR"`
	Password string `mapstructure:"PASSWORD"`
	DB       int    `mapstructure:"DB"`
}

type RefreshConfig struct {
	Tool      time.Duration `mapstructure:"TOOL"`
	Extension time.Duration `mapstructure:"EXTENSION"`
	Prompt    time.Duration `mapstructure:"PROMPT"`
	Environ   time.Duration `mapstructure:"ENVIRON"`
}

var cfg *Config

// Load 加载配置
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
