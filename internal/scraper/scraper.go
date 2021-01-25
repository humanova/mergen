package scraper

import (
	"log"
	"mergen/internal/post"
)

func ScrapeAll() {
	var eksiPosts []post.Post
	var newsPosts []post.Post

	newsPosts, err := scrapeNews()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeNews(): %s\n", err)
	}

	eksiPosts, err = scrapeEksiEntries()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeEksiEntries(): %s\n", err)
	}

	err = post.AddAll(append(newsPosts, eksiPosts...))
	if err != nil {
		log.Printf("[Scraper:main] error while inserting posts to db : %s\n", err)
	}
	log.Printf("[Scraper:main] Scraped %d from rss feeds, "+
		"%d from eksisozluk\n------\n", len(newsPosts), len(eksiPosts))
}
