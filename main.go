package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"log"
	"sync"
	"time"
)

type Article struct {
	Title string
	Content string
}

func (a *Article) putContent(content string) {
	a.Content = fmt.Sprint(a.Content, content)
}

func checkErr(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}

var nodeNameMap = map[string]string{
	"p"		: "<p>%s</p>",
	"h1"    : "<h3>%s</h3>",
	"div"	: "<div>%s</div>",
	"pre"	: "<pre><code>%s</code></pre>",
	"code"	: "<code>%s</code>",
	"figure": "<center>%s</center>",
}

func main()  {
	var mutex sync.Mutex
	url := "https://medium.com/flutter/the-power-of-webviews-in-flutter-a56234b57df2"
	c := colly.NewCollector(colly.AllowURLRevisit())
	rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1080")
	checkErr(err)

	c.SetProxyFunc(rp)

	c.OnHTML("article section>div", func(e *colly.HTMLElement) {
		article := new(Article)
		e.DOM.Children().Each(func(i int, selection *goquery.Selection) {
			var content string
			nodeName := goquery.NodeName(selection)
			if nodeName == "p" || nodeName == "div" || nodeName == "h1" {
				mutex.Lock()
				if i == 0 {
					content = selection.Find("h1").Text()
				} else {
					content = selection.Text()
				}
				response := Translate(content)
				if len(response.Result) > 0 {
					content = response.Result[0].Dst
				}
				if i == 0 {
					article.Title = content
				} else {
					content = fmt.Sprintf(nodeNameMap[nodeName], content)
					article.putContent(content)
				}
				time.Sleep(time.Second)
				mutex.Unlock()
			} else {
				content = fmt.Sprintf(nodeNameMap[nodeName], selection.Text())
				article.putContent(content)
			}
		})
		fmt.Println(article.Content)
	})

	err = c.Visit(url)
	checkErr(err)
}
