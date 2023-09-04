package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
)

type Configuration struct {
	Db struct {
		Ip          string `toml:"DbIp"`
		Port        string `toml:"DbPort"`
		Name        string `toml:"DbName"`
		User        string `toml:"DbUser"`
		Password    string `toml:"DbPassword"`
		AutoMigrate bool   `toml:"DbAutoMigrate"`
	}

	Listen struct {
		Ip   string `toml:"HttpIp"`
		Port string `toml:"HttpPort"`
	}

	Template struct {
		BaseName string `toml:"BaseTemplateName"`
	}
}

// LoadConfig loads and returns a configuration struct
func LoadConfig() Configuration {
	c := flag.String("c", "env.toml", "Path to the toml configuration file")
	flag.Parse()
	file, err := os.ReadFile(*c)
	if err != nil {
		panic("Unable to read TOML config file: " + err.Error())
	}

	var Config Configuration
	_, err = toml.Decode(string(file), &Config)
	if err != nil {
		panic("Unable to decode TOML config file: " + err.Error())
	}

	return Config
}
