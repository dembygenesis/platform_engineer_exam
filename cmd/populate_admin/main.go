package main

/**
This command inserts a new dummy user to the app with creds:
email: admin@gmail.com
pass:  123456
*/

import (
	"fmt"
	"github.com/dembygenesis/platform_engineer_exam/dependency_injection/dic"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql"
	"github.com/dembygenesis/platform_engineer_exam/src/persistence/mysql/models_schema"
	"github.com/dembygenesis/platform_engineer_exam/src/utils/strings"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func main() {
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

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(unhashedPassword))
	/*err = bcrypt.CompareHashAndPassword([]byte("$2a$10$U7Gu/i.MpomFEuGNPq/.OeyUiIEhNpTTinot/eWFO9UuK58weGp02"), []byte("123456"))
	if err != nil {
		panic("GG  no match")
	} else {
		panic("Match with len: " + strconv.Itoa(len(password)))
	}*/

	newUser := models_schema.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	err = newUser.Insert(mysql.BoilCtx, mysqlConnection.DB, boil.Infer())
	if err != nil {
		log.Fatalf("error inserting a dummy admin user: %v", err.Error())
	}
	fmt.Println("Successfully added a new admin with credentials of", map[string]string{
		email:    "admin@gmamil.com",
		password: "123456",
	})
}
