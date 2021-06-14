package post

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"mergen/internal/config"
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
		Addr:     config.Config.RedisHost,
		DB:       config.Config.RedisDB,
	})
	// ping the redis server
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		time.Sleep(3*time.Second)
		err := redisClient.Ping(context.Background()).Err()
		if err != nil {
			redisClient = nil
			return err
		}
	}
	ctx = context.Background()

	return nil
}

func PublishNewPosts() (int, error) {
	if redisClient == nil {
		return 0, errors.New("redis client is not initialized")
	}

	newPosts, err := getPostsUpdatedAfter(database, time.Now().UTC().Add(time.Duration(-config.Config.ScrapeInterval) * time.Minute))
	if err != nil {
		return 0, err
	}
	// publish to redis pub/sub
	if redisClient != nil {
		newPostsJson, err := json.Marshal(newPosts)
		if err != nil {
			return 0, err
		}
		err = redisClient.Publish(ctx, "new_posts", newPostsJson).Err()
		if err != nil {
			return 0, err
		}
	}

	return len(newPosts), nil
}

func Add(newPost Post) error {
	// insert to db
	err := createPost(database, newPost)
	if err != nil {
		return err
	}

	return nil
}

func AddAll(newPosts []Post) error {
	// insert to db
	err := createPosts(database, newPosts)
	if err != nil {
		return err
	}
	return nil
}

func GetPostsSince(timestamp int64) ([]Post, error) {
	var posts []Post
	posts, err := getPostsPublishedAfter(database, timestamp)
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
