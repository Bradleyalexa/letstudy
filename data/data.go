package data

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	_ "github.com/glebarez/go-sqlite"
)

var db *sql.DB


// Database Initialization
func OpenDB() error {
	var err error
	db, err = sql.Open("sqlite", "./sqlite-database.db")
	if err != nil {
		return err
	}
	return db.Ping()
}

func CloseDB() {
	db.Close()
}

// Create Tables
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

func CreateNotesTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS notes (
    	noteID INTEGER PRIMARY KEY AUTOINCREMENT,
    	content TEXT NOT NULL,
    	createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
    	updatedAt DATETIME
	);`
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatalf("❌ Failed to create notes table: %v", err)
	}
	statement.Exec()
}

func CreateReflectionTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS reflections (
		reflectionID INTEGER PRIMARY KEY AUTOINCREMENT,
		taskID INTEGER,
		taskContent TEXT,
		insight TEXT,
		improvement TEXT,
		rating INTEGER,
		createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	statement, err := db.Prepare(createTable)
	if err != nil {
		log.Fatalf("❌ Failed to create reflections table: %v", err)
	}
	statement.Exec()
}

// Helper
func formatTime(t string) string {
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return t
	}
	return parsed.Format("2 Jan 2006, 15:04")
}


// TASK FUNCTIONS
func InsertNote(taskContent string, taskEndDate *time.Time) error {
	insertNoteSQL := `INSERT INTO tasks (taskContent, taskEndDate) VALUES (?, ?)`
	statement, err := db.Prepare(insertNoteSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	var dateValue interface{}
	if taskEndDate != nil {
		dateValue = taskEndDate.Format(time.RFC3339)
	} else {
		dateValue = nil
	}

	_, err = statement.Exec(taskContent, dateValue)
	return err
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
			dateDisplay = formatTime(endDate.String)
		}

		fmt.Printf("[%d] %s (Due: %s) — %s\n", id, content, dateDisplay, taskStatus)
		found = true
	}

	if !found {
		if status == "done" {
			fmt.Println("✅ No completed tasks yet.")
		} else {
			fmt.Println("✅ No active tasks found.")
		}
	}
}

// func MarkTaskDone(taskID int) error {
// 	if err := OpenDB(); err != nil {
// 		return fmt.Errorf("failed to open database: %v", err)
// 	}
// 	defer CloseDB()

// 	CreateTable()

// 	query := `UPDATE tasks SET taskStatus = 'done' WHERE taskID = ?`
// 	result, err := db.Exec(query, taskID)
// 	if err != nil {
// 		return fmt.Errorf("failed to update task: %v", err)
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no task found with ID %d", taskID)
// 	}

// 	fmt.Printf("✅ Task #%d marked as done.\n", taskID)
// 	return nil
// }

//REMINDER FUNCTION
type TaskReminder struct {
	ID      int
	Content string
	Due     time.Time
}

func GetUpcomingTasks(within time.Duration) ([]TaskReminder, error) {
	CreateTable()

	rows, err := db.Query(`
		SELECT taskID, taskContent, taskEndDate 
		FROM tasks 
		WHERE taskEndDate IS NOT NULL 
		  AND taskStatus = 'not done'`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []TaskReminder
	now := time.Now()
	threshold := now.Add(within)

	for rows.Next() {
		var id int
		var content string
		var dueStr sql.NullString

		if err := rows.Scan(&id, &content, &dueStr); err != nil {
			return nil, err
		}

		if dueStr.Valid {
			due, err := time.Parse(time.RFC3339, dueStr.String)
			if err != nil {
				continue
			}
			if due.After(now) && due.Before(threshold) {
				tasks = append(tasks, TaskReminder{
					ID: id, Content: content, Due: due,
				})
			}
		}
	}

	return tasks, nil
}

//POMODORO FUNCTIONS

type PomodoroSession struct {
	ID          int
	SessionType string
	Duration    int
	StartTime   time.Time
	EndTime     time.Time
	Status      string
}

func InsertPomodoroSession(sessionType string, duration int, status string, start, end time.Time) error {
	insertSQL := `INSERT INTO pomodoro_sessions (sessionType, durationMinutes, startTime, endTime, status) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer statement.Close()

	_, err = statement.Exec(sessionType, duration, start, end, status)
	return err
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


// NOTES FUNCTIONS
func InsertQuickNote(content string) error {
	insertSQL := `INSERT INTO notes (content, createdAt) VALUES (?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(content, time.Now().Format(time.RFC3339))
	return err
}

func ListQuickNotes() error {
	rows, err := db.Query(`SELECT noteID, content, createdAt FROM notes ORDER BY createdAt ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	found := false
	fmt.Println("🗒️ Your Notes")
	fmt.Println("------------------------------------")

	for rows.Next() {
		var id int
		var content string
		var createdAt string
		if err := rows.Scan(&id, &content, &createdAt); err != nil {
			return err
		}
		fmt.Printf("[%d] %s (%s)\n", id, content, formatTime(createdAt))
		found = true
	}

	if !found {
		fmt.Println("✅ No notes yet.")
	}

	return nil
}

func ViewQuickNoteByID(id int) error {
	row := db.QueryRow(`SELECT noteID, content, createdAt, updatedAt FROM notes WHERE noteID = ?`, id)

	var noteID int
	var content, createdAt string
	var updatedAt sql.NullString

	err := row.Scan(&noteID, &content, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		fmt.Printf("❌ Note with ID %d not found.\n", id)
		return nil
	} else if err != nil {
		return err
	}

	fmt.Printf("🗒️ Note #%d\n------------------------------------\n%s\n\nCreated: %s\n",
		noteID, content, formatTime(createdAt))
	if updatedAt.Valid {
		fmt.Printf("Updated: %s\n", formatTime(updatedAt.String))
	}

	return nil
}

func DeleteQuickNoteByID(id int) error {
	result, err := db.Exec(`DELETE FROM notes WHERE noteID = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Printf("⚠️ No note found with ID %d.\n", id)
		return nil
	}

	fmt.Printf("🗑️ Note #%d deleted successfully.\n", id)
	return nil
}

func SearchQuickNotes(keyword string) error {
	rows, err := db.Query(`SELECT noteID, content, createdAt FROM notes WHERE content LIKE '%' || ? || '%' ORDER BY createdAt DESC`, keyword)
	if err != nil {
		return err
	}
	defer rows.Close()

	found := false
	fmt.Printf("🔍 Search results for \"%s\":\n", keyword)
	fmt.Println("------------------------------------")

	for rows.Next() {
		var id int
		var content string
		var createdAt string
		if err := rows.Scan(&id, &content, &createdAt); err != nil {
			return err
		}
		fmt.Printf("[%d] %s (%s)\n", id, content, formatTime(createdAt))
		found = true
	}

	if !found {
		fmt.Printf("❌ No notes found matching \"%s\".\n", keyword)
	}

	return nil
}

//reflection
func InsertReflection(taskID int, taskContent, insight, improvement string, rating int) error {
	insertSQL := `
	INSERT INTO reflections (taskID, taskContent, insight, improvement, rating, createdAt)
	VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(taskID, taskContent, insight, improvement, rating, time.Now().Format(time.RFC3339))
	return err
}

func MarkTaskDone(taskID int) error {
	if err := OpenDB(); err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer CloseDB()

	CreateTable()
	CreateReflectionTable()

	
	var taskContent, taskStatus string
	row := db.QueryRow(`SELECT taskContent, taskStatus FROM tasks WHERE taskID = ?`, taskID)
	if err := row.Scan(&taskContent, &taskStatus); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("❌ Task not found with ID %d", taskID)
		}
		return err
	}

	if taskStatus == "done" {
		return fmt.Errorf("⚠️ Task #%d is already marked as done", taskID)
	}

	
	query := `UPDATE tasks SET taskStatus = 'done' WHERE taskID = ?`
	result, err := db.Exec(query, taskID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no task found with ID %d", taskID)
	}

	fmt.Printf("✅ Task #%d marked as done.\n", taskID)
	fmt.Println("🪞 Let's reflect on this task!")

	reader := bufio.NewReader(os.Stdin)

	
	fmt.Print("✨ What went well? → ")
	insight, _ := reader.ReadString('\n')
	insight = strings.TrimSpace(insight)

	fmt.Print("🔧 What can be improved? → ")
	improvement, _ := reader.ReadString('\n')
	improvement = strings.TrimSpace(improvement)

	var rating int
	for {
		fmt.Print("⭐ How satisfied are you (1–5)? → ")
		_, err := fmt.Scan(&rating)
		if err != nil || rating < 1 || rating > 5 {
			fmt.Println("⚠️ Please enter a number between 1 and 5.")
			reader.ReadString('\n') 
			continue
		}
		break
	}

	if err := InsertReflection(taskID, taskContent, insight, improvement, rating); err != nil {
		return fmt.Errorf("failed to save reflection: %v", err)
	}

	fmt.Println("💾 Reflection saved successfully!")
	return nil
}

func ListReflectionSummaries() error {
	rows, err := db.Query(`
		SELECT reflectionID, taskID, taskContent
		FROM reflections
		ORDER BY createdAt ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("💭 Reflections")
	fmt.Println("------------------------------------")

	found := false
	for rows.Next() {
		var reflectionID, taskID int
		var taskContent string
		if err := rows.Scan(&reflectionID, &taskID, &taskContent); err != nil {
			return err
		}

		fmt.Printf("[%d] Task #%d — %s\n", reflectionID, taskID, taskContent)
		found = true
	}

	if !found {
		fmt.Println("✅ No reflections yet.")
	}
	return nil
}

func ViewReflectionByID(reflectionID int) error {
	row := db.QueryRow(`
		SELECT reflectionID, taskID, taskContent, insight, improvement, rating, createdAt
		FROM reflections
		WHERE reflectionID = ?`, reflectionID)

	var id, taskID, rating int
	var taskContent, insight, improvement, createdAt string

	err := row.Scan(&id, &taskID, &taskContent, &insight, &improvement, &rating, &createdAt)
	if err == sql.ErrNoRows {
		fmt.Printf("❌ Reflection with ID %d not found.\n", reflectionID)
		return nil
	} else if err != nil {
		return err
	}

	fmt.Printf("🪞 Reflection #%d — %s\n", id, taskContent)
	fmt.Println("------------------------------------")
	fmt.Printf("✨ What went well: %s\n", insight)
	fmt.Printf("🔧 What can be improved: %s\n", improvement)
	fmt.Printf("⭐ Rating: %d/5\n", rating)
	fmt.Printf("📅 %s\n", formatTime(createdAt))
	return nil
}

