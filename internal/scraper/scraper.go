package scraper

import (
	"log"
	"mergen/internal/post"
)

func ScrapeAll() {
	var posts []post.Post

	newsChan := make(chan []post.Post)
	eksiChan := make(chan []post.Post)
	twitterChan := make(chan []post.Post)
	redditChan := make(chan []post.Post)

	go scrapeNews(newsChan)
	go scrapeEksiEntries(eksiChan)
	go scrapeTwitter(twitterChan)
	go scrapeReddit(redditChan)

	newsPosts := <- newsChan
	eksiPosts := <- eksiChan
	twitterPosts := <- twitterChan
	redditPosts := <- redditChan

	for _, postCollection := range [][]post.Post{newsPosts, eksiPosts, twitterPosts, redditPosts} {
		posts = append(posts, postCollection...)
	}

	err := post.AddAll(posts)
	if err != nil {
		log.Printf("[Scraper:main] error while inserting posts to db : %s\n", err)
	}

	log.Printf("[Scraper:main] Scraped %d from rss feeds, %d from eksisozluk, %d from twitter, " +
	"%d from reddit\n------\n", len(newsPosts), len(eksiPosts), len(twitterPosts), len(redditPosts))
}
