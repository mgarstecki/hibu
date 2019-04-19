package main

import (
	"fmt"

	"github.com/mgarstecki/hibu/internal/crawler"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()

	myCrawler := crawler.NewCrawler(fs)

	results := make(chan crawler.CrawlResult)

	go myCrawler.Crawl(".", results)

	for res := range results {
		fmt.Printf("%s: %x (%v)\n", res.File, res.Hash, res.Err)
	}
}
