package token

import (
	"context"
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/models"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/common"
	"github.com/friendsofgo/errors"
	"github.com/sirupsen/logrus"
	"time"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . dataPersistence
type dataPersistence interface {
	GetAll(ctx context.Context) ([]models.Token, error)
	Generate(ctx context.Context, createdBy int, randomStringsOverride string, createdAtOverride *time.Time) (string, error)
	GetToken(ctx context.Context, key string) (*models.Token, error)
	UpdateTokenToExpired(ctx context.Context, token *models.Token) error
	RevokeToken(ctx context.Context, key string) error
}

type BusinessToken struct {
	dataLayer      dataPersistence
	tokenDaysValid int
}

var (
	errGetToken               = errors.New("error, get fails")
	errRevokeToken            = errors.New("error revoking token")
	errTokenRevoked           = errors.New("error, token is revoked")
	errTokenExpired           = errors.New("error, token has already expired")
	errTokenDeterminedExpired = errors.New("error, token has already expired")
	errUpdateTokenToExpired   = errors.New("error, updating token to expired failed")
)

func (b *BusinessToken) GetAll(ctx context.Context) ([]models.Token, error) {
	return b.dataLayer.GetAll(ctx)
}

func (b *BusinessToken) Generate(ctx context.Context, user *models_schema.User) (string, error) {
	return b.dataLayer.Generate(ctx, user.ID, "", nil)
}

func (b *BusinessToken) Revoke(ctx context.Context, key string) error {
	err := b.dataLayer.RevokeToken(ctx, key)
	if err != nil {
		return errors.Wrap(err, errRevokeToken.Error())
	}
	return nil
}

func (b *BusinessToken) Validate(ctx context.Context, key string) error {
	logger := common.GetLogger(ctx)
	token, err := b.dataLayer.GetToken(ctx, key)
	if err != nil {
		return errors.Wrap(err, errGetToken.Error())
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
	daysElapsed := int(token.ExpiresAt.Sub(token.CreatedAt).Hours() / 7)
	if daysElapsed > b.tokenDaysValid {
		defer func() {
			err = b.dataLayer.UpdateTokenToExpired(ctx, token)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Error("error_validate")
			}
		}()
		logger.WithFields(logrus.Fields{
			"err": errors.Wrap(err, errUpdateTokenToExpired.Error()),
		}).Error("error_validate")
		return errTokenDeterminedExpired
	}

	return nil
}

func NewBusinessToken(mysqlDataPersistence dataPersistence, tokenDaysValid int) *BusinessToken {
	return &BusinessToken{
		dataLayer:      mysqlDataPersistence,
		tokenDaysValid: tokenDaysValid,
	}
}
