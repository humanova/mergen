package handler

import (
	"log"
	"mergen/internal/post"
	"net/http"
	"time"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	var posts []post.Post

	timestamp := (time.Now().Add(time.Duration(-12) * time.Hour)).UTC().Unix()
	posts, err := post.GetPosts(timestamp)
	if err != nil {
		log.Println("[PostsHandler] couldn't get posts from DB : %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}
	log.Printf("%v", posts)
	return
}