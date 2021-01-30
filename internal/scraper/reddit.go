package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/turnage/graw/reddit"
	"io/ioutil"
	"log"
	"mergen/internal/post"
	"os"
)

type Subreddits struct {
	Subreddits []string `json:"subreddits"`
}

var bot reddit.Bot
var subList Subreddits

func getSubredditList(path string) error {
	subsFile, err := os.Open(path)
	if err != nil {
		log.Println("[Scraper:reddit] Couldn't open subreddit list file.")
		return err
	}
	defer subsFile.Close()

	byteValue, _ := ioutil.ReadAll(subsFile)

	err = json.Unmarshal(byteValue, &subList)
	if err != nil {
		log.Println("[Scraper:reddit] Couldn't unmarshal subreddit list.")
		return err
	}

	return nil
}

func initScraper() error {
	err := getSubredditList("reddit_list.json")
	if err != nil {
		return err
	}

	bot, err = reddit.NewBotFromAgentFile("mergenbot.agent", 0)
	if err != nil {
		log.Println("[Scraper:reddit] Failed to create bot handle: ", err)
		return err
	}

	return nil
}

func scrapeReddit() ([]post.Post, error) {
	if bot == nil {
		err := initScraper()
		if err != nil {
			return nil, err
		}
	}
	var posts []post.Post


	for _, subreddit := range subList.Subreddits {
		harvest, err := bot.Listing(subreddit, "")
		if err != nil {
			log.Printf("[Scraper:reddit] Failed to fetch %s : %s", subreddit, err)
			return nil, err
		}

		for _, submission := range harvest.Posts[:20] {
			p := post.Post{
				Title:     submission.Title,
				Source:    fmt.Sprintf("Reddit %s", subreddit),
				Author:    submission.Author,
				Text:      submission.SelfText,
				Url:       submission.Permalink,
				Timestamp: int64(submission.CreatedUTC),
				Score:     int64(submission.Score),
			}
			posts = append(posts, p)
		}
	}

	return posts, nil
}
