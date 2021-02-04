package post

import (
	"fmt"
	"gorm.io/gorm/clause"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func prepareDb() (*gorm.DB, error) {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	dsn := "host=localhost user=mergen password=mergen dbname=posts port=5432 sslmode=disable TimeZone=UTC"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})
	if err != nil {
		log.Println(fmt.Sprintf("[DB] couldn't connect to db : %s\n", err))
		return nil, err
	}

	err = database.AutoMigrate(&Post{})
	if err != nil {
		log.Println(fmt.Sprintf("[DB] couldn't create new table : %s\n", err))
		return nil, err
	}

	return database, nil
}

func createPost(database *gorm.DB, post Post) error {
	tx := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "text", "timestamp"}),
	}).Create(&post)


	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't insert new post : %s\n", tx.Error))
		return tx.Error
	}

	return nil
}

func createPosts(database *gorm.DB, posts []Post) error {
	tx := database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "url"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "text", "timestamp"}),
	}).Create(&posts)

	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't insert new posts : %s\n", tx.Error))
		return tx.Error
	}

	return nil
}

func updatePostScore(database *gorm.DB, post Post, score int64) error {
	tx := database.Update("Score", score).Where("url = ?", post.Url)
	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't update the post : %s\n", tx.Error))
		return tx.Error
	}

	return nil
}

func getPostsSince(database *gorm.DB, timestamp int64) ([]Post, error) {
	var posts []Post

	tx := database.Where("timestamp > ?", timestamp).Find(&posts)
	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't query any posts with given timestamp(%d) : %s\n", timestamp, tx.Error))
		return nil, tx.Error
	}

	return posts, nil
}
