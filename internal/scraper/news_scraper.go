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
)

type Websites struct {
	Websites []Website `json:"websites"`
}

type Website struct {
	Name string       `json:"name"`
	Feeds []string    `json:"feeds"`
}

var feedList Websites

func GetFeedList(path string) error {
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

func RemoveHtmlTag(in string) string {
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
//
func ScrapeNews() ([]post.Post, error) {
	err := GetFeedList("rss_list.json")
	if err != nil {
		return nil, err
	}

	var posts []post.Post

	fp := gofeed.NewParser()

	for _, website := range feedList.Websites {
		log.Printf("[Scraper] Getting feeds : %s\n", website.Name)

		for _, url := range website.Feeds {
			feed, err := fp.ParseURL(url)
			if err != nil {
				log.Printf("[Scraper] Error while parsing %s : %s\n", url, err)
				log.Printf("[Scraper] Skipping current '%s' feed...", website.Name)
				continue
				//return nil, err
			}

			for _, item := range feed.Items {
				p := post.Post{Title: item.Title,
					Author:    website.Name,
					Text:      RemoveHtmlTag(item.Description),
					Url:       item.Link,
					Timestamp: item.PublishedParsed.Unix(),
					Score:     0}
				posts = append(posts, p)
			}
		}
	}

	return posts, nil
}
