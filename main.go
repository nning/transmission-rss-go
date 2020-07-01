package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

var (
	configPath = ""
	singleRun  = false
	resetSeen  = false
)

func main() {
	app := &cli.App{
		Name:  "transmission-rss",
		Usage: "Transmission RSS is basically a workaround for transmission's lack of the ability to monitor RSS feeds and automatically add enclosed torrent links.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config, c",
				Value:       "transmission-rss.conf",
				Usage:       "Config file path",
				Destination: &configPath,
			},
			&cli.BoolFlag{
				Name:        "reset, r",
				Value:       false,
				Usage:       "Reset seen file on startup",
				Destination: &resetSeen,
			},
			&cli.BoolFlag{
				Name:        "signle-run, s",
				Value:       false,
				Usage:       "Start with daemon",
				Destination: &singleRun,
			},
		},
		Action: cli.ActionFunc(run),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func run(ctx *cli.Context) error {
	config := NewConfig(configPath)
	seenFile := NewSeenFile()

	updateInterval := config.UpdateInterval
	if updateInterval == 0 {
		updateInterval = 600
	}

	aggregator := NewAggregator(config, seenFile)

	if resetSeen {
		seenFile.Clear()
		fmt.Println("SEEN reset")
	}

	if singleRun {
		aggregator.Run()
	} else {
		for {
			aggregator.Run()
			time.Sleep(time.Duration(updateInterval) * time.Second)
		}
	}

	return nil
}
