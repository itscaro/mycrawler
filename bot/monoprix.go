package main

import (
	"net/http"
	"fmt"
	"time"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/PuerkitoBio/fetchbot"
)

func processMonoprix() {
	f := fetchbot.New(fetchbot.HandlerFunc(monoprixHandler))
	f.UserAgent = "Mozilla/5.0 (iPad; U; CPU OS 3_2_1 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Mobile/7B405"
	f.CrawlDelay = 100 * time.Microsecond

	queue := f.Start()
	queue.SendStringGet(
		"https://www.monoprix.fr/mangue-kent-null-2470229-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-vittel-2409711-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-evian-1068-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-contrex-2410891-p",
	)
	queue.Close()
}

func monoprixHandler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	productName1 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h1").Text(), "\t\n ")
	productName2 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h2").Text(), "\t\n ")
	productName3 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h4").Text(), "\t\n ")
	productPrice := strings.Trim(doc.Find("#priceChange").Text(), "\t\n ")
	productPromo := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > div.description.courses > span").Text(), "*\t\n ")

	productName := productName1 + " (" + productName2 + ") " + productName3
	//log.Printf("Visited: %s\n", ctx.URL().String())
	fmt.Printf("%s %s (%s)\n", productName, productPrice, productPromo)
	//log.Printf("%s", productNameAll)

	fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
}
