package post

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
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
	Language string
}

var database *gorm.DB
var redisClient *redis.Client
var ctx context.Context

func InitDb() error {
	db, err := prepareDb()
	if err != nil {
		return err
	}
	database = db
	return nil
}

func InitRedis() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
	})
	// ping the redis server
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		time.Sleep(3*time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			return err
		}
	}
	ctx = context.Background()

	return nil
}

func Add(newPost Post) error {
	// insert to db
	err := createPost(database, newPost)
	if err != nil {
		return err
	}

	// publish to redis pub/sub
	if redisClient != nil {
		newPostJson, err := json.Marshal(newPost)
		if err != nil {
			return err
		}
		err = redisClient.Publish(ctx, "new_posts", newPostJson).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func AddAll(newPosts []Post) error {
	// insert to db
	err := createPosts(database, newPosts)
	if err != nil {
		return err
	}

	// publish to redis pub/sub
	if redisClient != nil {
		newPostsJson, err := json.Marshal(newPosts)
		if err != nil {
			return err
		}
		err = redisClient.Publish(ctx, "new_posts", newPostsJson).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func GetPostsSince(timestamp int64) ([]Post, error) {
	var posts []Post
	posts, err := getPostsSince(database, timestamp)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostsFiltered(filters Filters) ([]Post, error) {
	var posts []Post
	posts, err := getPostsFiltered(database, filters)
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
