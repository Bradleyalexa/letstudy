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
	createSqlTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		taskID INTEGER PRIMARY KEY AUTOINCREMENT, 
		taskContent TEXT NOT NULL,
		taskEndDate DATE,
		taskStatus TEXT DEFAULT 'not done'
	);`

	statement, err := db.Prepare(createSqlTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	fmt.Println("Table Created")
}

func CreatePomodoroTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS pomodoro_sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sessionType TEXT,
		durationMinutes INTEGER,
		startTime DATETIME,
		endTime DATETIME,
		status TEXT
	);`
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatalf("❌ Failed to create pomodoro table: %v", err)
	}
	statement.Exec()
}

func InsertPomodoroSession(sessionType string, duration int, status string, start, end time.Time) error {
	insertSQL := `INSERT INTO pomodoro_sessions (sessionType, durationMinutes, startTime, endTime, status) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(sessionType, duration, start, end, status)
	if err != nil {
		return fmt.Errorf("failed to insert pomodoro session: %v", err)
	}

	return nil
}

type PomodoroSession struct {
	ID          int
	SessionType string
	Duration    int
	StartTime   time.Time
	EndTime     time.Time
	Status      string
}

func GetAllPomodoroSessions() ([]PomodoroSession, error) {
	rows, err := db.Query(`SELECT id, sessionType, durationMinutes, startTime, endTime, status FROM pomodoro_sessions ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []PomodoroSession
	for rows.Next() {
		var s PomodoroSession
		err := rows.Scan(&s.ID, &s.SessionType, &s.Duration, &s.StartTime, &s.EndTime, &s.Status)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
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


func DisplayTasks(status string) {
	if err := OpenDB(); err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer CloseDB()

	CreateTable()

	query := `SELECT taskID, taskContent, taskEndDate, taskStatus FROM tasks WHERE taskStatus = ?`
	rows, err := db.Query(query, status)
	if err != nil {
		log.Fatalf("❌ Failed to query tasks: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var id int
		var content string
		var endDate sql.NullString
		var taskStatus string

		err := rows.Scan(&id, &content, &endDate, &taskStatus)
		if err != nil {
			log.Fatal(err)
		}

		dateDisplay := "—"
		if endDate.Valid {
			dateDisplay = endDate.String
		}

		fmt.Printf("[%d] %s (Due: %s) — %s\n", id, content, dateDisplay, taskStatus)
		found = true
	}

	if !found {
		fmt.Println("✅ No tasks found for this category.")
	}
}

func MarkTaskDone(taskID int) error {
	if err := OpenDB(); err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer CloseDB()

	CreateTable()

	query := `UPDATE tasks SET taskStatus = 'done' WHERE taskID = ?`

	result, err := db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	//error checking if the rows affected, if there is no rows affected == no change == no task found
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %d", taskID)
	}

	fmt.Printf("✅ Task #%d marked as done.\n", taskID)
	return nil
}


func CloseDB(){
	db.Close()
}