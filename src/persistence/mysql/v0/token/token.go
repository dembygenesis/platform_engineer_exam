package token

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
)

type PersistenceToken struct {
	db *sql.DB
}

const (
	min = 6
	max = 12
)

// Generate returns a return string in the length range of 6-12
func (p *PersistenceToken) Generate() (string, error) {
	charLength := rand.Intn(max-min) + min
	return generateRandomCharacters(charLength), nil
}

func (p *PersistenceToken) Validate(s string) error {
	return nil
}

func NewPersistenceToken(db *sql.DB) (*PersistenceToken, error) {
	return &PersistenceToken{db}, nil
}
