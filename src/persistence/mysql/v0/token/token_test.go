package token

import (
	"context"
	"fmt"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

var (
// I'll have a bunch of mocks here
)

// I'll have an assembler function here and combines these mocks,
// and pushes them to my mock function
func validatePassFetchToken(mock sqlmock.Sqlmock) {
	/*sqlInsertToken := `

	`*/
}

/*func TestPersistenceToken_GenerateFailCheckUniqueToken_ForTestsIgnore(t *testing.T) {
	db, mock, err := sqlmock.New()

	// Play with mock
	row := sqlmock.NewRows([]string{"id", "name"})
	row.AddRow(999, "abc")

	sqlToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?);"
	mock.ExpectQuery(regexp.QuoteMeta(sqlToken)).WithArgs("123").WillReturnRows(row)
	// Play with mock

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), 3)
	if err != nil {
		fmt.Println("=========== err ===========", err)
		require.Error(t, err, errCheckUniqueToken)
		errMsg := err.Error()
		wantErrMsg := errCheckUniqueToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	}
}*/

func configureMockGenerateFailFetchToken(mock sqlmock.Sqlmock) {
	sqlToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?);"
	mock.ExpectQuery(regexp.QuoteMeta(sqlToken)).WithArgs("123").WillReturnError(errFetchToken)
}

func TestPersistenceToken_GenerateFailCheckUniqueToken_ForTestsIgnore(t *testing.T) {
	db, mock, err := sqlmock.New()
	configureMockGenerateFailFetchToken(mock)

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), 3)
	if err != nil {
		fmt.Println("=======", err)
		require.Error(t, err, errCheckUniqueToken)

		errMsg := err.Error()
		wantErrMsg := errCheckUniqueToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	}
}

func TestPersistenceToken_GenerateErrCheckUniqueToken(t *testing.T) {

}

func TestPersistenceToken_GenerateErrInsertNewToken(t *testing.T) {

}
