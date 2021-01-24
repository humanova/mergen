package scraper

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	_ "github.com/mmcdole/gofeed/rss"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"mergen/internal/post"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Websites struct {
	Websites []Website `json:"websites"`
}

type Website struct {
	Name  string   `json:"name"`
	Feeds []string `json:"feeds"`
}

var feedList Websites

func removeHtmlTag(in string) string {
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)

	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func getFeedList(path string) error {
	feedFile, err := os.Open(path)
	if err != nil {
		log.Println("[Scraper] Couldn't open feed list file.")
		return err
	}
	defer feedFile.Close()

	byteValue, _ := ioutil.ReadAll(feedFile)
	err = json.Unmarshal(byteValue, &feedList)
	if err != nil {
		log.Println("[Scraper] Couldn't unmarshal feed list.")
		return err
	}

	return nil
}

func ScrapeNews() ([]post.Post, error) {
	err := getFeedList("rss_list.json")
	if err != nil {
		return nil, err
	}

	var posts []post.Post
	keys := make(map[string]bool)

	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 5 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   2 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
		},
	}

	// get posts from rss feeds with concurrency
	wg := sync.WaitGroup{}

	for _, website := range feedList.Websites {
		log.Printf("[Scraper] Getting feeds : %s\n", website.Name)

		wg.Add(1)
		go func(website Website, keys *map[string]bool, posts *[]post.Post) {
			for _, feedUrl := range website.Feeds {
				feed, err := fp.ParseURL(feedUrl)
				if err != nil {
					log.Printf("[Scraper] Skipping current '%s' feed : %s", website.Name, err)
					continue
				}

				for _, item := range feed.Items {
					// check if it's a duplicate
					if _, value := (*keys)[item.Link]; !value {
						(*keys)[item.Link] = true

						itemAuthor := "None"
						if item.Author != nil {
							itemAuthor = item.Author.Name
						}
						itemTitle := html.UnescapeString(item.Title)
						itemText := removeHtmlTag(html.UnescapeString(item.Description))
						urlParsed, err := url.Parse(item.Link)
						if err != nil {
							log.Printf("[Scraper] Skipping current '%s' feed : %s", website.Name, err)
							continue
						}
						itemUrl := urlParsed.Host + urlParsed.Path

						p := post.Post{
							Title:     itemTitle,
							Source:    website.Name,
							Author:    itemAuthor,
							Text:      itemText,
							Url:       itemUrl,
							Timestamp: item.PublishedParsed.Unix(),
							Score:     0}
						*posts = append(*posts, p)
					}
				}
			}
			wg.Done()
		}(website, &keys, &posts)
	}
	wg.Wait()

	return posts, nil
}
