package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB

func OpenDB() error {
	var err error
	// glebarez registers the driver name "sqlite"
	db, err = sql.Open("sqlite", "./sqlite-database.db")
	if err != nil {
		return err
	}
	return db.Ping()
}

func CreateTable() {
	createSqlTable := `CREATE TABLE IF NOT EXISTS tasks (
		taskID INTEGER PRIMARY KEY AUTOINCREMENT, 
		taskContent TEXT NOT NULL,
		taskEndDate DATE
		taskStatus TEXT DEFAULT 'not done')`

	statement, err := db.Prepare(createSqlTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	fmt.Println("Table Created")
}

func InsertNote(taskContent string, taskEndDate *time.Time) error {
	insertNoteSQL := `INSERT INTO tasks (taskContent, taskEndDate) VALUES (?, ?)`
	statement, err := db.Prepare(insertNoteSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(taskContent, taskEndDate)
	if err != nil {
		return err
	}

	return nil
}
