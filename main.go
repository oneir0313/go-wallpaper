package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gocolly/colly"
	"github.com/robfig/cron/v3"
)

var (
	fileName    string
	fullURLFile string
)

const picPath = "c:/Users/lennychang/Pictures/Saved Pictures/"

func main() {
	cr := cron.New()
	c := colly.NewCollector()
	_, err := cr.AddFunc("0, 0, 1, *, *", func() {
		removeContents(picPath)
		c.OnHTML(".media-download > a", func(e *colly.HTMLElement) {
			monthStr := time.Now().Month().String()
			imageHref := e.Attr("href")

			if strings.Contains(strings.ToLower(imageHref), "2022_"+strings.ToLower(monthStr)) && e.Text == "1920 x 1080" {
				fullURLFile = "https://www.taiwan.net.tw" + imageHref

				// Build fileName from fullPath
				fileURL, err := url.Parse(fullURLFile)
				if err != nil {
					log.Fatal(err)
				}
				path := fileURL.Path
				segments := strings.Split(path, "/")
				fileName = picPath + segments[len(segments)-1]

				// Create blank file
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
				// Put content on file
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
		err := c.Visit("https://www.taiwan.net.tw/m1.aspx?sNo=0012076")
		if err != nil {
			fmt.Printf("Visit web url failed: %s\n", err.Error())
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	go cr.Start()

	fmt.Println("Cron job start")

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
}

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
