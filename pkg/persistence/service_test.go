package persistence

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDBInstance(t *testing.T) {

	// Test error with wrong db connection string
	db, err := GetDBInstance("")
	assert.Nil(t, db)
	if assert.Error(t, err) {
		expectedError := errors.New(DB_CONNECTION_VALID_ERROR)
		assert.Equal(t, expectedError, err)
	}

	// Test with correct db connection string
	db, err = GetDBInstance("postgres://admin:admin123@localhost:5432/example-test?sslmode=disable")
	assert.Nil(t, err)
	assert.NotNil(t, db)

	db = nil
}
