package db

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/peter-mount/go-kernel/v2"
	"os"
	"time"
)

// DBService database/sql bound with github.com/lib/pq as a Kernel Service
type DBService struct {
	postgresURI *string
	db          *sql.DB
	maxOpen     int
	maxIdle     int
	maxLifetime time.Duration
	// Set to true to enable additional debugging
	Debug bool
}

func (s *DBService) Init(_ *kernel.Kernel) error {
	s.postgresURI = flag.String("db", "", "The database to connect to")
	return nil
}

func (s *DBService) Start() error {
	if *s.postgresURI == "" {
		*s.postgresURI = os.Getenv("POSTGRESDB")
	}
	if *s.postgresURI == "" {
		return fmt.Errorf("No database uri provided")
	}

	if s.maxOpen < 0 {
		s.maxOpen = 1
	}

	if s.maxIdle < 0 {
		s.maxIdle = 1
	} else if s.maxIdle > s.maxOpen {
		s.maxIdle = s.maxOpen
	}

	db, err := sql.Open("postgres", *s.postgresURI)
	if err != nil {
		return err
	}
	s.db = db

	db.SetMaxOpenConns(s.maxOpen)
	db.SetMaxIdleConns(s.maxIdle)
	if s.maxLifetime > 0 {
		db.SetConnMaxLifetime(s.maxLifetime)
	}

	return nil
}

func (s *DBService) Stop() {
	if s.db != nil {

		_ = s.db.Close()
		s.db = nil
	}
}

// GetDB returns the underlying sql.DB
func (s *DBService) GetDB() *sql.DB {
	return s.db
}

func (s *DBService) SetDB(postgresURI string) *DBService {
	s.postgresURI = &postgresURI
	return s
}

func (s *DBService) MaxOpen(maxOpen int) *DBService {
	s.maxOpen = maxOpen
	return s
}

func (s *DBService) MaxIdle(maxIdle int) *DBService {
	s.maxIdle = maxIdle
	return s
}

func (s *DBService) MaxLifetime(maxLifetime time.Duration) *DBService {
	s.maxLifetime = maxLifetime
	return s
}

func (s *DBService) Exec(query string, args ...interface{}) (sql.Result, error) {
	r, e := s.db.Exec(query, args...)
	return r, e
}

func (s *DBService) Query(query string, args ...interface{}) (*sql.Rows, error) {
	r, e := s.db.Query(query, args...)
	return r, e
}

func (s *DBService) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *DBService) Prepare(sql string) (*sql.Stmt, error) {
	return s.db.Prepare(sql)
}

func (s *DBService) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}
