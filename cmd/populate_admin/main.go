package main

import (
	"context"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/common"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/strings"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
	logger := common.GetLogger(context.Background())
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the builder: %v", err.Error())
	}
	ctn := builder.Build()

	mysqlConnection, err := ctn.SafeGetMysqlConnection()
	if err != nil {
		log.Fatalf("error getting the mysql_connection from the container: %v", err.Error())
	}

	name := "Admin User"
	email := "admin@gmail.com"
	unhashedPassword := "123456"
	password, _ := strings.Encrypt(unhashedPassword)
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(unhashedPassword)); err != nil {
		log.Fatalf("Failed to synchronize passwords: %v", err)
	}

	newUser := models_schema.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = newUser.Insert(mysql.BoilCtx, mysqlConnection.DB, boil.Infer())
	if err != nil {
		log.Fatalf("error inserting a dummy admin user: %v", err.Error())
	}

	logger.WithFields(logrus.Fields{
		email:    "admin@gmail.com",
		password: "123456",
	}).Info("add_admin_user")


}
