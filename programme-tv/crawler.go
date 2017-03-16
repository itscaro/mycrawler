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

var rxOk = regexp.MustCompile(`http://(www\.)?programme-tv\.net(.*)`)
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
	productName := strings.Trim(doc.Find("#content > div > div > div.span4 > aside > h1").Text(), "\t\n ")
	productPrice := strings.Trim(doc.Find("#priceChange").Text(), "\t\n ")

	log.Printf("Visited: %s: %s - %s\n", ctx.URL().String(), productName, productPrice)

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

	// Play nice with ddgo when running the test!
	opts.MaxVisits = 0

	// Create crawler and start at root of duckduckgo
	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run([]string{"http://www.programme-tv.net/", "http://www.programme-tv.net/videos"})
}
