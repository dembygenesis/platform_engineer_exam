package token

import (
	"context"
	"database/sql"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/dembygenesis/platform_engineer_exam/models"
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

func TestPersistenceToken_Generate_HappyPath(t *testing.T) {
	db, mock, err := sqlmock.New()
	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	configureMockGeneratePassFetchToken(mock, randomString)
	configureMockGeneratePassInsertToken(mock, randomString, createdById, createdAt)

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
	t.Run("Test Generate Happy Path", func(t *testing.T) {
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestPersistenceToken_Generate_FailCheckUniqueToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	configureMockGenerateFailFetchToken(mock)

	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	persistenceToken := PersistenceToken{db: db}
	_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
	t.Run("Test Generate Fail Check Unique Token", func(t *testing.T) {
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errFetchToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestPersistenceToken_Generate_FailInsertNewToken(t *testing.T) {
	randomString := generateRandomCharacters(12)
	createdAt := time.Now()
	createdById := 3

	db, mock, err := sqlmock.New()
	configureMockGeneratePassFetchToken(mock, randomString)
	configureMockGenerateFailInsertToken(mock, randomString, createdAt)

	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Generate Fail Insert New Token", func(t *testing.T) {
		_, err = persistenceToken.Generate(context.Background(), createdById, randomString, &createdAt)
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errInsertNewToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func configureMockValidatePassGetToken(mock sqlmock.Sqlmock, key string) {
	sqlGetToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?) LIMIT 1;"
	rows := sqlmock.NewRows([]string{"key", "expired", "revoked"}).AddRow(key, false, false)
	mock.ExpectQuery(regexp.QuoteMeta(sqlGetToken)).WithArgs(key).WillReturnRows(rows)
}

func TestPersistenceToken_Validate_HappyPath(t *testing.T) {
	randomString := generateRandomCharacters(12)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	lapseLimit := models.SevenDaysLapse
	lapseType := "7 Days"

	configureMockValidatePassGetToken(mock, randomString)
	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Validate Happy Path", func(t *testing.T) {
		err = persistenceToken.Validate(context.Background(), randomString, lapseLimit, lapseType)
		require.NoError(t, err)
	})
}

func configureMockValidateFailErrTokenNotFound(mock sqlmock.Sqlmock, key string) {
	sqlGetToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?) LIMIT 1;"
	mock.ExpectQuery(regexp.QuoteMeta(sqlGetToken)).WithArgs(key).WillReturnError(sql.ErrNoRows)
}

func TestPersistenceToken_Validate_FailErrTokenNotFound(t *testing.T) {
	randomString := generateRandomCharacters(12)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	lapseLimit := models.SevenDaysLapse
	lapseType := "7 Days"

	configureMockValidateFailErrTokenNotFound(mock, randomString)
	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Validate Fail Err Token Not Found", func(t *testing.T) {
		err = persistenceToken.Validate(context.Background(), randomString, lapseLimit, lapseType)
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenNotFound.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func configureMockValidateFailErrFetchToken(mock sqlmock.Sqlmock, key string) {
	sqlGetToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?) LIMIT 1;"
	mock.ExpectQuery(regexp.QuoteMeta(sqlGetToken)).WithArgs(key).WillReturnError(errFetchToken)
}

func TestPersistenceToken_Validate_FailErrFetchToken(t *testing.T) {
	randomString := generateRandomCharacters(12)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	lapseLimit := models.SevenDaysLapse
	lapseType := "7 Days"

	configureMockValidateFailErrFetchToken(mock, randomString)
	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Validate Fail Err Fetch Token", func(t *testing.T) {
		err = persistenceToken.Validate(context.Background(), randomString, lapseLimit, lapseType)
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errFetchToken.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func configureMockValidatePassGetTokenFailRevoked(mock sqlmock.Sqlmock, key string) {
	sqlGetToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?) LIMIT 1;"
	rows := sqlmock.NewRows([]string{"key", "expired", "revoked"}).AddRow(key, false, true)
	mock.ExpectQuery(regexp.QuoteMeta(sqlGetToken)).WithArgs(key).WillReturnRows(rows)
}

func configureMockValidatePassGetTokenFailExpired(mock sqlmock.Sqlmock, key string) {
	sqlGetToken := "SELECT `token`.* FROM `token` WHERE (`token`.`key` = ?) LIMIT 1;"
	rows := sqlmock.NewRows([]string{"key", "expired", "revoked"}).AddRow(key, true, false)
	mock.ExpectQuery(regexp.QuoteMeta(sqlGetToken)).WithArgs(key).WillReturnRows(rows)
}

func TestPersistenceToken_Validate_FailErrTokenRevoked(t *testing.T) {
	randomString := generateRandomCharacters(12)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	lapseLimit := models.SevenDaysLapse
	lapseType := "7 Days"

	configureMockValidatePassGetTokenFailRevoked(mock, randomString)
	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Validate Fail Err Token Revoked", func(t *testing.T) {
		err = persistenceToken.Validate(context.Background(), randomString, lapseLimit, lapseType)
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenRevoked.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestPersistenceToken_Validate_FailErrTokenExpired(t *testing.T) {
	randomString := generateRandomCharacters(12)

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	lapseLimit := models.SevenDaysLapse
	lapseType := "7 Days"

	configureMockValidatePassGetTokenFailExpired(mock, randomString)
	persistenceToken := PersistenceToken{db: db}
	t.Run("Test Validate Fail Err Token Expired", func(t *testing.T) {
		err = persistenceToken.Validate(context.Background(), randomString, lapseLimit, lapseType)
		require.Error(t, err)

		errMsg := err.Error()
		wantErrMsg := errTokenExpired.Error()
		assert.Containsf(t, errMsg, wantErrMsg, "expected error containing %q, got %s", wantErrMsg, err)
	})
}

func TestPersistenceToken_Validate_FailErrDeterminedExpired(t *testing.T) {

}
