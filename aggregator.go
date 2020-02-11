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

func aggregate(config Config, seenFile *SeenFile) {
	parser := gofeed.NewParser()

	for _, feedConfig := range config.Feeds {
		feed, err := parser.ParseURL(feedConfig.Url)

		if err != nil {
			fmt.Println("AGGREGATE ERROR " + err.Error() + " (" + feedConfig.Url + ")")
			continue
		}

		fmt.Println("AGGREGATE " + feed.Title + " (" + feedConfig.Url + ")")

		sessionId := getSessionId(config)

		for _, item := range feed.Items {
			if !seenFile.IsPresent(item.Link) {
				addTorrent(config, sessionId, item.Link)
				logTorrent(item.Link)

				seenFile.Add(item.Link)
			}
		}
	}
}
