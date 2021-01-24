package main

import (
	"log"
	"mergen/internal/post"
	"mergen/internal/scraper"
)

func main() {
	err := post.InitDb()
	if err != nil {
		panic(err)
	}

	//test scraper funcs
	var posts []post.Post
	posts, err = scraper.ScrapeNews()
	if err != nil {
		log.Fatal(err)
	}

	err = post.AddAll(posts)
	if err != nil {
		log.Fatal(err)
	}

}