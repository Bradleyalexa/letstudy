package data

import (
	"database/sql"
	"fmt"
	"log"

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
		taskTitle TEXT NOT NULL,
		taskContent TEXT NOT NULL,
		taskDATE DATE)`

	statement, err := db.Prepare(createSqlTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	fmt.Println("Table Created")
}