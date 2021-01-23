package scraper

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	_ "github.com/mmcdole/gofeed/rss"
	"io/ioutil"
	"log"
	"mergen/internal/post"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
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

	// get posts from rss feeds with concurrency
	wg := sync.WaitGroup{}

	for _, website := range feedList.Websites {
		log.Printf("[Scraper] Getting feeds : %s\n", website.Name)

		wg.Add(1)
		go func(website Website, keys *map[string]bool, posts *[]post.Post) {
			for _, url := range website.Feeds {
				feed, err := fp.ParseURL(url)
				if err != nil {
					log.Printf("[Scraper] Skipping current '%s' feed : %s", website.Name, err)
					continue
				}

				for _, item := range feed.Items {
					// check if it's a duplicate
					if _, value := (*keys)[item.Link]; !value {
						(*keys)[item.Link] = true

						author := "None"
						if item.Author != nil {
							author = item.Author.Name
						}
						p := post.Post{
							Title:     item.Title,
							Source:    website.Name,
							Author:    author,
							Text:      removeHtmlTag(item.Description),
							Url:       item.Link,
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
