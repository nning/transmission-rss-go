package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Feeds []struct {
		Url string `yaml:"url"`
	} `yaml:"feeds"`

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
}

type RequestBody struct {
	Method    string `json:"method"`
	Arguments struct {
		Filename string `json:"filename"`
	} `json:"arguments"`
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func loadConfig(configPath string) Config {
	yamlData, err := ioutil.ReadFile(configPath)
	panicOnError(err)

	var config Config
	err = yaml.Unmarshal(yamlData, &config)
	panicOnError(err)

	return config
}

func getUrl(config Config) string {
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

func getSessionId(config Config) string {
	client := &http.Client{}

	url := getUrl(config)

	request, err := http.NewRequest("GET", url, nil)
	panicOnError(err)

	request.SetBasicAuth(config.Login.Username, config.Login.Password)
	response, err := client.Do(request)
	panicOnError(err)

	_, err = ioutil.ReadAll(response.Body)
	panicOnError(err)

	sessionId := response.Header["X-Transmission-Session-Id"][0]

	fmt.Println("SESSION_ID " + sessionId)

	return sessionId
}

func parseFeeds(config Config) {
	parser := gofeed.NewParser()

	for _, feedConfig := range config.Feeds {
		feed, _ := parser.ParseURL(feedConfig.Url)

		fmt.Println("AGGREGATE " + feedConfig.Url + " (" + feed.Title + ")")

		sessionId := getSessionId(config)

		for _, item := range feed.Items {

			client := &http.Client{}

			url := getUrl(config)

			var requestBody RequestBody
			requestBody.Method = "torrent-add"
			requestBody.Arguments.Filename = item.Link

			jsonData, err := json.Marshal(requestBody)
			panicOnError(err)

			request, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
			panicOnError(err)

			request.SetBasicAuth(config.Login.Username, config.Login.Password)

			request.Header.Add("Content-Type", "application/json")
			request.Header.Add("X-Transmission-Session-Id", sessionId)

			_, err = client.Do(request)
			panicOnError(err)

			fmt.Println("LINK " + item.Link)
		}
	}
}

func main() {
	var configPath = flag.String("c", "transmission-rss.conf", "Config file path")
	var help = flag.Bool("h", false, "Print help")
	var once = flag.Bool("o", false, "Run once")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *configPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	config := loadConfig(*configPath)

	if *once {
		parseFeeds(config)
	} else {
		for {
			parseFeeds(config)
			time.Sleep(5 * time.Second)
		}
	}
}
