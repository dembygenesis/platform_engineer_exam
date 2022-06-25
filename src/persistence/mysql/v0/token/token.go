package token

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/common"
	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	errCheckUniqueToken       = errors.New("error checking for unique tokens")
	errFetchToken             = errors.New("error fetching token")
	errFetchTokens            = errors.New("error fetching tokens")
	errTokenRevoked           = errors.New("error, token is revoked")
	errInsertNewToken         = errors.New("error inserting new token")
	errTokenNotFound          = errors.New("error, token not found")
	errTokenExpired           = errors.New("error, token has already expired")
	errTokenDeterminedExpired = errors.New("error, token has already expired")
)

// GetAll returns all the tokens
func (p *PersistenceToken) GetAll(ctx context.Context) ([]models.Token, error) {
	var container []models.Token
	err := models_schema.Tokens(
		qm.InnerJoin("user u ON u.id = token.created_by"),
		qm.Select([]string{
			"token.id AS id",
			"token.key AS `key`",
			"token.created_at AS created_at",
			"token.revoked AS revoked",
			"token.expired AS expired",
			"token.expires_at AS expires_at",
			"u.name AS created_by",
		}...),
	).Bind(mysql.BoilCtx, p.db, &container)
	if err != nil {
		return nil, errors.Wrap(err, errFetchTokens.Error())
	}

	return container, nil
}

// Generate returns a unique string in the length range of 6-12 characters
func (p *PersistenceToken) Generate(ctx context.Context, createdBy int, randomStringsOverride string, createdAtOverride *time.Time) (string, error) {
	logger := common.GetLogger(ctx)
	var randomString string
	tokenVerifiedUnique := false
	loops := 0
	for !tokenVerifiedUnique {
		if loops > 0 {
			logger.WithFields(logrus.Fields{
				"msg":   "number of loops it took to generate the unique key",
				"loops": loops,
			}).Info("info_generate")
		}
		loops++

		randomizedCharLength := rand.Intn(max-min) + min
		if len(randomStringsOverride) == 0 {
			randomString = generateRandomCharacters(randomizedCharLength)
		} else {
			randomString = randomStringsOverride
		}

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
	var createdAt time.Time
	if createdAtOverride != nil {
		createdAt = *createdAtOverride
	}
	newToken := models_schema.Token{
		Key:       randomString,
		CreatedBy: createdBy,
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(7 * time.Hour * 24),
	}

	err := newToken.Insert(mysql.BoilCtx, p.db, boil.Infer())
	if err != nil {
		return "", errors.Wrap(err, errInsertNewToken.Error())
	}

	return randomString, nil
}

// Validate checks if the token has hit it's 7 day expiry, and if it has - updates the flag "expired" to be true
func (p *PersistenceToken) Validate(ctx context.Context, str string, lapseLimit float64, lapseType string) error {
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
	if token.Expired {
		logger.WithFields(logrus.Fields{
			"msg": fmt.Sprintf("Token: '%v', has already expired.", token.Key),
		}).Error("error_validate")
		return errTokenExpired
	}
	daysElapsed := token.ExpiresAt.Sub(token.CreatedAt).Hours() / 7
	if daysElapsed > lapseLimit {
		defer func() {
			token.Expired = true
			_, err = token.Update(mysql.BoilCtx, p.db, boil.Infer())
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("error_validate")
			}
		}()
		logger.WithFields(logrus.Fields{
			"msg": fmt.Sprintf("Lapsed token detected, with lapse type: '%v'.", lapseType),
		}).Error("error_validate")
		return errTokenDeterminedExpired
	}
	return nil
}

// NewPersistenceToken returns a new *PersistenceToken instance
func NewPersistenceToken(db *sql.DB) *PersistenceToken {
	return &PersistenceToken{db}
}
