package auth

import (
	"time"

	"github.com/swithek/sessionup"
	pgstore "github.com/swithek/sessionup-pgstore"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

type Sessioner struct {
	StoreManager *sessionup.Manager
}

func NewSessionManager(dbdns string) (*Sessioner, error) {
	db, err := persistence.GetDBInstance(dbdns)
	if err != nil {
		return nil, err
	}
	s, err := pgstore.New(db, "sessions", time.Hour*1)
	if err != nil {
		return nil, err
	}

	mgr := sessionup.NewManager(s)
	sm := Sessioner{mgr}
	return &sm, nil
}
