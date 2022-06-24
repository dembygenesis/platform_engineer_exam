package token

import (
	"context"
	"database/sql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/common"
	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
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
	sevenDaysKey = (7 * time.Hour * 24).Hours()

	// For testing
	// thirtySecondsKey = (1 * time.Minute / 2).Hours()
)

var (
	errCheckUniqueToken = errors.New("error checking for unique tokens")
	errFetchToken       = errors.New("error fetching token")
	errFetchTokens      = errors.New("error fetching tokens")
	errTokenRevoked     = errors.New("token is revoked")
	errInsertNewToken   = errors.New("error inserting new token")
	errTokenNotFound    = errors.New("error, token not found")
	errTokenExpired     = errors.New("error, token has already expired")
)

// GetAll returns all the tokens
func (p *PersistenceToken) GetAll(ctx context.Context) ([]models_schema.Token, error) {
	var container []models_schema.Token
	err := models_schema.Tokens().Bind(mysql.BoilCtx, p.db, &container)
	if err != nil {
		return nil, errors.Wrap(err, errFetchTokens.Error())
	}
	return container, nil
}

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
			return "", errors.Wrap(err, errCheckUniqueToken.Error())
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

// Validate checks if the token has hit it's 7 day expiry, and if it has - updates the flag "expired" to be true
func (p *PersistenceToken) Validate(ctx context.Context, str string) error {
	logger := common.GetLogger(ctx)
	token, err := models_schema.Tokens(
		models_schema.TokenWhere.Key.EQ(str),
	).One(mysql.BoilCtx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errTokenNotFound
		} else {
			return errors.Wrap(err, errFetchToken.Error())
		}
	}
	if token.Revoked {
		return errTokenRevoked
	}
	daysElapsed := token.ExpiresAt.Sub(token.CreatedAt).Hours() / 7
	if daysElapsed > sevenDaysKey {
		defer func() {
			token.Expired = true
			_, err = token.Update(mysql.BoilCtx, p.db, boil.Infer())
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("error_validate")
			}
		}()
		return errTokenExpired
	}
	return nil
}

// NewPersistenceToken returns a new *PersistenceToken instance
func NewPersistenceToken(db *sql.DB) *PersistenceToken {
	return &PersistenceToken{db}
}
