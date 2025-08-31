# Go web crawler
This is a guided projected by [boot.dev](https://www.boot.dev/). 
Simple web crawler that just visits a site, crawls the pages and prints out a report of the number of times a page was referenced by internal links in descending order.

To run the web crawler pass the website url to crawl like so:
```sh
go build -o crawler && ./crawler <website-url>
```
The crawler also takes 2 positional arguments after the url, maxConcurency and maxPages. The first one defines the max amount of goroutines to spawn and the second the max amount of pages to crawl on that website.
Note that the crawler crawl each page only once and for each visit after that just increments the time found counter.