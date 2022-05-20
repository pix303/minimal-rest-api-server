package persistence

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/lib/pq"
)

const (
	DB_CONNECTION_VALID_ERROR = "no db connection valid"
)

type Item struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name"`
	Descripion string `json:"description"`
	Quantity   int    `json:"quantity,omitempty"`
}

// ItemPersistencer rappresents Item CRUD actions
type ItemPersistencer interface {
	GetItems(offset, limit int) ([]Item, error)
	GetItem(id string) (*Item, error)
	PostItem(item Item) (string, error)
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
			err = errors.New(DB_CONNECTION_VALID_ERROR)
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

func (ps *PersistenceService) ResetDBInstance() {
	ps.db.Close()
	ps.db = nil
	dbInstance = nil
}
