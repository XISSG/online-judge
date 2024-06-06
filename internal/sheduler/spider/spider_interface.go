package spider

type Spider interface {
	CrawlImages(keyword string) []string
	CrawlArticles(keyword string) []string
}
