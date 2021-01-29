package scraper

import (
	"fmt"
	"github.com/turnage/graw/reddit"
	"log"
	"mergen/internal/post"
)

type announcer struct{}

func (a *announcer) Post(post *reddit.Post) error {
	fmt.Printf(`%s posted "%s"\n`, post.Author, post.Title)
	return nil
}

func ScrapeReddit() ([]post.Post, error) {
	var posts []post.Post

	bot, err := reddit.NewBotFromAgentFile("mergenbot.agent", 0)
	if err != nil {
		log.Println("[Scraper:reddit] Failed to create bot handle: ", err)
		return nil, err
	}

	harvest, err := bot.Listing("/r/wallstreetbets", "")
	if err != nil {
		log.Println("[Scraper:reddit] Failed to fetch /r/golang: ", err)
		return nil, err
	}

	for _, post := range harvest.Posts[:5] {
		log.Printf("[%s] posted [%s]\n", post.Author, post.Title)
	}

	return posts, nil
}
