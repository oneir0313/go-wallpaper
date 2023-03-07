package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

var (
	fileName    string
	fullURLFile string
	run bool
)

func init() {
  flag.BoolVar(&run, "run", false, "run one time")
}


func main() {
	flag.Parse()

	picPath := "./pictures/"

	cr := cron.New()
	c := colly.NewCollector()

	crawlPic := func() {
		// 刪除上一個月份檔案
		err := removeContents(picPath)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// 抓取觀光局網站上的下載錨點元素
		c.OnHTML(".media-download > a", func(e *colly.HTMLElement) {
			imageHref := e.Attr("href")
			yearStr := strconv.Itoa(time.Now().Year())
			monthStr := time.Now().Month().String()

			// 判斷錨點元素內的年份月份有相符且符合1920x1080的字串
			if strings.Contains(strings.ToLower(imageHref), yearStr+"_"+strings.ToLower(monthStr)) && e.Text == "1920 x 1080" {
				fullURLFile = "https://www.taiwan.net.tw" + imageHref

				// 建立檔案名稱
				fileURL, err := url.Parse(fullURLFile)
				if err != nil {
					log.Fatal(err)
				}
				path := fileURL.Path
				segments := strings.Split(path, "/")
				fileName = picPath + segments[len(segments)-1]

				// 建立空白檔案
				file, err := os.Create(fileName)
				if err != nil {
					log.Fatal(err)
				}
				client := http.Client{
					CheckRedirect: func(r *http.Request, via []*http.Request) error {
						r.URL.Opaque = r.URL.Path
						return nil
					},
				}
				// 將內容物下載至檔案
				resp, err := client.Get(fullURLFile)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				size, _ := io.Copy(file, resp.Body)

				defer file.Close()

				fmt.Printf("Downloaded a file %s with size %d", fileName, size)
			}
		})
		// Colly訪問台灣觀光局網站
		err = c.Visit("https://www.taiwan.net.tw/m1.aspx?sNo=0012076")
		if err != nil {
			fmt.Printf("Visit web url failed: %s\n", err.Error())
		}
	}

	if run {
		crawlPic()
		return
	}
	// 執行排程每個月一號零點整執行下載
	_, err := cr.AddFunc("0, 0, 1, *, *", crawlPic)
	if err != nil {
		log.Fatal(err)
	}

	// goroutine執行cronjob
	go cr.Start()

	fmt.Println("Cron job start")

	// 直到收到終止訊號保持服務執行
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
}

// 刪除目錄上的所有檔案
func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}