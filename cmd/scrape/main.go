package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/torbendury/scrape-go/pkg/scrape"
)

func main() {
	scrapeUrls := flag.Bool("urls", false, "Set to true if you want to scrape deep link URLs.")
	urlFile := flag.String("url-outfile", "urls.txt", "The file to write scraped URLs to.")
	baseUrl := flag.String("base-url", "", "The base URL to start with.")
	linkDepth := flag.Int("link-depth", 5, "Maximum scraping depth.")
	allowUrlDuplicates := flag.Bool("allow-duplicate-urls", false, "Allow duplicated links. Only takes effect if URL scraping is active.")
	scrapeImages := flag.Bool("images", false, "NOT IMPLEMENTED YET: Set to true if you want to scrape images.")
	flag.Parse()

	if *baseUrl == "" || *linkDepth < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	scraper := scrape.NewScraper(*baseUrl, *scrapeUrls, *urlFile, *linkDepth, *allowUrlDuplicates, *scrapeImages)

	fmt.Printf("ScrapeUrls:\t\t%v\nScrapeImages:\t\t%v\nLinkDepth:\t\t%v\nAllowUrlDuplicates:\t%v\n", *scrapeUrls, *scrapeImages, *linkDepth, *allowUrlDuplicates)
	err := scraper.StartScrape()
	if err != nil {
		panic(err)
	}
	scraper.SaveResults()
}
