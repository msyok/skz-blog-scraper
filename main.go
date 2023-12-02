package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/msyok/skz-blog-scraper/lib"
	"golang.org/x/net/html"
)

const hostname = "https://sakurazaka46.com"
const blogListPath = "/s/s46/diary/blog/list?ima=0000"

func main() {
	var id int
	var c int
	var dryRun bool
	var postPath string
	flag.IntVar(&id, "id", 0, "member id")
	flag.IntVar(&c, "c", 3, "concurrency")
	flag.BoolVar(&dryRun, "dry", false, "dry run")
	flag.StringVar(&postPath, "path", "", "post path")

	flag.Parse()

	if id == 0 {
		log.Fatal("set `id` argument")
	}

	saver := lib.NewSaver(id)
	if err := os.MkdirAll(saver.GetPostDir(), 0755); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	if err := os.MkdirAll(saver.GetImageDir(), 0755); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	if postPath != "" {
		srapingPost(postPath, saver)
		return
	}

	page := 0
	for {
		if !sraping(id, page, saver, c, dryRun) || dryRun {
			log.Println("done")
			return
		}
		page++
	}
}

func getPostIdFromUrl(postUrl string) string {
	lastSlachIdx := strings.LastIndex(postUrl, "/")
	queryParemeterIdx := strings.Index(postUrl, "?")
	return postUrl[lastSlachIdx+1 : queryParemeterIdx]
}

func sraping(id, page int, saver *lib.Saver, concurrency int, dryRun bool) bool {
	processing := false
	blogPageUrl := fmt.Sprintf("%s%s&ct=%d&page=%d", hostname, blogListPath, id, page)

	c := colly.NewCollector()

	c.OnHTML("ul.com-blog-part.box3.fxpc", func(e *colly.HTMLElement) {
		postPathList := make([]string, 0)
		e.ForEach("li.box", func(i int, e *colly.HTMLElement) {
			postPathList = append(postPathList, e.ChildAttr("a", "href"))
		})
		if len(postPathList) == 0 {
			processing = false
			return
		}

		processing = true

		if dryRun {
			postPathList = postPathList[:1]
		}

		var wg sync.WaitGroup
		guard := make(chan struct{}, concurrency)
		for i := 0; i < len(postPathList); i += 1 {
			guard <- struct{}{}
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				srapingPost(postPathList[index], saver)
				<-guard
			}(i)
		}
		wg.Wait()
	})

	c.Visit(blogPageUrl)
	c.Wait()
	return processing
}

func srapingPost(postPath string, saver *lib.Saver) {
	log.Println("start post", postPath)

	postId := getPostIdFromUrl(postPath)

	var (
		year  string
		month string
		day   string
		title string
	)

	c := colly.NewCollector()

	c.OnHTML("div.ym-txt", func(e *colly.HTMLElement) {
		year = e.ChildText("span.ym-year")
		month = e.ChildText("span.ym-month")
	})
	c.OnHTML("div.ym-inner", func(e *colly.HTMLElement) {
		day = e.ChildText("p.date.wf-a")
	})
	c.OnHTML("h1.title", func(e *colly.HTMLElement) {
		title = lib.RemoveUnnecessaryTokens(e.Text)
	})
	c.OnHTML("div.box-article", func(e *colly.HTMLElement) {
		node := e.DOM.Nodes[0]
		content, imageUrls := parsePostContent(node, postId, saver)
		if content == "" {
			return
		}
		if err := saver.SavePost(postId, year, month, day, title, content); err != nil {
			log.Fatal(err)
		}
		for _, imageUrl := range imageUrls {
			if err := saver.SaveImage(postId, imageUrl); err != nil {
				log.Fatal(err)
			}
		}
	})

	c.Visit(fmt.Sprintf("%s%s", hostname, postPath))

	log.Println("end post", postPath)
}

func parsePostContent(node *html.Node, postId string, saver *lib.Saver) (string, []string) {
	w := lib.NewWriter()

	imageUrls := make([]string, 0)

	child := node.FirstChild
	for child != nil {
		data := child.Data
		if data == "a" {
			name, link := lib.GetLink(child)
			w.WriteLink(name, link)
		} else if data == "img" {
			name, url := lib.GetImage(child)
			if url != "" && lib.IsValidImageUrl(url) {
				w.WriteImage(name, saver.GetImagePath(postId, url))
				imageUrls = append(imageUrls, fmt.Sprintf("%s%s", hostname, url))
			}
		} else if data == "br" {
			w.WriterLineBreak()
		} else if data == "div" || data == "span" || data == "p" {
			childContent, childImageUrls := parsePostContent(child, postId, saver)
			childContent = lib.RemoveUnnecessaryTokens(childContent)
			w.WriteText(childContent)
			if ((childContent == "" && data == "div") || data == "p") {
				w.WriterLineBreak()
			}
			imageUrls = append(imageUrls, childImageUrls...)
		} else {
			data = lib.RemoveUnnecessaryTokens(data)
			w.WriteText(data)
		}

		child = child.NextSibling
	}

	return w.String(), imageUrls
}
