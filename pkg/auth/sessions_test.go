package auth

import (
	"errors"
	"testing"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
	"github.com/stretchr/testify/assert"
)

func TestNewSessionManager(t *testing.T) {
	sessioner, err := NewSessionManager("")
	assert.Nil(t, sessioner)
	if assert.Error(t, err) {
		expError := errors.New(persistence.DB_CONNECTION_VALID_ERROR)
		assert.Equal(t, expError, err)
	}

	sessioner, err = NewSessionManager("postgres://admin:admin123@localhost:5432/example-test?sslmode=disable")
	assert.Nil(t, err)
	assert.NotNil(t, sessioner)
	assert.NotNil(t, sessioner.StoreManager)
}
