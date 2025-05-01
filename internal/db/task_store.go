package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type TaskStore struct {
	DB *sql.DB
}

func NewPostgresStore(connStr string) (*TaskStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &TaskStore{DB: db}, nil
}

func (s *TaskStore) InitSchema() error {
	_, err := s.DB.Exec(`CREATE TABLE IF NOT EXISTS task_status (
        task_id TEXT PRIMARY KEY,
        status TEXT,
        selected_node TEXT,
        updated_at TIMESTAMP DEFAULT now()
    )`)
	return err
}

func (s *TaskStore) UpdateStatus(taskID, status, node string) error {
	_, err := s.DB.Exec(`
        INSERT INTO task_status (task_id, status, selected_node)
        VALUES ($1, $2, $3)
        ON CONFLICT (task_id) DO UPDATE
        SET status = $2, selected_node = $3, updated_at = now()`,
		taskID, status, node)
	return err
}
