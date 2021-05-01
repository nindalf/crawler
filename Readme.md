# Crawler

## Problem statement

> Given a starting URL, the crawler should visit each URL it finds on the same domain. It should print each URL visited, and a list of links found on that page. The crawler should be limited to one subdomain - so when you start with *https://example.com/*, it would crawl all pages on the example.com website, but not follow external links, for example to google.com or github.com.

## Usage

```
➜ go build
➜ ./crawler --help
Usage of ./crawler:
  -u string
    	The URL to be crawled (default "https://example.com")
  -w int
    	The number of workers to run simultaneously (default 4)
```

## Example

```
➜ ./crawler -u https://blog.nindalf.com -w 25
2021/05/01 10:16:10 Crawl starting
2021/05/01 10:16:12 Crawl complete
2021/05/01 10:16:12 Found 207 urls
[https://pages.cloudflare.com https://store.steampowered.com/app/427520/Factorio/ ...
```
