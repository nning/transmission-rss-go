package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/mmcdole/gofeed"
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

	self := Aggregator{
		Client:   client,
		Config:   config,
		Parser:   parser,
		SeenFile: seenFile,
	}

	return &self
}

func logTorrent(link string) {
	u, err := url.Parse(link)
	panicOnError(err)

	desc := filepath.Base(u.Path)
	if u.Scheme == "magnet" {
		desc = u.Query().Get("dn")
	}

	fmt.Println("ADD_TORRENT " + desc)
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

	if match(item.Title, feedConfig.RegExp) && !self.SeenFile.IsPresent(link) {
		self.Client.AddTorrent(link)
		logTorrent(link)

		self.SeenFile.Add(link)

		if feedConfig.SeedRatioLimit > 0 {
			// TODO setTorrent
			// self.Client.SetTorrent(link, options)
		}
	}
}

func (self *Aggregator) processFeed(feedConfig *Feed) {
	feed, err := self.Parser.ParseURL(feedConfig.Url)

	if err != nil {
		fmt.Println("AGGREGATE ERROR " + err.Error() + " (" + feedConfig.Url + ")")
		return
	}

	fmt.Println("AGGREGATE " + feed.Title + " (" + feedConfig.Url + ")")

	for _, item := range feed.Items {
		self.processItem(feedConfig, item)
	}
}

func (self *Aggregator) Run() {
	for _, feedConfig := range self.Config.Feeds {
		self.processFeed(&feedConfig)
	}
}
