package post

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func prepareDb() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("posts.db"), &gorm.Config{})
	if err != nil {
		log.Print(fmt.Sprintf("[DB] couldn't connect to db : %s", err))
		return nil, err
	}

	err = database.AutoMigrate(&Post{})
	if err != nil {
		log.Print(fmt.Sprintf("[DB] couldn't create new table : %s", err))
		return nil, err;
	}

	return database, nil
}

func createPost(database *gorm.DB, post Post) error {
	tx := database.Create(&post)
	if tx.Error != nil {
		log.Print(fmt.Sprintf("[DB] couldn't insert new post : %s", tx.Error))
		return tx.Error
	}

	return nil
}

func createPosts(database *gorm.DB, posts []Post) error {
	tx := database.Create(&posts)
	if tx.Error != nil {
		log.Print(fmt.Sprintf("[DB] couldn't insert new posts : %s", tx.Error))
		return tx.Error
	}

	return nil
}

func updatePostScore(database *gorm.DB, post Post, score int64) error {
	var p Post
	tx := database.First(&p, "Url = ?", post.Url)
	if tx.Error != nil {
		log.Print(fmt.Sprintf("[DB] couldn't find the post to be updated : %s", tx.Error))
		return tx.Error
	}

	tx = database.Model(&p).Update("Score", score)
	if tx.Error != nil {
		log.Print(fmt.Sprintf("[DB] couldn't update the post : %s", tx.Error))
		return tx.Error
	}

	return nil
}