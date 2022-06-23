package v0

import (
	"database/sql"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/friendsofgo/errors"
	"time"
)

type MYSQLConnection struct {
	DB *sql.DB
}

var (
	errDatabasePropertyNil = errors.New("db property is nil")
	errDatabasePingError   = errors.New("db ping failure")
)

// Ping checks if the sql instance can be reached
func (c *MYSQLConnection) Ping() error {
	if c.DB == nil {
		return errDatabasePropertyNil
	}
	err := c.DB.Ping()
	if err != nil {
		return errors.Wrap(errDatabasePingError, err.Error())
	}
	return nil
}

// NewMYSQLConnection returns a struct with a mysql instance
func NewMYSQLConnection(c config.DatabaseCredentials) (*MYSQLConnection, error) {
	conn := MYSQLConnection{}
	connString := c.User + ":" +
		c.Pass + "@tcp(" +
		c.Host + ":" +
		c.Port + ")/" +
		c.Database + "?charset=utf8&parseTime=true"
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, errors.Wrap(err, "error establishing mysql connection")
	}
	db.SetConnMaxLifetime(time.Minute * 60)
	conn.DB = db
	return &conn, nil
}
