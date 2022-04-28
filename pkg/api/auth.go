package api

import (
	"context"

	"github.com/gorilla/sessions"
	abclientstate "github.com/volatiletech/authboss-clientstate"
	"github.com/volatiletech/authboss/v3"
	"github.com/volatiletech/authboss/v3/defaults"

	"github.com/pix303/minimal-rest-api-server/pkg/domain"
)

var (
	ab           = authboss.New()
	database     = NewMemStorer()
	sessionStore abclientstate.SessionStorer
	cookieStore  abclientstate.CookieStorer
)

func SetupAuthConfig() (*authboss.Authboss, error) {
	ab.Config.Paths.RootURL = "http://localhost:8080/"
	ab.Config.Paths.Mount = "/authboss"
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}

	sessionStore = abclientstate.NewSessionStorer("session-name", []byte("session-test"))
	cookieSessionStore := sessionStore.Store.(*sessions.CookieStore)
	cookieSessionStore.MaxAge(60)
	cookieSessionStore.Options.HttpOnly = false
	cookieStore = abclientstate.NewCookieStorer([]byte("cookie-test"), nil)
	cookieStore.HTTPOnly = false

	defaults.SetCore(&ab.Config, true, false)
	if err := ab.Init(); err != nil {
		return nil, err
	}

	return ab, nil

}

type MemStorer struct {
	Users  map[string]domain.User
	Tokens map[string][]string
}

func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users:  make(map[string]domain.User),
		Tokens: make(map[string][]string),
	}
}

func (m MemStorer) Save(_ context.Context, user authboss.User) error {

	u := user.(*domain.User)
	m.Users[u.GetPID()] = *u
	return nil
}

func (m MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {
	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserFound
	}
	return &u, nil
}
