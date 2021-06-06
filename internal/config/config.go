package config

import (
	"github.com/tkanos/gonfig"
	"log"
)

var (
	Config Configuration
)

func init () {
	err := GetConfig("./configs/config.json", &Config)
	if err != nil {
		log.Panicf("couldn't get/parse the config : %v", err)
	}
}

type Configuration struct {
	ScrapeInterval   uint64    // minutes
	DbName           string
	DbHost           string
	DbPort           int
	DbUser           string
	DbPassword       string
	DbSSLMode        string
	RedisHost        string // ip:port
	RedisDB          int
	RedditConfigPath string // user agent file
	RSSListPath      string
	RedditListPath   string
	TwitterListPath  string
}

func GetConfig(configPath string, config *Configuration) error {
	err := gonfig.GetConf(configPath, config)
	if err != nil {
		return err
	}

	return nil
}