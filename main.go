package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func main() {
	// 1. Ensure all three arguments are provided
	if len(os.Args) < 4 {
		fmt.Println("Error: Missing arguments.")
		fmt.Println("Usage: mood_track <mood_string> <trigger_type: active|scheduled|random> <device_name>")
		os.Exit(1)
	}

	// 2. Extract and clean inputs
	mood := strings.TrimSpace(os.Args[1])
	triggerType := strings.ToLower(os.Args[2])
	device := strings.TrimSpace(os.Args[3])

	// 3. Validate that mood string is not empty
	if mood == "" {
		fmt.Println("Error: Mood string cannot be empty.")
		os.Exit(1)
	}

	// 4. Validate trigger type based on ESM methodology
	if triggerType != "active" && triggerType != "scheduled" && triggerType != "random" {
		fmt.Println("Error: Invalid trigger type. Allowed values: active, scheduled, random.")
		os.Exit(1)
	}

	// 5. Validate device name is not empty
	if device == "" {
		fmt.Println("Error: Device name cannot be empty.")
		os.Exit(1)
	}

	// 6. Capture current local time with correct timezone offset (ISO 8601)
	currentTime := time.Now().Format(time.RFC3339)

	// 7. Set up database file path in user's home directory (~/mood_tracker.db)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error: Failed to resolve home directory: %v\n", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(homeDir, "mood_tracker.db")

	// 8. Connect to SQLite database (creates file if missing)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Printf("Error: Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// 9. Create table with 'mood' as TEXT to support flexible string entries
	createTableSQL := `CREATE TABLE IF NOT EXISTS mood_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp TEXT NOT NULL,
		mood TEXT NOT NULL,
		trigger_type TEXT NOT NULL,
		device TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		fmt.Printf("Error: Failed to verify table structure: %v\n", err)
		os.Exit(1)
	}

	// 10. Insert data safely using placeholders
	insertSQL := `INSERT INTO mood_logs (timestamp, mood, trigger_type, device) VALUES (?, ?, ?, ?);`
	_, err = db.Exec(insertSQL, currentTime, mood, triggerType, device)
	if err != nil {
		fmt.Printf("Error: Failed to write log entry: %v\n", err)
		os.Exit(1)
	}

	// Quiet success exit for clean background execution
	os.Exit(0)
}
