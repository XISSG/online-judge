package scheduler

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// 并发爬取
func ConcurrencyCrawl(totalPage int, savePath string, concurrency int) {
	start := time.Now()
	var wg sync.WaitGroup

	ch := make(chan struct{}, concurrency)
	for i := 1; i <= totalPage; i++ {
		wg.Add(1)
		ch <- struct{}{}
		page := i
		go func(page int) {
			defer wg.Done()
			defer func() { <-ch }()
			Crawl(page, savePath)
		}(page)
	}

	wg.Wait()
	close(ch)
	cost := time.Since(start)
	fmt.Printf("total cost time: %v\n", cost)
}

// 单线程爬取
func Crawl(page int, savePath string) {
	//爬取网页并解析网页
	images := CrawlWeb(page)

	//下载图片
	start := time.Now()
	DownloadAvatar(images, savePath)
	duration := time.Since(start)
	fmt.Printf("page %v , cost time: %v\n", page, duration)
}

// 该函数需根据实际网站进行修改
func CrawlWeb(page int) []string {
	// 创建一个上下文和取消函数
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 登录页面的URL和目标页面的URL
	url := fmt.Sprintf("https://www.duitang.com/category/?cat=avatar#!hot-p%v", page)

	// 定义需要执行的任务
	var pageContent string
	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div.mask-close`), // 等待目标页面的某个元素出现
		chromedp.Click(`div.mask-close`),       // 点击登录按钮
		chromedp.Sleep(time.Second * 2),
		chromedp.OuterHTML(`html`, &pageContent), // 获取整个页面内容
	}

	// 执行任务
	if err := chromedp.Run(ctx, tasks); err != nil {
		log.Fatal(err)
	}

	imageUrls := parseDom(pageContent)
	return imageUrls
}

// 需根据实际网页进行对应解析
func parseDom(pageContent string) []string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(pageContent))
	if err != nil {
		log.Fatalln(err)
	}

	var images []string
	dom.Find("div.mbpho img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			images = append(images, src)
		}
	})
	return images
}

// 访问图片网址并保存
func DownloadAvatar(urls []string, savePath string) {
	for _, url := range urls {
		url = strings.TrimSuffix(url, "_webp")
		err := downloadFile(url, savePath)
		if err != nil {
			continue
		}
	}
}

func downloadFile(url string, savePath string) error {
	fileName := path.Base(url)
	filePath := filepath.Join(savePath, fileName)

	resp, _ := http.Get(url)

	err := os.MkdirAll(savePath, 0755)
	fd, err := os.Create(filePath)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	defer fd.Close()

	reader := bufio.NewReader(resp.Body)
	body, _ := io.ReadAll(reader)
	_, err = fd.Write(body)
	return err
}
