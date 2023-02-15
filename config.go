package main

import (
	"errors"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CALENDAR_ID string `yaml:"CALENDAR_ID"`
	Discord     struct {
		Enable     bool     `yaml:"enable"`
		Webhook    []string `yaml:"webhook"`
		ChannelIDs []int64  `yaml:"channel_ids"`
	} `yaml:"discord"`
	Line struct {
		Enable bool `yaml:"enable"`
		// TODO add line support
	} `yaml:"line"`
}

var (
	ConfigData = &Config{}
	StopWatch  = false
)

const ConfigFilePath = "./config.yaml"

func init() {
	// os.ReadFile("./config.json")

	yamlFile, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
		config := NewConfig()
		data, _ := yaml.Marshal(config)
		os.WriteFile(ConfigFilePath, data, 0o644)
	} else {
		yaml.Unmarshal(yamlFile, &ConfigData)
	}

	// watch config file
	go func() {
		for !StopWatch {
			err := watchFile(ConfigFilePath)
			if err != nil {
				log.Println("config watch error", err)
			}

			// reload config
			yamlFile, err := os.ReadFile(ConfigFilePath)
			if err != nil {
				log.Println("config watch error", err)
			} else {
				yaml.Unmarshal(yamlFile, &ConfigData)
			}
		}
	}()
}

func NewConfig() *Config {
	return &Config{
		CALENDAR_ID: "<請輸入自己的日曆編號>",
	}
}

// is ease watch file func
func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		// sleep 1 second
		time.Sleep(time.Second)
	}
	return nil
}
