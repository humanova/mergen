package scraper

import (
	"context"
	"encoding/json"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"io/ioutil"
	"log"
	"mergen/internal/post"
	"os"
	"sync"
)

type Accounts struct {
	accounts []string `json:"accounts"`
}

var accountList Accounts

func getAccountList(path string) error {
	accFile, err := os.Open(path)
	if err != nil {
		log.Println("[Scraper:twitter] Couldn't open twitter list file.")
		return err
	}
	defer accFile.Close()

	byteValue, _ := ioutil.ReadAll(accFile)
	err = json.Unmarshal(byteValue, &accountList)
	if err != nil {
		log.Println("[Scraper:twitter] Couldn't unmarshal twitter list.")
		return err
	}

	return nil
}

func scrapeTwitter() ([]post.Post, error)  {
	err := getAccountList("twitter_list.json")
	if err != nil {
		return nil, err
	}
	var posts []post.Post

	scraper := twitterscraper.New()
	log.Printf("[Scraper:twitter] Getting new tweets from %d twitter accounts\n", len(accountList.accounts))

	wg := sync.WaitGroup{}
	for _, username := range accountList.accounts {
		for tweet := range scraper.GetTweets(context.Background(), username, 50) {
			wg.Add(1)
			if tweet.Error != nil {
				log.Printf("[Scraper:twitter] Couldn't scrape a tweet, skipping : %s", tweet.Error)
			}

			if !tweet.IsRetweet || tweet.IsQuoted {
				p := post.Post{
					Title: tweet.Username + "'s tweet",
					Source: "Twitter",
					Author: tweet.Username,
					Text: tweet.Text,
					Url: tweet.PermanentURL,
					Timestamp: tweet.TimeParsed.Unix(),
					Score: int64(tweet.Retweets*10 + tweet.Likes),
				}
				posts = append(posts, p)
			}
		wg.Done()
		}
	}
	wg.Wait()

	return posts, nil
}