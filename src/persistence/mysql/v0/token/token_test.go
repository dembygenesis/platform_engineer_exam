package token

import (
	"context"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
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
	sqlToken := "SELECT (.+) FROM `token` WHERE .*"
	mock.ExpectQuery(regexp.QuoteMeta(sqlToken)).WithArgs("123").WillReturnError(errFetchToken)
}

func configureMockGeneratePassFetchToken(mock sqlmock.Sqlmock, randomString string) {
	sqlToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?);"
	rows := sqlmock.NewRows([]string{
		"id",
		"key",
		"created_at",
		"revoked",
		"expired",
		"created_by",
		"expires_at",
	})
	mock.ExpectQuery(regexp.QuoteMeta(sqlToken)).WithArgs(randomString).WillReturnRows(rows)
}

func configureMockGeneratePassInsertToken(mock sqlmock.Sqlmock, randomString string, createdBy int, createdAt time.Time) {
	var mockIdReturned int64 = 1
	sqlInsert := "INSERT INTO `token` (`key`,`created_at`,`created_by`,`expires_at`) VALUES (?,?,?,?)"
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).WithArgs(
		randomString,
		createdAt,
		3,
		createdAt.Add(7*time.Hour*24),
	).WillReturnResult(sqlmock.NewResult(mockIdReturned, 1))

	sqlPostSelectAfterSQLBoilerInsert := "SELECT `id`,`revoked`,`expired` FROM `token` WHERE `id`=?"
	rows := sqlmock.NewRows([]string{"id", "revoked", "expired"})
	rows.AddRow(mockIdReturned, false, false)
	mock.ExpectQuery(regexp.QuoteMeta(sqlPostSelectAfterSQLBoilerInsert)).WithArgs(
		mockIdReturned,
	).WillReturnRows(rows)
}

func TestPersistenceToken_GenerateFailCheckUniqueToken_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	configureMockGeneratePassFetchToken(mock, randomString)
	configureMockGeneratePassInsertToken(mock, randomString, createdById, createdAt)

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPersistenceToken_GenerateErrCheckUniqueToken(t *testing.T) {

}

func TestPersistenceToken_GenerateErrInsertNewToken(t *testing.T) {

}
