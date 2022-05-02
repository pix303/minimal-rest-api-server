package auth

import (
	"net/http"
	"time"

	"github.com/swithek/sessionup"
	pgstore "github.com/swithek/sessionup-pgstore"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

// Sessioner wraps session store manager
type Sessioner struct {
	StoreManager *sessionup.Manager
}

// NewSessionManager build and return db session store manager
func NewSessionManager(dbdns string) (*Sessioner, error) {
	db, err := persistence.GetDBInstance(dbdns)
	if err != nil {
		return nil, err
	}

	s, err := pgstore.New(db, "sessions", 60*time.Second)
	if err != nil {
		return nil, err
	}

	mgr := sessionup.NewManager(s, sessionup.SameSite(http.SameSiteNoneMode))
	sm := Sessioner{mgr}
	return &sm, nil
}
