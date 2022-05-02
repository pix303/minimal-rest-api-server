package persistence

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

// UserPersistencer rappresents users persist actions
type UserPersistencer interface {
	GetUsers(offset, limit int) ([]domain.User, error)
}

// PersistenceService wrap db connector
type PersistenceService struct {
	db *sql.DB
}

var dbInstance *sql.DB

// GetDBInstance return postgresql db connection or create it if not exists
func GetDBInstance(dbdns string) (*sql.DB, error) {
	var err error
	if dbInstance == nil {
		if strings.HasPrefix(dbdns, "postgres") {
			dbInstance, err = sql.Open("postgres", dbdns)
		} else {
			err = errors.New("no db connection valid")
		}

		if err != nil {
			return nil, err
		}

		err = dbInstance.Ping()
		if err != nil {
			return nil, err
		}
	}
	return dbInstance, nil
}
