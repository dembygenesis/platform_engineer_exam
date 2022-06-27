package mysql

import (
	"database/sql"
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/src/config"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/strings"
	"github.com/friendsofgo/errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
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
		return err
	}
	return nil
}

var (
	connString2 = fmt.Sprintf(`%v:%v@tcp(%v:%v)/%v`,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SCHEMA"),
	)
)

func establishConnection() {
	db, err := sql.Open("mysql", connString2)

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print("============== wtf 1", err, connString2)
		panic("FUCK")

	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println("======== Gg ping", err)
	} else {
		fmt.Println("======== PING PING PING!!!")
	}
}

// NewMYSQLConnection returns a struct with a mysql instance
func NewMYSQLConnection(c config.DatabaseCredentials) (*MYSQLConnection, error) {
	establishConnection()

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
	err = conn.Ping()
	if err != nil {
		fmt.Println("============ Ping err with creds", strings.GetJSON(c))
		return nil, errors.Wrap(err, errDatabasePingError.Error())
	}
	return &conn, nil
}
