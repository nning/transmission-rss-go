package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Feed struct {
	Url string `yaml:"url"`
}

type Config struct {
	Feeds []Feed `yaml:"feeds"`

	Server struct {
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
		Tls     bool   `yaml:"tls"`
		RpcPath string `yaml:"rpc_path"`
	} `yaml:"server"`

	Login struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"login"`

	UpdateInterval int `yaml:"update_interval"`

	Paused bool `yaml:"add_paused"` // TODO
}

// TODO NewConfig
func loadConfig(configPath string) Config {
	yamlData, err := ioutil.ReadFile(configPath)
	panicOnError(err)

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	panicOnError(err)

	return config
}

func getUrl(config *Config) string {
	url := ""

	if config.Server.Tls {
		url += "https://"
	} else {
		url += "http://"
	}

	url += config.Server.Host
	url += config.Server.RpcPath

	return url
}
