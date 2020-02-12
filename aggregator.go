package main

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/mmcdole/gofeed"
)

func logTorrent(link string) {
	u, err := url.Parse(link)
	panicOnError(err)

	desc := filepath.Base(u.Path)
	if u.Scheme == "magnet" {
		desc = u.Query().Get("dn")
	}

	fmt.Println("ADD_TORRENT " + desc)
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
			link := item.Link

			if len(item.Enclosures) > 0 {
				link = item.Enclosures[0].URL
			}

			if !seenFile.IsPresent(link) {
				addTorrent(config, sessionId, link)
				logTorrent(link)

				seenFile.Add(link)
			}
		}
	}
}
