package persistence

import (
	"fmt"
	"io/ioutil"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

// NewPersistenceService manage dbrms connecton and requests
func NewPostgresqlPersistenceService(dbdns string) (*PersistenceService, error) {
	db, err := GetDBInstance(dbdns)
	if err != nil {
		return nil, err
	}
	service := PersistenceService{db}
	err = service.initDB()
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (ps *PersistenceService) initDB() error {
	stmt, err := ioutil.ReadFile("./assets/init.sql")
	if err != nil {
		stmt, err = ioutil.ReadFile("../../assets/init.sql")
		if err != nil {
			return err
		}
	}

	_, err = ps.db.Exec(string(stmt))
	if err != nil {
		return err
	}

	return nil
}

// GetUsers retrive all users
func (ps *PersistenceService) GetUsers(offset, limit int) ([]domain.User, error) {
	users := make([]domain.User, 0)

	rows, err := ps.db.Query(fmt.Sprintf("SELECT id, username FROM users LIMIT %d OFFSET %d", limit, offset))
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
