package scraper

import (
	"encoding/json"
	"fmt"
	"github.com/turnage/graw/reddit"
	"io/ioutil"
	"log"
	"mergen/internal/post"
	"os"
	"sync"
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

func initRedditScraper() error {
	err := getSubredditList("reddit_list.json")
	if err != nil {
		log.Println("[Scraper:reddit] Couldn't load subreddit list: ", err)
		return err
	}

	bot, err = reddit.NewBotFromAgentFile("mergenbot.agent", 0)
	if err != nil {
		log.Println("[Scraper:reddit] Failed to create bot handle: ", err)
		return err
	}

	return nil
}

func scrapeReddit(result chan []post.Post) {
	if bot == nil {
		err := initRedditScraper()
		if err != nil {
			result <- nil
		}
	}
	var posts []post.Post

	log.Printf("[Scraper:reddit] Getting new posts from %d subreddits\n", len(subList.Subreddits))

	wg := sync.WaitGroup{}
	for _, subreddit := range subList.Subreddits {
		harvest, err := bot.Listing(subreddit, "")
		if err != nil {
			log.Printf("[Scraper:reddit] Failed to fetch %s : %s", subreddit, err)
			return
		}

		wg.Add(1)
		go func(harvest reddit.Harvest, posts *[]post.Post) {
			for _, submission := range harvest.Posts[:20] {
				if !submission.Stickied {
					text := submission.SelfText
					if !submission.IsSelf {
						text = submission.URL
					}

					p := post.Post{
						Title:     submission.Title,
						Source:    fmt.Sprintf("Reddit %s", subreddit),
						Author:    submission.Author,
						Text:      text,
						Url:       "https://reddit.com" + submission.Permalink,
						Timestamp: int64(submission.CreatedUTC),
						Score:     int64(submission.Score),
					}
					*posts = append(*posts, p)
				}
			}
			wg.Done()
		}(harvest, &posts)
	}
	wg.Wait()

	result <- posts
}
