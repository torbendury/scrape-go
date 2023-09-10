package main

import (
	"flag"
	"fmt"

	"github.com/torbendury/scrape-go/pkg/scrape"
)

func main() {
	scrapeUrls := flag.Bool("urls", false, "Set to true if you want to scrape deep link URLs.")
	baseUrl := flag.String("base-url", "", "The base URL to start with.")
	linkDepth := flag.Int("link-depth", 5, "Maximum scraping depth.")
	allowUrlDuplicates := flag.Bool("allow-duplicate-urls", false, "Allow duplicated links. Only takes effect if URL scraping is active.")
	scrapeImages := flag.Bool("images", false, "NOT IMPLEMENTED YET: Set to true if you want to scrape images.")
	flag.Parse()

	scraper := scrape.NewScraper(*baseUrl, *scrapeUrls, *linkDepth, *allowUrlDuplicates, *scrapeImages)

	fmt.Println("Created scraper config:")
	fmt.Printf("ScrapeUrls:\t\t%v\n", *scrapeUrls)
	fmt.Printf("ScrapeImages:\t\t%v\n", *scrapeImages)
	fmt.Printf("LinkDepth:\t\t%v\n", *linkDepth)
	fmt.Printf("AllowUrlDuplicates:\t%v\n", *allowUrlDuplicates)

	fmt.Println("Starting scrape...")
	err := scraper.StartScrape()
	if err != nil {
		panic(err)
	}
	scraper.PrintResults()
}
