package main

import (
	"errors"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CALENDAR_ID string        `yaml:"CALENDAR_ID"`
	Discord     DiscordConfig `yaml:"discord"`
	Line        LineConfig    `yaml:"line"`
	Options     OptionConfig  `yaml:"options"`
}

type DiscordConfig struct {
	Enable     bool     `yaml:"enable"`
	TOKEN      string   `yaml:"TOKEN"`
	Webhooks   []string `yaml:"webhook"`
	ChannelIDs []int64  `yaml:"channel_ids"`
}

type LineConfig struct {
	Enable bool   `yaml:"enable"`
	TOKEN  string `yaml:"TOKEN"`
}

type OptionConfig struct {
	AdvanceReminder     bool `yaml:"advance_reminder"`
	AdvanceReminderDays int  `yaml:"advance_reminder_days"`
}

var (
	ConfigData = &Config{}
	StopWatch  = false
)

const (
	ConfigFilePath = "./data/config.yaml"
	TmpFilePath    = "./data/tmp"
)

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
		Discord:     DiscordConfig{Enable: false},
		Line:        LineConfig{Enable: false},
		Options: OptionConfig{
			AdvanceReminder:     true,
			AdvanceReminderDays: 7,
		},
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

func GetTmpDate() string {
	if data, err := os.ReadFile(TmpFilePath); err == nil {
		return string(data)
	}
	return "" // if file not exist
}

func WriteTmpDate(date time.Time) {
	os.WriteFile(TmpFilePath, []byte(date.Format("2006-01-02")), 0o644)
}
