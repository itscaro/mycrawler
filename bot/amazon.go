package main

import (
	"github.com/PuerkitoBio/fetchbot"
	"time"
	"net/http"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func processAmazon() {
	f := fetchbot.New(fetchbot.HandlerFunc(amazonHandler))
	f.UserAgent = "Mozilla/5.0 (iPad; U; CPU OS 3_2_1 like Mac OS X; en-us) AppleWebKit/531.21.10 (KHTML, like Gecko) Mobile/7B405"
	f.CrawlDelay = 100 * time.Microsecond

	queue := f.Start()
	queue.SendStringGet(
		"https://www.amazon.fr/Faber-Castell-167100-Feutre-PITT-artist/dp/B000TKEZDO/",
	)
	queue.Close()
}

func amazonHandler(ctx *fetchbot.Context, res *http.Response, err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}

	productName := strings.Trim(doc.Find("#productTitle").Text(), "\t\n ")
	productBrand := strings.Trim(doc.Find("#brand").Text(), "\t\n ")
	productPrice := strings.Trim(doc.Find("#priceblock_ourprice").Text(), "\t\n ")
	productPriceRefurbished := strings.Trim(doc.Find("#olp_feature_div > div > span:nth-child(2) > span").Text(), "\t\n ")

	//log.Printf("Visited: %s\n", ctx.URL().String())
	fmt.Printf("%s (%s) %s - %s\n", productName, productBrand, productPrice, productPriceRefurbished)
	//log.Printf("%s", productNameAll)

	fmt.Printf("[%d] %s %s\n", res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL())
}
