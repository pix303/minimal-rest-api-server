package persistence

import (
	"io/ioutil"
	"strconv"
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

// GetItems retrive all items with offset limit
func (ps *PersistenceService) GetItems(offset, limit int) ([]Item, error) {
	items := make([]Item, 0)

	rows, err := ps.db.Query("SELECT id, name, description, quantity FROM items LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Descripion,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetItem retrive a single item by id
func (ps PersistenceService) GetItem(id string) (*Item, error) {
	var item Item
	row := ps.db.QueryRow("SELECT id, name, description, quantity FROM items WHERE id = $1", id)
	err := row.Scan(
		&item.ID,
		&item.Name,
		&item.Descripion,
		&item.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (ps PersistenceService) PostItem(item Item) (string, error) {
	stmt, err := ps.db.Prepare("INSERT INTO public.items (name, description, quantity) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		return "", err
	}
	var id int64
	err = stmt.QueryRow(item.Name, item.Descripion, item.Quantity).Scan(&id)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(id, 10), nil
}

func (ps PersistenceService) PutItem(item Item) (*Item, error) {
	stmt, err := ps.db.Prepare("UPDATE public.items SET name=$1, description=$2, quantity=$3 WHERE id=$4 RETURNING id, name, description, quantity")
	if err != nil {
		return nil, err
	}
	var updatedItem Item
	err = stmt.QueryRow(item.Name, item.Descripion, item.Quantity, item.ID).Scan(
		&updatedItem.ID,
		&updatedItem.Name,
		&updatedItem.Descripion,
		&updatedItem.Quantity,
	)
	if err != nil {
		return nil, err
	}

	return &updatedItem, nil
}
