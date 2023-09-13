package scrape

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	// with Go 1.21 this moves into the standard library
	"golang.org/x/exp/slices"

	"github.com/PuerkitoBio/goquery"
	"github.com/torbendury/scrape-go/pkg/utility"
)

// Scraper mainly holds scraping information about URLs, images and file paths. It gets constructed by NewScraper().
type Scraper struct {
	allowUrlDuplicates bool
	baseUrl            string
	imageFileList      []string
	imagesDirectory    string
	linkDepth          int
	scrapeImages       bool
	scrapeUrls         bool
	urlFile            string
	urlList            []string
}

// NewScraper constructs a new Scraper. It includes users information about the configuration for web scraping and creates empty slices for URLs and image paths.
func NewScraper(baseUrl string, scrapeUrls bool, urlFile string, linkDepth int, allowUrlDuplicates bool, scrapeImages bool, imagesDirectory string) *Scraper {
	return &Scraper{
		allowUrlDuplicates: allowUrlDuplicates,
		baseUrl:            baseUrl,
		imageFileList:      make([]string, 0),
		imagesDirectory:    imagesDirectory,
		linkDepth:          linkDepth,
		scrapeImages:       scrapeImages,
		scrapeUrls:         scrapeUrls,
		urlFile:            urlFile,
		urlList:            make([]string, 0),
	}
}

// StartScrape inits scraping for URLs and images. Also, if `allowUrlDuplicates` is set to `false`, it cleans the list of found URLs.
func (s *Scraper) StartScrape() error {
	err := s.startUrlScrape()
	if err != nil {
		return err
	}
	if !s.allowUrlDuplicates {
		fmt.Printf("Original len:\t%d\n", len(s.urlList))
		s.urlList = utility.RemoveDuplicateStr(s.urlList)
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
	if len(s.urlList) > 0 {
		for _, url := range s.urlList {
			if !strings.HasPrefix(url, "https") || !strings.HasPrefix(url, s.baseUrl) {
				continue
			}
			results, err := s.scrapeImage(url)
			if err != nil {
				return err
			}
			s.imageFileList = append(s.imageFileList, results...)
		}
	} else {
		return errors.New("cannot scrape images from empty url list")
	}
	return nil
}

// SaveResults saves found URLs and images to files.
func (s *Scraper) SaveResults() {
	if s.scrapeUrls {
		file, err := os.OpenFile(s.urlFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		writer := bufio.NewWriter(file)
		for _, url := range s.urlList {
			_, err := writer.WriteString(url + "\n")
			if err != nil {
				panic(err)
			}
		}
		writer.Flush()
	}
	if s.scrapeImages {
		for _, url := range s.imageFileList {
			fileEnding := utility.ExtractProbableFileEnding(url)
			file, err := os.OpenFile(s.imagesDirectory+utility.HashUrlToFileName(url)+"."+fileEnding, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			resp, err := http.Get(url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				panic(errors.New("received non 200 response code"))
			}
			//Write the bytes to the file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				panic(err)
			}
		}
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
			} else if url != "/" && url != "" && !strings.HasPrefix(url, "#") {
				res := s.baseUrl + url
				if !slices.Contains(s.urlList, res) && !slices.Contains(results, res) {
					results = append(results, url)
				}
			}
		}
	})
	return results, nil
}

func (s *Scraper) scrapeImage(baseUrl string) ([]string, error) {
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
		return nil, errors.New("status code error when initializing image scrape")
	}
	// Read and parse response data
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	doc.Find("img").Each(func(i int, selector *goquery.Selection) {
		// for each item found, get the url
		url, found := selector.Attr("src")
		if found {
			if strings.HasPrefix(url, "/") && url != "/" {
				res := s.baseUrl + url
				if !slices.Contains(s.imageFileList, res) && !slices.Contains(results, res) {
					results = append(results, res)
				}
			} else if url != "/" && url != "" && !strings.HasPrefix(url, "#") {
				res := s.baseUrl + url
				if !slices.Contains(s.imageFileList, res) && !slices.Contains(results, res) {
					results = append(results, url)
				}
			}
		}
	})
	return results, nil
}

// TODO: refactor scrapeUrl and scrapeImage into this function
// func (s *Scraper) scrape(baseUrl string, container string, find string, attr string) ([]string, error) {
//   return nil, nil
// }
