package scraper

import (
	"github.com/gocolly/colly"
	"log"
	"mergen/internal/post"
	"strconv"
	"strings"
	"time"
)

func scrapeEksiEntries(result chan []post.Post) {
	const baseUrl = "https://eksisozluk.com"
	const timeLayout = "02.01.2006 15:04"
	const source = "Eksisozluk"
	var posts []post.Post
	c := colly.NewCollector(
		colly.AllowedDomains("eksisozluk.com"),
	)

	// parse popular topics
	c.OnHTML("ul[class]", func(e *colly.HTMLElement) {
		if e.Attr("class") != "topic-list" {
			return
		}
		log.Println("[Scraper:eksi] Getting topics from eksisozluk#haber")

		topicPaths := e.ChildAttrs("a", "href")
		if len(topicPaths) == 0 {
			log.Println("[Scraper:eksi] Couldn't parse any topics")
			return
		}

		for _, path := range topicPaths {
			c.Visit(baseUrl + path)
		}
	})

	// parse the first entry in the topic page
	c.OnHTML("div[id]", func(e *colly.HTMLElement) {
		if e.Attr("id") != "topic" {
			return
		}
		topicTitle := e.DOM.Find("span[itemprop]").Text()
		entry := e.DOM.Find("li[data-id]").First()

		author, _ := entry.Attr("data-author")
		text := entry.Find("div .content").First().Text()
		text = strings.TrimSuffix(strings.TrimPrefix(text, "\n    "), "\n")

		favString, _ := entry.Attr("data-favorite-count")
		score, err := strconv.ParseInt(favString, 10, 64)

		dateStr := e.DOM.Find("a[class='entry-date permalink']").First().Text()
		urlPath := e.ChildAttr("a[class='entry-date permalink']", "href")

		url := baseUrl + urlPath

		if strings.Contains(dateStr, "~") {
			dateStr = strings.Split(dateStr, " ~")[0]
		}
		if err != nil {
			log.Println("[Scraper:eksi] Couldn't parse entry score")
			return
		}
		entryTime, _ := time.Parse(timeLayout, dateStr)
		entryTime = entryTime.Add(time.Duration(-3) * time.Hour)

		p := post.Post{
			Title:     topicTitle,
			Source:    source,
			Author:    author,
			Text:      text,
			Url:       url,
			Timestamp: entryTime.Unix(),
			Score:     score,
			Language:  "tr",
		}

		posts = append(posts, p)
	})

	c.Visit(baseUrl + "/basliklar/kanal/haber")
	result <- posts
}
