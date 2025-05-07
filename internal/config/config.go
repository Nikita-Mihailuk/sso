package config

import (
	"flag"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storagePath" env-required:"true"`
	TokenTTL    time.Duration `yaml:"tokenTTL" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		path := fetchConfigPath()
		if path == "" {
			panic("config path is empty")
		}

		instance = LoadConfigByPath(path)
	})

	return instance
}

func LoadConfigByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}

	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetTypeByDefaultValue(true)

	if err := viper.ReadInConfig(); err != nil {
		panic("config read error")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic("config unmarshal error")
	}

	return &config
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
