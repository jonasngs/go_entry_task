package config

import (
    "io/ioutil"
    "path/filepath"
	"log"
    "gopkg.in/yaml.v3"
)

type Config struct {
	TCP_Server string `yaml:"tcp_server"`
}

var config Config

func OpenConfigFile() {
	config_file, _ := filepath.Abs("config/config.yml")
    yamlFile, err := ioutil.ReadFile(config_file)
    err = yaml.Unmarshal(yamlFile, &config)
    if err != nil {
        log.Printf("Error %s when opening config file: ", err)
    }
}

func GetTCPServer() string {
	return config.TCP_Server
}