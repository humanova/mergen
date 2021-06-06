package post

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"mergen/internal/config"
	"os"
	"time"
)

type Filters struct {
	CreatedAfter  string
	CreatedBefore string
	QueryText     string
	Sources       []string
	Languages     []string
	Author        string
}

func prepareDb() (*gorm.DB, error) {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		   config.Config.DbHost, config.Config.DbUser, config.Config.DbPassword, config.Config.DbName,
		   config.Config.DbPort, config.Config.DbSSLMode)

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

func getPostsPublishedAfter(database *gorm.DB, timestamp int64) ([]Post, error) {
	var posts []Post

	tx := database.Where("timestamp > ?", timestamp).Find(&posts)
	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't query any posts with given timestamp(%d) : %s\n", timestamp, tx.Error))
		return nil, tx.Error
	}

	return posts, nil
}

func getPostsUpdatedAfter(database *gorm.DB, qTime time.Time) ([]Post, error) {
	var posts []Post

	tx := database.Where("updated_at > ?", qTime).Find(&posts)
	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't query any posts with given time(%d) : %s\n", qTime, tx.Error))
		return nil, tx.Error
	}

	return posts, nil
}

func getPostsFiltered(database *gorm.DB, filters Filters) ([]Post, error) {
	var posts []Post

	tx := database.Where("")
	if filters.CreatedAfter != "" {
		tx.Where("timestamp > ?", filters.CreatedAfter)
	}
	if filters.CreatedBefore != "" {
		tx = tx.Where("timestamp < ?", filters.CreatedBefore)
	}
	if filters.Author != "" {
		tx = tx.Where("author ILIKE ?", fmt.Sprintf("%%%s%%", filters.Author))
	}
	if filters.QueryText != "" {
		tx = tx.Where("text ILIKE ? OR title ILIKE ?",
			fmt.Sprintf("%%%s%%", filters.QueryText),
			fmt.Sprintf("%%%s%%", filters.QueryText))
	}
	if filters.Sources != nil {
		tx = tx.Where("source IN ?", filters.Sources)
	}
	if filters.Languages != nil {
		tx = tx.Where("language IN ?", filters.Languages)
	}

	// if nothing is passed as a filter, return posts from last 12 hours
	if filters.CreatedAfter == "" && filters.CreatedBefore == "" &&
		filters.Author == "" && filters.QueryText == "" && filters.Sources == nil {
		tx = tx.Where("timestamp > ?", (time.Now().Add(-12 * time.Hour)).UTC().Unix())
	}

	tx.Find(&posts)

	if tx.Error != nil {
		log.Println(fmt.Sprintf("[DB] couldn't query any posts with given filters(%v) : %s\n", filters, tx.Error))
		return nil, tx.Error
	}

	return posts, nil
}
