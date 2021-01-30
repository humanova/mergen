package scraper

import (
	"log"
	"mergen/internal/post"
)

func ScrapeAll() {
	var eksiPosts []post.Post
	var newsPosts []post.Post
	var twitterPosts []post.Post
	var redditPosts []post.Post

	newsPosts, err := scrapeNews()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeNews(): %s\n", err)
	}

	eksiPosts, err = scrapeEksiEntries()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeEksiEntries(): %s\n", err)
	}

	twitterPosts, err = scrapeTwitter()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeTwitter(): %s\n", err)
	}

	redditPosts, err = scrapeReddit()
	if err != nil {
		log.Printf("[Scraper:main] error in scrapeReddit(): %s\n", err)
	}

	// mess below should be fixed :^) TODO: rewrite this
	err = post.AddAll(append(append(newsPosts, eksiPosts...), append(twitterPosts, redditPosts...)...))
	if err != nil {
		log.Printf("[Scraper:main] error while inserting posts to db : %s\n", err)
	}

	log.Printf("[Scraper:main] Scraped %d from rss feeds, %d from eksisozluk, %d from twitter, " +
	"%d from reddit\n------\n", len(newsPosts), len(eksiPosts), len(twitterPosts), len(redditPosts))

}
