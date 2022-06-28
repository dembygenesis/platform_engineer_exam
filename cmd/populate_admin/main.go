package main

import (
	"context"
	"fmt"
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
	fmt.Println("=============== 1")
	logger := common.GetLogger(context.Background())
	builder, err := dic.NewBuilder()
	if err != nil {
		log.Fatalf("error trying to initialize the builder: %v", err.Error())
	}
	ctn := builder.Build()

	fmt.Println("=============== 2")
	mysqlConnection, err := ctn.SafeGetMysqlConnection()
	if err != nil {
		log.Fatalf("error getting the mysql_connection from the container: %v", err.Error())
	}

	fmt.Println("=============== 3")
	name := "Admin User"
	email := "admin@gmail.com"
	unhashedPassword := "123456"
	password, _ := strings.Encrypt(unhashedPassword)
	fmt.Println("=============== 4")
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(unhashedPassword)); err != nil {
		log.Fatalf("Failed to synchronize passwords: %v", err)
	}

	fmt.Println("=============== 5")
	newUser := models_schema.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = newUser.Insert(mysql.BoilCtx, mysqlConnection.DB, boil.Infer())
	if err != nil {
		log.Fatalf("error inserting a dummy admin user: %v", err.Error())
	}
	fmt.Println("=============== 6")
	logger.WithFields(logrus.Fields{
		email:    "admin@gmamil.com",
		password: "123456",
	})

	/*fmt.Println("Successfully added a new admin with credentials of", map[string]string{
		email:    "admin@gmamil.com",
		password: "123456",
	})*/
}
