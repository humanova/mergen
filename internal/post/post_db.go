package post

import (
	"database/sql"
	"fmt"
	_	"github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

func PrepareDb() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "./posts.db")
	if err != nil {
		log.Fatal(fmt.Sprintf("[DB] couldn't connect to db : %s", err))
		return nil, err
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY, title TEXT, author TEXT, text TEXT, url TEXT, timestamp DATETIME, score INTEGER)")
	if err != nil {
		log.Fatal(fmt.Sprintf("[DB] couldn't create new table : %s", err))
		return nil, err;
	}
	statement.Exec()

	return database, nil
}

func InsertNewPost(database *sql.DB, post Post) error {
	statement, err := database.Prepare("INSERT INTO posts (title, author, text, url, timestamp, score) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(fmt.Sprintf("[DB] couldn't insert new post : %s", err))
		return err
	}
	statement.Exec(post.Title, post.Author, post.Text, post.Url, post.Timestamp, post.Score)

	return nil
}

func InsertNewPosts(database *sql.DB, posts []Post) error {
	sqlStr := "INSERT INTO posts (title, author, text, url, timestamp, score) VALUES"
	insertArgs := []interface{}{}

	for _, post := range posts {
		sqlStr += "(?, ?, ?, ?, ?, ?),"
		insertArgs = append(insertArgs, post.Title, post.Author, post.Text, post.Url, post.Timestamp, post.Score)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	statement, err := database.Prepare(sqlStr)
	if err != nil {
		log.Fatal(fmt.Sprintf("[DB] couldn't insert new post : %s", err))
		return err
	}
	statement.Exec(insertArgs...)
	return nil
}