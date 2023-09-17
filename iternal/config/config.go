package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Env     string `yaml:"env" env-default:"dev"`
	IsDebug *bool  `yaml:"is_debug"`
	Listen  struct {
		Type string `yaml:"type"`

		Port string `yaml:"port"`
	} `yaml:"listen"`

	Storage StorageConfig `yaml:"storage"`
}
type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Создание синглтона
var instance *Config
var once sync.Once

func GetConfig() *Config {
	//код ниже выполниться только 1 раз
	once.Do(func() {
		log.Println("read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}
	})
	return instance
}
