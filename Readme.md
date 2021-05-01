# Crawler

![Build passing](https://github.com/nindalf/crawler/actions/workflows/go.yml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/nindalf/crawler)](https://goreportcard.com/report/github.com/nindalf/crawler) 
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/nindalf/crawler.svg)](https://github.com/nindalf/crawler)


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

## Note on implementation

This runs on one node, using in-memory queues and storage. If it needs to run on multiple nodes, it can be extended. Implement the `Queue` and `Storage` interfaces with RabbitMQ/Postgres/Redis and use these backends in main.go. `Queue` needs to be implemented in a thread-safe way but `Storage` does not.
