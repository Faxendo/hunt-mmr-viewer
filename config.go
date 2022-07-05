package main

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DefaultFolder   string `yaml:"default_folder"`
	SendToElastic   bool   `yaml:"send_to_elastic"`
	ElasticUrl      string `yaml:"elastic_url"`
	ElacticUser     string `yaml:"elactic_user"`
	ElasticPassword string `yaml:"elactic_password"`
	LastHashMatch   int    `yaml:"last_hash_match"`
	UUID            string `yaml:"uuid"`
}

func loadConfig() Config {
	var cfg Config

	if _, err := os.Stat("./config.yml"); errors.Is(err, os.ErrNotExist) {
		cfg.DefaultFolder = "C:/Program Files (x86)/Steam/steamapps/common/Hunt Showdown/user/profiles/default"
		cfg.SendToElastic = false
		cfg.LastHashMatch = 0
		cfg.ElasticUrl = "https://elasticserver.local:9200"
		cfg.ElacticUser = "user"
		cfg.ElasticPassword = "password"
		cfg.UUID = uuid.New().String()
		saveConfig(cfg)
		return cfg
	}

	f, err := os.Open("config.yml")
	check(err)
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	check(err)

	return cfg
}

func saveConfig(config Config) {
	f, err := os.Create("config.yml")
	check(err)
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	err = encoder.Encode(config)
	check(err)
}
