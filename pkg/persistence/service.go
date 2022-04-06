package persistence

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

type PersistenceService struct {
	db *bun.DB
}

// NewPersistenceService manage dbrms connecton and requests
func NewPersistenceService(dbdns string) (*PersistenceService, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbdns)))
	db := bun.NewDB(sqldb, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return nil, err
	}

	service := PersistenceService{db}

	return &service, nil
}

// GetUsers retrive all users
func (ps PersistenceService) GetUsers(ctx context.Context) ([]domain.User, error) {
	users := make([]domain.User, 0)
	err := ps.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
