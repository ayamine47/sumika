package config

import (
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Discord struct {
		Token   string
		Status  string
	}
	Guild struct {
		Prefix string
		Channel string
	}
	NextCloud struct {
		Enable   bool
		Url      string
		Username string
		Password string
		Path     string
	}
	GoogleDrive struct {
		Enable     bool
		SecretFile string
		TokenFile  string
	}
	UrlWhiteList []string `yaml:"url_whitelist"`
}

const configFile = "./config.yaml"

var CurrentConfig Config

func init() {
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed: ", err)
	}

	//verify
	if CurrentConfig.Discord.Token == "" {
		log.Fatal("Token is empty")
	}

	_, _ = strings.CutSuffix(CurrentConfig.NextCloud.Url, "/")
	_, _ = strings.CutSuffix(CurrentConfig.NextCloud.Path, "/")
}
