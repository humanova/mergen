package post

import "database/sql"

type Post struct {
	Title 	string
	Author 	string
	Text 	string
	Url 	string
	Timestamp 	int64 // unix time UTC
	Score   int64
}

var database *sql.DB

func InitDb() error {
	db, err := PrepareDb()
	if err != nil {
		return err
	}
	database = db
	return nil
}

func AddPost(newPost Post) error {
	err := InsertNewPost(database, newPost)
	if err != nil {
		return err
	}
	return nil
}

func AddPosts(newPosts []Post) error {
	err := InsertNewPosts(database, newPosts)
	if err != nil {
		return err
	}
	return nil
}