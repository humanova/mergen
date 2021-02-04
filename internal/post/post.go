package post

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string
	Source    string
	Author    string
	Text      string
	Url       string `gorm:"unique"`
	Timestamp int64 // unix time UTC
	Score     int64
}

var database *gorm.DB

func InitDb() error {
	db, err := prepareDb()
	if err != nil {
		return err
	}
	database = db
	return nil
}

func Add(newPost Post) error {
	err := createPost(database, newPost)
	if err != nil {
		return err
	}
	return nil
}

func AddAll(newPosts []Post) error {
	err := createPosts(database, newPosts)
	if err != nil {
		return err
	}
	return nil
}

func GetPosts(timestamp int64) ([]Post, error) {
	var posts []Post
	posts, err := getPostsSince(database, timestamp)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func UpdateScore(post Post, score int64) error {
	err := updatePostScore(database, post, score)
	if err != nil {
		return err
	}
	return nil
}
