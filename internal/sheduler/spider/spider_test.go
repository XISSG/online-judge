package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"testing"
)

func TestSpider(t *testing.T) {
	//crawl := NewBingSpider(
	//	colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0"),
	//)
	crawl := NewBingSpider(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0"),
		colly.AllowedDomains("bing.com"),
	)
	images := crawl.CrawlImages("原神")
	fmt.Println(images)
}
