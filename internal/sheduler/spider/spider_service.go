package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"net/url"
)

type Spider interface {
	CrawlPictures(keyword string) []string
	CrawlArticles(keyword string) []string
}
type spider struct {
	collector *colly.Collector
}

func NewBingSpider(options ...colly.CollectorOption) Spider {
	return &spider{
		collector: colly.NewCollector(options...),
	}
}

func (spider *spider) CrawlPictures(keyword string) []string {
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

func (spider *spider) CrawlArticles(keyword string) []string {

	return nil
}
