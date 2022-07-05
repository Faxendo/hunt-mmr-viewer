package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DefaultFolder string `yaml:"default_folder"`
	SendToElastic bool   `yaml:"send_to_elastic"`
}

func loadConfig() Config {
	var cfg Config

	if _, err := os.Stat("./config.yml"); errors.Is(err, os.ErrNotExist) {
		cfg.DefaultFolder = "C:/Program Files (x86)/Steam/steamapps/common/Hunt Showdown/user/profiles/default"
		cfg.SendToElastic = false
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
