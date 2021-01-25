package main

import (
	"github.com/go-co-op/gocron"
	"log"
	"mergen/internal/post"
	"mergen/internal/scraper"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := post.InitDb()
	if err != nil {
		panic(err)
	}
	log.Println("Starting scraper cron job...")

	scraperCron := gocron.NewScheduler(time.UTC)
	_, err = scraperCron.Every(5).Minutes().Do(scraper.ScrapeAll)
	if err != nil {
		log.Fatalf("couldn't create cron job : %s", err)
	}
	scraperCron.StartAsync()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
