package main

import (
	"flag"
	"os"

	"github.com/torbendury/scrape-go/pkg/scrape"
)

func main() {
	allowUrlDuplicates := flag.Bool("allow-duplicate-urls", false, "Allow duplicated links. Only takes effect if URL scraping is active.")
	baseUrl := flag.String("base-url", "", "The base URL to start with.")
	imagesDirectory := flag.String("images-dir", "./images/", "The directory to save scraped images to.")
	linkDepth := flag.Int("link-depth", 5, "Maximum scraping depth.")
	scrapeImages := flag.Bool("images", false, "Set to true if you want to scrape images.")
	scrapeUrls := flag.Bool("urls", false, "Set to true if you want to scrape deep link URLs.")
	urlFile := flag.String("url-outfile", "urls.txt", "The file to write scraped URLs to.")
	flag.Parse()

	if *baseUrl == "" || *linkDepth < 1 || len(*urlFile) < 1 || len(*imagesDirectory) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	scraper := scrape.NewScraper(*baseUrl, *scrapeUrls, *urlFile, *linkDepth, *allowUrlDuplicates, *scrapeImages, *imagesDirectory)

	err := scraper.StartScrape()
	if err != nil {
		panic(err)
	}
	scraper.SaveResults()
}
