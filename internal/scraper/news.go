package scraper

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	_ "github.com/mmcdole/gofeed/rss"
	"html"
	"io/ioutil"
	"log"
	"mergen/internal/post"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
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
		log.Println("[Scraper:news] Couldn't open feed list file.")
		return err
	}
	defer feedFile.Close()

	byteValue, _ := ioutil.ReadAll(feedFile)
	err = json.Unmarshal(byteValue, &feedList)
	if err != nil {
		log.Println("[Scraper:news] Couldn't unmarshal feed list.")
		return err
	}

	return nil
}

func scrapeNews() ([]post.Post, error) {
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
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   4 * time.Second,
			ResponseHeaderTimeout: 2 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true,
		},
	}

	log.Printf("[Scraper:news] Getting news feeds from %d websites\n", len(feedList.Websites))

	for _, website := range feedList.Websites {
		for _, feedUrl := range website.Feeds {
			feed, err := fp.ParseURL(feedUrl)

			defer fp.Client.CloseIdleConnections()

			if err != nil {
				log.Printf("[Scraper:news] Skipping current '%s' feed : %s", website.Name, err)
				continue
			}

			for _, item := range feed.Items {
				// check if it's a duplicate
				urlParsed, err := url.Parse(item.Link)
				if err != nil {
					log.Printf("[Scraper:news] Skipping current item of feed '%s' : %s", website.Name, err)
					continue
				}
				itemUrl := urlParsed.Host + urlParsed.Path

				if _, value := keys[itemUrl]; !value {
					keys[itemUrl] = true

					itemAuthor := "None"
					if item.Author != nil {
						itemAuthor = item.Author.Name
					}
					itemTitle := html.UnescapeString(item.Title)
					itemText := removeHtmlTag(html.UnescapeString(item.Description))

					p := post.Post{
						Title:     itemTitle,
						Source:    website.Name,
						Author:    itemAuthor,
						Text:      itemText,
						Url:       itemUrl,
						Timestamp: item.PublishedParsed.Unix(),
						Score:     0}
					posts = append(posts, p)
				}
			}
		}
	}

	return posts, nil
}
