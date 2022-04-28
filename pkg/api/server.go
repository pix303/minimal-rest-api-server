package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/volatiletech/authboss/v3"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

const apiSessionKey = "api-session-key"
const apiSessionUserKey = "api-session-key-user"

// contextKey rappersents specific context key type to override default interface{} type for key in Context
type contextKey struct {
	name string
}

// PersistenceHandler take care of persistence requests
type PersistenceHandler struct {
	UserService persistence.UserPersistencer
}

// tokenWrapper brings auth token resources during process
type tokenWrapper struct {
	SecretKey []byte
	Source    string
}

var contextKeyUsernameKey = &contextKey{"username"}
var authToken tokenWrapper

func newServer(dbdns string) (*PersistenceHandler, error) {
	ps, err := persistence.NewPostgresqlPersistenceService(dbdns)
	if err != nil {
		return nil, err
	}
	return &PersistenceHandler{UserService: ps}, nil
}

// NewRouter return new Router/Multiplex to handler api request endpoint
// secretKey is needed to sign auth token, dbDns is the url for connect dbrms
func NewRouter(secretKey, dbDns string, ab *authboss.Authboss) (*mux.Router, error) {

	authToken.SecretKey = []byte(secretKey)

	s, err := newServer(dbDns)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.Use(ab.LoadClientStateMiddleware)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})
	r.HandleFunc("/", welcomeHandler)
	r.PathPrefix("/authboss").Handler(http.StripPrefix("/authboss", ab.Config.Core.Router))

	subr := r.PathPrefix("/api/v1").Subrouter()
	subr.Use(authboss.Middleware2(ab, authboss.RequireNone, authboss.RespondUnauthorized))
	subr.HandleFunc("/", welcomeAuthedHandler).Methods("GET")
	subr.HandleFunc("/users", s.usersGetHandler).Methods("GET")

	return r, nil
}

func welcomeHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Welcome to minimal web api, you are not logged in"))
}

func welcomeAuthedHandler(rw http.ResponseWriter, rq *http.Request) {
	usernameRaw := rq.Context().Value(contextKeyUsernameKey)
	fmt.Println(usernameRaw)
	if usernameRaw != nil {
		if username, ok := usernameRaw.(string); ok {
			EncodeBody(rw, rq, fmt.Sprintf("Welcome %s to minimal web api authenticated", username))
		}
	}
}

func (s *PersistenceHandler) usersGetHandler(rw http.ResponseWriter, rq *http.Request) {
	users, err := s.UserService.GetUsers()
	if err != nil {
		RespondError(rw, rq, err, "Error on retrive users", http.StatusInternalServerError)
		return
	}

	EncodeBody(rw, rq, users)
}
