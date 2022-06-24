package token

import (
	"context"
	"database/sql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"math/rand"
	"time"
)

type PersistenceToken struct {
	db *sql.DB
}

const (
	min = 6
	max = 12
)

var (
	errCheckUniqueToken = errors.New("error checking for unique tokens")
	errInsertNewToken   = errors.New("error inserting new token")
)

// Generate returns a unique string in the length range of 6-12 characters
func (p *PersistenceToken) Generate(ctx context.Context, createdBy int) (string, error) {
	var randomString string
	tokenVerifiedUnique := false
	for !tokenVerifiedUnique {
		randomizedCharLength := rand.Intn(max-min) + min
		randomString = generateRandomCharacters(randomizedCharLength)

		token, err := models_schema.Tokens(
			models_schema.TokenWhere.Key.EQ(randomString),
		).All(mysql.BoilCtx, p.db)
		if err != nil {
			return "", errors.Wrap(errCheckUniqueToken, err.Error())
		}
		if len(token) == 0 {
			tokenVerifiedUnique = true
		}
	}
	newToken := models_schema.Token{
		Key:       randomString,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * time.Hour * 24),
	}
	err := newToken.Insert(mysql.BoilCtx, p.db, boil.Infer())
	if err != nil {
		return "", errors.Wrap(errInsertNewToken, err.Error())
	}
	return randomString, nil
}

func (p *PersistenceToken) Validate(ctx context.Context, s string) error {
	/*token := models_schema.Tokens(
		models_schema.TokenWhere.Key.EQ(s),
	).Count(mysql.BoilCtx, )*/
	return nil
}

func NewPersistenceToken(db *sql.DB) *PersistenceToken {
	return &PersistenceToken{db}
}
