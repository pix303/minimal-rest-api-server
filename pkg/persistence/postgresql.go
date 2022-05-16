package persistence

import (
	"fmt"
	"io/ioutil"
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
func (ps *PersistenceService) GetItems(offset, limit int) ([]Item, error) {
	items := make([]Item, 0)

	rows, err := ps.db.Query(fmt.Sprintf("SELECT id, name, description, pieces FROM items LIMIT %d OFFSET %d", limit, offset))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id string
	var name string
	var desc string
	var pcs int
	for rows.Next() {
		err = rows.Scan(&id, &name, &desc, &pcs)
		if err != nil {
			return nil, err
		}
		items = append(items, Item{ID: id, Name: name, Descripion: desc, Pieces: pcs})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}
