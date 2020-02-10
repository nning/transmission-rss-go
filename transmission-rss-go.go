package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
  "github.com/mmcdole/gofeed"
)

type Config struct {
	Feeds []struct {
		Url string `yaml:"url"`
	} `yaml:"feeds"`
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

func parseFeeds(config Config) {
	parser := gofeed.NewParser()

	for _, feedConfig := range config.Feeds {
		feed, _ := parser.ParseURL(feedConfig.Url)

		fmt.Println("AGGREGATE " + feedConfig.Url + " (" + feed.Title + ")")

		// for _, item := range feed.Items {
		// 	fmt.Println("LINK " + item.Link)
		// }
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
