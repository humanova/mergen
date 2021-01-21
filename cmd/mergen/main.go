package main

import (
	"log"
	"mergen/internal/post"
	"time"
)

func main() {
	err := post.InitDb()
	if err != nil {
		panic(err)
	}

	// test adding a new post to db
	p1 := post.Post{"Title", "Author", "blabla", "url_blabla", time.Now().UTC().Unix(), 100}
	p2 := post.Post{"Title2", "Author2", "blabla2", "url_blabla2", time.Now().UTC().Unix(), 101}

	err = post.AddPosts([]post.Post{p1, p2})
	if err != nil {
		log.Fatal(err)
	}
}