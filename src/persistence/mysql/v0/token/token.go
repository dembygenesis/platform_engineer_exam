package token

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type dataPersistence interface {
	// Generate creates a new 6-12 digit authentication token
	Generate() (string, error)

	// Validate checks if a string is registered
	Validate(s string) error
}

type PersistenceToken struct {
	db *sql.DB
}

func (p *PersistenceToken) Generate() (string, error) {
	return "Generate", nil
}

func (p *PersistenceToken) Validate(s string) error {
	return nil
}

func NewPersistenceToken(db *sql.DB) (*PersistenceToken, error) {
	return &PersistenceToken{db}, nil
}
