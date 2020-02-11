package main

import (
	"fmt"
	"net/url"

	"github.com/mmcdole/gofeed"
)

func logTorrent(link string) {
	u, err := url.Parse(link)
	panicOnError(err)

	fmt.Println("ADD_TORRENT " + u.Query().Get("dn"))
}

func aggregate(config Config) {
	parser := gofeed.NewParser()

	for _, feedConfig := range config.Feeds {
		feed, _ := parser.ParseURL(feedConfig.Url)

		fmt.Println("AGGREGATE " + feed.Title + " (" + feedConfig.Url + ")")

		sessionId := getSessionId(config)

		for _, item := range feed.Items {
			addTorrent(config, sessionId, item.Link)
			logTorrent(item.Link)
		}
	}
}
