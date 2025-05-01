
package db

import (
    "database/sql"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./task.db")
    if err != nil {
        log.Fatal(err)
    }
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS task_status (
        task_id TEXT PRIMARY KEY,
        status TEXT,
        updated_at TEXT,
        selected_node TEXT
    )`)
    if err != nil {
        log.Fatal(err)
    }
}

func UpdateTaskStatus(taskID string, status string, selected string) {
    stmt, err := DB.Prepare(`INSERT OR REPLACE INTO task_status(task_id, status, updated_at, selected_node) VALUES (?, ?, ?, ?)`)
    if err != nil {
        log.Println(err)
        return
    }
    _, err = stmt.Exec(taskID, status, time.Now().Format(time.RFC3339), selected)
    if err != nil {
        log.Println(err)
    }
}
