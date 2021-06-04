package main

import (
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"log"
	"mergen/internal/handler"
	"mergen/internal/post"
	"mergen/internal/scraper"
	"net/http"
	"time"
)

func main() {
	err := post.InitDb()
	if err != nil {
		panic(err)
	}
	err = post.InitRedis()
	if err != nil {
		log.Println("couldn't connect to the redis server")
	}

	log.Println("Starting scraper cron job...")

	scraperCron := gocron.NewScheduler(time.UTC)
	_, err = scraperCron.Every(5).Minutes().Do(scraper.ScrapeAll)
	if err != nil {
		log.Fatalf("couldn't create cron job : %s", err)
	}
	scraperCron.StartAsync()

	r := mux.NewRouter()
	router := r.PathPrefix("/mergen").Subrouter()
	router.HandleFunc("/web", handler.WebHandler).Methods("GET")
	router.HandleFunc("/posts", handler.PostsHandler).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/mergen/static/", http.FileServer(http.Dir("./web/static/"))))

	err = http.ListenAndServe(":5005", router)
}
