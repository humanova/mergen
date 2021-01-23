package main

import (
	"log"
	"mergen/internal/post"
	"mergen/internal/scraper"
	"time"
)

func main() {
	err := post.InitDb()
	if err != nil {
		panic(err)
	}

	// test post functions
	p1 := post.Post{Title:"Title", Author:"Author", Text:"blabla", Url:"url_blabla", Timestamp:time.Now().UTC().Unix(), Score:100}
	p2 := post.Post{Title:"Title2", Author:"Author2", Text:"blabla2", Url:"url_blabla2", Timestamp:time.Now().UTC().Unix(), Score:101}

	err = post.Add(p1)
	if err != nil {
		log.Fatal(err)
	}

	err = post.AddAll([]post.Post{p1, p2})
	if err != nil {
		log.Fatal(err)
	}

	err = post.UpdateScore(p2, 420)
	if err != nil {
		log.Fatal(err)
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