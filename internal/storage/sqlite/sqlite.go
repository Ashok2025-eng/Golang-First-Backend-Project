package sqlite

import (
	"database/sql" //PROVIDE INFORMATION RELATED TO DATABASE

	"github.com/Ashok2025-eng/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath) //you can pass postgres also if you use
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    age INTEGER
)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age)VALUES (?,?,?)")
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}
	LastId, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}

	return LastId, nil
}
