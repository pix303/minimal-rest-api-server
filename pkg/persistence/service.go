package persistence

import (
	"database/sql"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

type UserPersistencer interface {
	GetUsers() ([]domain.User, error)
}

type PersistenceService struct {
	db *sql.DB
}
