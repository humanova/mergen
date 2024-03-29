package scraper

import (
	"context"
	"encoding/json"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"io/ioutil"
	"log"
	"mergen/internal/config"
	"mergen/internal/post"
	"net/http"
	"os"
	"sync"
)

type Accounts struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Name string `json:"name"`
	Lang string `json:"lang"`
}

var accountList Accounts
var cookies []*http.Cookie

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

func initTwitterScraper() error {
	err := getAccountList(config.Config.TwitterListPath)
	if err != nil {
		log.Printf("[Scraper:twitter] Couldn't load account list: %s", err)
		return err
	}
	return nil
}

func scrapeTwitter(result chan []post.Post) {
	if accountList.Accounts == nil {
		err := initTwitterScraper()
		if err != nil {
			result <- nil
		}
	}
	var posts []post.Post

	scraper := twitterscraper.New()

	if cookies != nil {
		scraper.SetCookies(cookies)
	} else {
		err := scraper.Login(config.Config.TwitterUsername, config.Config.TwitterPassword)
		if err != nil {
			log.Printf("[Scraper:twitter] Couldn't login to Twitter account %s\n", config.Config.TwitterUsername)
			result <- nil
		}
		if scraper.IsLoggedIn() {
			cookies = scraper.GetCookies()
		}
	}

	log.Printf("[Scraper:twitter] Getting new tweets from %d twitter accounts\n", len(accountList.Accounts))

	wg := sync.WaitGroup{}
	for _, acc := range accountList.Accounts {
		for tweet := range scraper.GetTweets(context.Background(), acc.Name, 50) {
			wg.Add(1)

			go func(tweet *twitterscraper.TweetResult, posts *[]post.Post) {
				if tweet.Error != nil {
					log.Println("[Scraper:twitter] Couldn't scrape a tweet, skipping : ", tweet.Error)
				}

				if !tweet.IsRetweet || tweet.IsQuoted {
					p := post.Post{
						Title:     tweet.Username + "'s tweet",
						Source:    "Twitter",
						Author:    tweet.Username,
						Text:      tweet.Text,
						Url:       tweet.PermanentURL,
						Timestamp: tweet.TimeParsed.Unix(),
						Score:     int64(tweet.Retweets*10 + tweet.Likes),
						Language:  acc.Lang,
					}
					*posts = append(*posts, p)
				}
				wg.Done()
			}(tweet, &posts)
		}
	}
	wg.Wait()

	result <- posts
}
