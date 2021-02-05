package handler

import (
	"html/template"
	"log"
	"mergen/internal/post"
	"net/http"
	"sort"
)

var tpl *template.Template

func PostsHandler(w http.ResponseWriter, r *http.Request) {
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
		log.Println("[PostsHandler] couldn't get posts from DB : %s", err)
		http.Error(w, "500, internal error", http.StatusInternalServerError)
		return
	}
	// sort posts by timestamp
	sort.Slice(posts, func(i, j int) bool {return posts[i].Timestamp > posts[j].Timestamp})

	// if template isn't parsed yet, parse it
	if tpl == nil {
		tpl, err = template.ParseGlob("web/template/*.gohtml")
		if err != nil {
			log.Println("[PostsHandler] couldn't parse templates %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}

	tpl.ExecuteTemplate(w, "posts.gohtml", posts)
	return
}