package services

import (
	"main/types"
	"sort"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
)

func StartNewsUpdater(feedUrls []string, updateChan chan []types.NewsItem, errorChan chan error) {
	for {
		items := make([]types.NewsItem, 0)
		for _, url := range feedUrls {
			nextItems, err := loadNewsItems(url)
			if err != nil {
				errorChan <- err
			} else {
				items = append(items, nextItems...)
			}
		}
		sort.Slice(items, func(a int, b int) bool {
			return items[b].Date.Before(items[a].Date)
		})
		updateChan <- items
		time.Sleep(time.Minute * 5)
	}
}

func loadNewsItems(url string) ([]types.NewsItem, error) {
	items := make([]types.NewsItem, 0)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	} else {
		for _, item := range feed.Items {
			pubdate := time.Now()
			if item.PublishedParsed != nil {
				pubdate = *item.PublishedParsed
			}
			newitem := types.NewsItem{Headline: item.Title, Description: strip.StripTags(item.Description), Date: pubdate, Source: feed.Title}
			items = append(items, newitem)
		}
	}
	return items, nil
}
