package persistence

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

var dbInstance *sql.DB

func GetDBInstance(dbdns string) (*sql.DB, error) {
	var err error
	if dbInstance == nil {
		dbInstance, err = sql.Open("postgres", dbdns)
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

// NewPersistenceService manage dbrms connecton and requests
func NewPostgresqlPersistenceService(dbdns string) (*PersistenceService, error) {
	db, err := GetDBInstance(dbdns)
	if err != nil {
		return nil, err
	}
	service := PersistenceService{db}

	return &service, nil
}

// GetUsers retrive all users
func (ps *PersistenceService) GetUsers() ([]domain.User, error) {
	users := make([]domain.User, 0)

	rows, err := ps.db.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int32
	var username string
	for rows.Next() {
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}
		users = append(users, domain.User{ID: id, Username: username})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}
