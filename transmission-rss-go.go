package main

import (
	"flag"
	"os"
	"time"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
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

	updateInterval := config.UpdateInterval
	if updateInterval == 0 {
		updateInterval = 600
	}

	if *once {
		aggregate(config)
	} else {
		for {
			aggregate(config)
			time.Sleep(time.Duration(updateInterval) * time.Second)
		}
	}
}
