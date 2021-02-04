package handler

import (
	"log"
	"mergen/internal/post"
	"net/http"
	"html/template"
	"sort"
	"time"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	var posts []post.Post

	tpl, err := template.ParseGlob("web/template/*.gohtml")
	if err != nil {
		log.Println("[PostsHandler] couldn't parse templates %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	timestamp := (time.Now().Add(time.Duration(-12) * time.Hour)).UTC().Unix()
	posts, err = post.GetPosts(timestamp)
	if err != nil {
		log.Println("[PostsHandler] couldn't get posts from DB : %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}
	// sort posts by their date (new to old)
	sort.Slice(posts, func(i, j int) bool {return posts[i].Timestamp > posts[j].Timestamp})

	tpl.ExecuteTemplate(w, "posts.gohtml", posts)
	return
}