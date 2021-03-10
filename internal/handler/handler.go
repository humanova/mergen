package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"mergen/internal/post"
	"net/http"
	"sort"
)

var tpl *template.Template

func getPosts(r *http.Request) ([]post.Post, error) {
	var posts []post.Post
	var sources []string
	query := r.URL.Query()

	created_after := query.Get("created_after")
	created_before := query.Get("created_before")
	query_text := query.Get("query")
	sources = query["source"]  // can be an array
	author := query.Get("author")

	filters := post.Filters{created_after, created_before, query_text, sources, author}
	posts, err := post.GetPostsFiltered(filters)
	if err != nil {
		return nil, err
	}
	// sort posts by timestamp
	sort.Slice(posts, func(i, j int) bool {return posts[i].Timestamp > posts[j].Timestamp})

	return posts, nil
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts(r)
	if err != nil {
		log.Println("[PostsHandler] couldn't get posts from DB : %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(posts)
	if err != nil {
		log.Println("[PostsHandler] couldn't marshal posts: %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
	return
}

func WebHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := getPosts(r)
	if err != nil {
		log.Println("[WebHandler] couldn't get posts from DB : %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}

	// if template isn't parsed yet, parse it
	if tpl == nil {
		tpl, err = template.ParseGlob("web/template/*.gohtml")
		if err != nil {
			log.Println("[WebHandler] couldn't parse templates %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(w, "posts.gohtml", posts)
	return
}