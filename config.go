package main

import (
	"fmt"
	"io/ioutil"

	"net/url"

	"github.com/nning/transmission-rss-go/utils"
	"gopkg.in/yaml.v2"
)

type Feed struct {
	Url            string  `yaml:"url"`
	RegExp         string  `yaml:"regexp"`
	SeedRatioLimit float32 `yaml:"seed_ratio_limit"`
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

	Paused bool `yaml:"add_paused"`
}

func NewConfig(configPath string) *Config {
	utils.TouchIfNotExist(configPath, defaultConf)

	yamlData, err := ioutil.ReadFile(configPath)
	panicOnError(err)

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	panicOnError(err)

	config.UpdateInterval = utils.ValueOrDefaultInt(config.UpdateInterval, 600)
	config.Server.Host = utils.ValueOrDefaultString(config.Server.Host, "localhost")
	config.Server.RpcPath = utils.ValueOrDefaultString(config.Server.RpcPath, "/transmission/rpc")
	config.Server.Port = utils.ValueOrDefaultInt(config.Server.Port, 9091)

	return &config
}

func (config *Config) ServerURL() string {
	uri := url.URL{
		Host: fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Path: config.Server.RpcPath,
	}
	if config.Server.Tls {
		uri.Scheme = "https"
	} else {
		uri.Scheme = "http"
	}
	if len(config.Login.Username) > 0 {
		uri.User = url.UserPassword(config.Login.Username, config.Login.Password)
	}

	return uri.String()
}
