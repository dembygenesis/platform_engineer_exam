package token

import (
	"database/sql"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"math"
	"math/rand"
)

type PersistenceToken struct {
	db *sql.DB
}

const (
	min = 6
	max = 12
)

// generateRandomCharacters generates a string of length "l"
// src: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func (p *PersistenceToken) generateRandomCharacters(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/2)))
	_, _ = rand.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:l]
}

// Generate returns a return string in the length range of 6-12
func (p *PersistenceToken) Generate() (string, error) {
	charLength := rand.Intn(max-min) + min
	return p.generateRandomCharacters(charLength), nil
}

func (p *PersistenceToken) Validate(s string) error {
	return nil
}

func NewPersistenceToken(db *sql.DB) (*PersistenceToken, error) {
	return &PersistenceToken{db}, nil
}
