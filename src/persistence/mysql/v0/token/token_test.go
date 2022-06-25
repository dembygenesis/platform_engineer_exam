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

func configureMockGenerateFailInsertToken(mock sqlmock.Sqlmock, randomString string, createdAt time.Time) {
	var mockIdReturned int64 = 1
	sqlInsert := "INSERT INTO `token` (`key`,`created_at`,`created_by`,`expires_at`) VALUES (?,?,?,?)"
	mock.ExpectExec(regexp.QuoteMeta(sqlInsert)).WithArgs(
		randomString,
		createdAt,
		3,
		createdAt.Add(7*time.Hour*24),
	).WillReturnResult(sqlmock.NewResult(mockIdReturned, 1)).WillReturnError(errInsertNewToken)

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

func TestPersistenceToken_Generate_FailCheckUniqueToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	configureMockGenerateFailFetchToken(mock)

	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
	require.Error(t, err)

	errMsg := err.Error()
	wantErrMsg := errFetchToken.Error()
	assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
}

func TestPersistenceToken_Generate_FailInsertNewToken(t *testing.T) {
	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	db, mock, err := sqlmock.New()
	configureMockGeneratePassFetchToken(mock, randomString)
	configureMockGenerateFailInsertToken(mock, randomString, createdAt)

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
	require.Error(t, err)

	errMsg := err.Error()
	wantErrMsg := errInsertNewToken.Error()
	assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
}

func TestPersistenceToken_Validate_HappyPath(t *testing.T) {

}

func TestPersistenceToken_Validate_FailErrFetchToken(t *testing.T) {

}

func TestPersistenceToken_Validate_FailErrRevoked(t *testing.T) {

}

func TestPersistenceToken_Validate_FailErrExpired(t *testing.T) {

}

func TestPersistenceToken_Validate_FailErrDeterminedExpired(t *testing.T) {

}
