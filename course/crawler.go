package main

import (
	"github.com/PuerkitoBio/gocrawl"
	"time"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"log"
	"strings"
)

var rxOk = regexp.MustCompile(`https?://(www\.)?monoprix\.fr/([^/]+)-([0-9]+)(-p)?$`)
//var rxOk = regexp.MustCompile(`https?://(www\.)?monoprix\.fr/([^/]*)-([0-9]+)-p?$`)

// Create the Extender implementation, based on the gocrawl-provided DefaultExtender,
// because we don't want/need to override all methods.
type ExampleExtender struct {
	gocrawl.DefaultExtender // Will use the default implementation of all but Visit and Filter
}

// Override Visit for our need.
func (x *ExampleExtender) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	// Use the goquery document or res.Body to manipulate the data
	// ...
	//productNameAll := strings.Trim(doc.Find("#content > div > div > div.span4 > aside").Text(), "\t\r\n ")
	productName1 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h1").Text(), "\t\n ")
	productName2 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h2").Text(), "\t\n ")
	productName3 := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h4").Text(), "\t\n ")
	productPrice := strings.Trim(doc.Find("#priceChange").Text(), "\t\n ")
	productPromo := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > div.description.courses > span").Text(), "*\t\n ")

	productName := productName1 + " (" + productName2 + ") " + productName3
	//log.Printf("Visited: %s\n", ctx.URL().String())
	log.Printf("%s %s (%s)\n", productName, productPrice, productPromo)
	//log.Printf("%s", productNameAll)

	// Return nil and true - let gocrawl find the links
	return nil, true
}

// Override Filter for our need.
func (x *ExampleExtender) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	return !isVisited && rxOk.MatchString(ctx.NormalizedURL().String())
}

func main() {
	// Set custom options
	opts := gocrawl.NewOptions(new(ExampleExtender))

	// should always set your robot name so that it looks for the most
	// specific rules possible in robots.txt.
	opts.RobotUserAgent = "Example"
	// and reflect that in the user-agent string used to make requests,
	// ideally with a link so site owners can contact you if there's an issue
	opts.UserAgent = "Mozilla/5.0 (compatible; Example/1.0; +http://example.com)"

	opts.CrawlDelay = 100 * time.Microsecond
	opts.LogFlags = gocrawl.LogError

	seeds := []string{
		"https://www.monoprix.fr/mangue-kent-null-2470229-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-vittel-2409711-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-evian-1068-p",
		"https://www.monoprix.fr/eau-minerale-naturelle-contrex-2410891-p",
	}

	// Play nice with ddgo when running the test!
	opts.MaxVisits = len(seeds) - 1

	// Create crawler and start at root of duckduckgo
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run(seeds)
}
