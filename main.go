package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// BUG: The description tag is partially implemented but most no tech orientened feeds opt to use
// content:encoded tag and the <![CDATA[ and ]]> wrappers to workaround the limitatiopns of rss feed
// format.

type Feed struct {
	Title       string
	Description string
	Link        string
	Date        string
}

var urls = []string{
	"https://rumble.com/c/",
	"https://rumble.com/user/",
}

// Build date and all other date/times supplied in RFC822 format.
var buildDate = time.Now().Format(time.RFC822)

func parseURLS(urls []string, channelName string) *http.Response {
	var rtn *http.Response
	for _, url := range urls {
		resp, err := http.Get(url + channelName)
		if resp.StatusCode == 200 && err == nil {
			rtn = resp
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	return rtn
}

func scrapeHTML() {
	// Request the HTML page.

	// TODO: Add arguemnt parsing to allow of multiple sites and filtering of the feed.
	channelName := os.Args[len(os.Args)-1]
	if channelName == "" {
		panic("Please provide a channel name.")
	}

	resp := parseURLS(urls, channelName)

	// data, err := ioutil.ReadFile("index.html")
	// For testing local copy of the pages HTML can be sourcced to test your code.
	// Without the need to make a request to the server.

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Load the HTML document.
	// doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(data)))
	// Uncomment the above line and comment out the below line to test with local copy of the page.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	feed := []Feed{}

	// Find the review items.
	doc.Find(".videostream").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title, link and description.
		title := s.Find("h3").Text()
		link := s.Find("a").AttrOr("href", "")
		duration := s.Find(".videostream__status--duration").Text()
		date := s.Find("time").AttrOr("datetime", "")

		parseDate, _ := time.Parse(time.RFC3339, date)

		feed = append(feed, Feed{
			Title:       strings.TrimSpace(title),
			Description: "[Duration: " + strings.TrimSpace(duration) + "]",
			// BUG: RSS requires <pubDate> to be in RFC822 format.  Date supplied from the page is in RFC3339 format.
			Date: parseDate.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
			// NOTE: Leave off trailing slash link is prefixed with slash during the scraping.
			Link: "https://rumble.com" + link,
		})

	})

	var rssTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
<channel>
	<title>` + channelName + ` RSS Feed</title>
	<link>https://rumble.com/user/` + channelName + `</link>
	<description>Golang goquery rss feed.</description>
	<language>en-us</language>
	<lastBuildDate>` + buildDate + `</lastBuildDate>
		{{- range . }}
		<item>
			<title>{{ .Title }}</title>
			<link>{{ .Link }}</link>
			<pubDate>{{ .Date }}</pubDate>
			<guid>{{ .Link }}</guid>
			<description>{{ .Description }}</description>
		</item>
	{{- end }}
</channel>
</rss>`

	tmpl, err := template.New(rssTmpl).Parse(rssTmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, feed)
	if err != nil {
		panic(err)
	}

}

func main() {
	scrapeHTML()
}
