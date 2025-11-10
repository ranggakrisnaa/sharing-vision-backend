package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL membuka koneksi ke MySQL menggunakan DSN yang disediakan.
// Contoh DSN: user:pass@tcp(127.0.0.1:3306)/sharing_vision?parseTime=true&charset=utf8mb4&loc=Local
func NewMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
