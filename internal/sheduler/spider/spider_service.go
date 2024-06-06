package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
)

type BingSpider struct {
	collector *colly.Collector
}

func NewBingSpider(options ...colly.CollectorOption) Spider {
	return &BingSpider{
		collector: colly.NewCollector(options...),
	}
}

func (spider *BingSpider) CrawlImages(keyword string) []string {
	var images []string
	spider.collector.OnHTML("div.img_cont.hoff", func(e *colly.HTMLElement) {
		e.DOM.Find("img").Each(func(_ int, s *goquery.Selection) {
			src, exists := s.Attr("src")
			if exists {
				images = append(images, src)
			}
		})
	})

	query := url.QueryEscape(keyword)
	urls := fmt.Sprintf("https://cn.bing.com/images/search?q=%s&first=1&cw=1897&ch=1005", query)
	spider.collector.Visit(urls)

	return images
}

func (spider *BingSpider) CrawlArticles(keyword string) []string {

	return nil
}
