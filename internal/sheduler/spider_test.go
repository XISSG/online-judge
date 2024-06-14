package scheduler

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSpider(t *testing.T) {
	curPath, _ := os.Getwd()
	path := filepath.Join(curPath, "../../public/pictures/avatar")
	start := time.Now()
	//for i := 1; i < 30; i++ {
	//	Crawl(i, path)
	//}
	ConcurrencyCrawl(30, path, 5)
	cost := time.Since(start)
	fmt.Println("total cost time: ", cost)

}
