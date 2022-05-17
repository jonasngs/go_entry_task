package config

import (
    "io/ioutil"
    "path/filepath"
	"log"
    "gopkg.in/yaml.v3"
)

type Config struct {
	TCP_Server string `yaml:"tcp_server"`
	MySQL_Server MySQL  `yaml:"database"`
	Redis_Server Redis `yaml:"redis"`
}

type MySQL struct {
	Driver string `yaml:"driver"`
	DB_name string `yaml:"db_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	Max_connections int `yaml:"max_connections"`
	Max_idle_connections int `yaml:"max_idle_connections"`
}

type Redis struct {
	Address string `yaml:"address"`
	Password string `yaml:"password"`
	DB int `yaml:"db"`
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

func GetMYSQLServer() MySQL {
	return config.MySQL_Server
}

func GetRedisServer() Redis {
	return config.Redis_Server
}