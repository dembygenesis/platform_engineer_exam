package user

import (
	"database/sql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/friendsofgo/errors"
	"golang.org/x/crypto/bcrypt"
)

type PersistenceUser struct {
	db *sql.DB
}

var (
	errMatchingEmail = errors.New("error matching email")
)

// BasicAuth authenticates by first matching the user to "email", and the password to it's "encrypted" version
func (p *PersistenceUser) BasicAuth(user, pass string) (bool, *models_schema.User, error) {
	match, err := models_schema.Users(
		models_schema.UserWhere.Email.EQ(user),
	).One(mysql.BoilCtx, p.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, match, nil
		} else {
			return false, match, errors.Wrap(errMatchingEmail, err.Error())
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(match.Password), []byte(pass))
	if err != nil {
		return false, match, nil
	}
	return true, match, nil
}

func NewPersistenceUser(db *sql.DB) *PersistenceUser {
	return &PersistenceUser{db}
}
