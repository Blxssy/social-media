package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string     `yaml:"env" envDefault:"local"`
	Database Database   `yaml:"database"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Redis    Redis      `yaml:"redis"`
	Token    Token      `yaml:"token"`
}

type Database struct {
	Dialect   string `yaml:"dialect" default:"postgres"`
	Host      string `yaml:"host" default:"localhost"`
	Port      string `yaml:"port"`
	Name      string `yaml:"name"`
	Username  string `yaml:"username"`
	Password  string
	Migration bool `yaml:"migration"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Token struct {
	accessTokenTTL  time.Duration `yaml:"accessTokenTTL"`
	refreshTokenTTL time.Duration `yaml:"refreshTokenTTL"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.SetDefault("Database.Password", os.Getenv("DB_PASSWORD"))

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}
