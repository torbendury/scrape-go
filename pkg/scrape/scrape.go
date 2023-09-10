package scrape

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	// with Go 1.21 this moves into the standard library
	"golang.org/x/exp/slices"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	baseUrl            string
	scrapeUrls         bool
	linkDepth          int
	allowUrlDuplicates bool
	scrapeImages       bool
	urlList            []string
	imageFileList      []string
}

type ScrapeResults struct {
	Links      []string
	ImagePaths []string
}

func NewScraper(baseUrl string, scrapeUrls bool, linkDepth int, allowUrlDuplicates bool, scrapeImages bool) *Scraper {
	return &Scraper{
		scrapeUrls:         scrapeUrls,
		linkDepth:          linkDepth,
		allowUrlDuplicates: allowUrlDuplicates,
		scrapeImages:       scrapeImages,
		urlList:            make([]string, 0),
		imageFileList:      make([]string, 0),
		baseUrl:            baseUrl,
	}
}

func (s *Scraper) StartScrape() error {
	err := s.startUrlScrape()
	if err != nil {
		return err
	}
	if !s.allowUrlDuplicates {
		fmt.Printf("Original len:\t%d\n", len(s.urlList))
		s.urlList = removeDuplicateStr(s.urlList)
		fmt.Printf("Cleaned len:\t%d\n", len(s.urlList))
	}
	err = s.startImageScrape()
	if err != nil {
		return err
	}
	return nil
}

func (s *Scraper) startUrlScrape() error {
	// Perform first request to get a URL pool for scraping
	results, err := s.scrapeUrl(s.baseUrl)
	if err != nil {
		return err
	}
	s.urlList = append(s.urlList, results...)
	// dirty hack: iterate over all URLs * linkDepth - this causes some unnecessary duplicate traffic
	for i := 0; i < s.linkDepth; i++ {
		fmt.Printf("Link depth %v\n", i+1)
		for _, url := range s.urlList {
			if !strings.HasPrefix(url, "https") || !strings.HasPrefix(url, s.baseUrl) {
				continue
			}
			res, err := s.scrapeUrl(url)
			if err != nil {
				return err
			}
			s.urlList = append(s.urlList, res...)
		}
	}
	return nil
}

func (s *Scraper) startImageScrape() error {
	return nil
}

func (s *Scraper) PrintResults() {
	if s.scrapeUrls {
		fmt.Println("URL list:")
		for _, url := range s.urlList {
			fmt.Println(url)
		}
	}
	if s.scrapeImages {
		fmt.Println("Image files list:")
		fmt.Printf("%v", s.imageFileList)
	}
}

func (s *Scraper) scrapeUrl(baseUrl string) ([]string, error) {
	results := make([]string, 0)
	rand.Seed(time.Now().Unix())
	userAgent := userAgents[rand.Intn(len(userAgents))]
	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("status code error when initializing url scrape")
	}
	// Read and parse response data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	doc.Find("a").Each(func(i int, selector *goquery.Selection) {
		// for each item found, get the url
		url, found := selector.Attr("href")
		if found {
			if strings.HasPrefix(url, "/") && url != "/" {
				res := s.baseUrl + url
				if !slices.Contains(s.urlList, res) && !slices.Contains(results, res) {
					results = append(results, res)
				}
			} else if url != "/" && url != "" {
				res := s.baseUrl + url
				if !slices.Contains(s.urlList, res) && !slices.Contains(results, res) {
					results = append(results, url)
				}
			}
		}
	})
	return results, nil
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
