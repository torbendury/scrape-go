# scrape-go

A tiny scraper built with Golang. Work in progress!

I'm very aware of [gocolly/colly](https://github.com/gocolly/colly) but would love to experiment with scraping and learning a bit more than by using a framework directly from the start.

See [TODO](TODO).

## License

See [LICENSE](LICENSE).

## Usage

```bash
$ scrape-go -h

Usage of scrape-go:
  -allow-duplicate-urls
        Allow duplicated links. Only takes effect if URL scraping is active.
  -base-url string
        The base URL to start with.
  -images
        NOT IMPLEMENTED YET: Set to true if you want to scrape images.
  -link-depth int
        Maximum scraping depth. (default 5)
  -url-outfile string
        The file to write scraped URLs to. (default "urls.txt")
  -urls
        Set to true if you want to scrape deep link URLs.
```

## Try it

```bash
go run cmd/scrape/main.go -urls -base-url=https://torbentechblog.com
```
