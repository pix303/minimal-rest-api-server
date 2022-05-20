package persistence

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresInitDB(t *testing.T) {
	_, err := NewPostgresqlPersistenceService("")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New(DB_CONNECTION_VALID_ERROR), err)
	}

	ps, err := NewPostgresqlPersistenceService("postgres://admin:admin123@localhost:5432/example-test?sslmode=disable")
	assert.Nil(t, err)
	assert.NotNil(t, ps)
	assert.NotNil(t, ps.db)

	ps.ResetDBInstance()
	assert.Nil(t, ps.db)
}
