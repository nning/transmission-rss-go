package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/nning/transmission-rss-go/logger"
)

type Aggregator struct {
	Client   *Client
	Config   *Config
	Parser   *gofeed.Parser
	SeenFile *SeenFile
}

func NewAggregator(config *Config, seenFile *SeenFile) *Aggregator {
	client := NewClient(config)
	parser := gofeed.NewParser()
	parser.Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return fmt.Errorf("302")
		},
	}

	self := Aggregator{
		Client:   client,
		Config:   config,
		Parser:   parser,
		SeenFile: seenFile,
	}

	return &self
}

func match(title string, expr string) bool {
	re, err := regexp.Compile(expr)
	panicOnError(err)

	return re.Match([]byte(title))
}

func (self *Aggregator) processItem(feedConfig *Feed, item *gofeed.Item) {
	link := item.Link

	if len(item.Enclosures) > 0 {
		link = item.Enclosures[0].URL
	}

	if !self.SeenFile.IsPresent(link) {
		if !match(item.Title, feedConfig.RegExp) {
			self.SeenFile.Add(link)
			return
		}

		logger.Info("ADD", item.Title)
		id, err := self.Client.AddTorrent(link, feedConfig.DownloadPath)
		if err != nil {
			logger.Error(err)
			return
		}

		self.SeenFile.Add(link)

		if feedConfig.SeedRatioLimit > 0 {
			arguments := make(map[string]interface{})

			arguments["ids"] = []int{id}
			arguments["seedRatioLimit"] = feedConfig.SeedRatioLimit
			arguments["seedRatioMode"] = 1

			self.Client.SetTorrent(arguments)
		}
	}
}

func (self *Aggregator) processFeed(feedConfig *Feed) {
	logger.Info("Fetching", feedConfig.Url)
	feed, err := self.Parser.ParseURL(feedConfig.Url)

	if err != nil {
		logger.Error("Fetching", err.Error())
		return
	}

	logger.Info("Found", len(feed.Items), "items")

	for _, item := range feed.Items {
		self.processItem(feedConfig, item)
	}
}

func (self *Aggregator) Run() {
	for _, feedConfig := range self.Config.Feeds {
		self.processFeed(&feedConfig)
	}
}
