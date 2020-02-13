package main

import (
	"flag"
	"fmt"
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
	var singleRun = flag.Bool("s", false, "Single run mode")
	var resetSeen = flag.Bool("r", false, "Reset seen file on startup")

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

	seenFile := NewSeenFile()
	aggregator := NewAggregator(&config, seenFile)

	if *resetSeen {
		seenFile.Clear()
		fmt.Println("SEEN reset")
	}

	if *singleRun {
		aggregator.Run()
	} else {
		for {
			aggregator.Run()
			time.Sleep(time.Duration(updateInterval) * time.Second)
		}
	}
}
