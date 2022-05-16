package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"

	"github.com/pix303/minimal-rest-api-server/pkg/auth"
	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

// contextKey rappersents specific context key type to override default interface{} type for key in Context
type contextKey struct {
	name string
}

// Handler take care of persistence, auth and session requests
type Handler struct {
	Sessioner   *auth.Sessioner
	ItemService persistence.ItemPersistencer
}

func newHandler(dbdns string) (*Handler, error) {
	ps, err := persistence.NewPostgresqlPersistenceService(dbdns)
	if err != nil {
		return nil, err
	}

	s, err := auth.NewSessionManager(dbdns)
	if err != nil {
		return nil, err
	}

	return &Handler{ItemService: ps, Sessioner: s}, nil
}

// NewRouter return new Router/Multiplex to handler api request endpoint
// secretKey is needed to sign auth token, dbDns is the url for connect dbrms
func NewRouter(dbDns string) (*mux.Router, error) {

	handler, err := newHandler(dbDns)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.Use(handler.Sessioner.StoreManager.Public)
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})
	r.HandleFunc("/", welcomeHandler)
	r.HandleFunc("/auth/{action}/{provider}", handler.loginHandler).Methods("GET")
	r.HandleFunc("/auth/logout", handler.logoutHandler).Methods("GET")

	subr := r.PathPrefix("/api/v1").Subrouter()
	subr.Use(handler.Sessioner.StoreManager.Auth)
	subr.HandleFunc("/", welcomeAuthedHandler).Methods("GET")
	subr.HandleFunc("/items", handler.usersGetHandler).Methods("GET")

	return r, nil
}

func welcomeHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Welcome to minimal web api, you are NOT logged in"))
}

func (h *Handler) loginHandler(rw http.ResponseWriter, rq *http.Request) {
	params := mux.Vars(rq)
	action := params["action"]
	provider := params["provider"]

	switch action {
	case "login":
		_, err := gothic.CompleteUserAuth(rw, rq)
		if err != nil {
			gothic.BeginAuthHandler(rw, rq)
		} else {
			Respond(rw, rq, "you are logged in", http.StatusOK)
		}
		return

	case "callback":
		u, err := gothic.CompleteUserAuth(rw, rq)
		if err != nil {
			RespondError(rw, rq, err, fmt.Sprintf("failed authorization for %s", provider), http.StatusUnauthorized)
			return
		}
		if err != nil {
			RespondError(rw, rq, err, fmt.Sprintf("failed sessioning auth for %s", provider), http.StatusUnauthorized)
			return
		}
		if err != nil {
			RespondError(rw, rq, err, fmt.Sprintf("failed sessioning auth for %s", provider), http.StatusUnauthorized)
			return
		}

		err = h.Sessioner.StoreManager.Init(rw, rq, u.UserID)
		if err != nil {
			RespondError(rw, rq, err, "failed session init", http.StatusUnauthorized)
			return
		}

		http.Redirect(rw, rq, "/api/v1/", http.StatusSeeOther)
		return
	}
}

func (h *Handler) logoutHandler(rw http.ResponseWriter, rq *http.Request) {
	err := gothic.Logout(rw, rq)
	if err != nil {
		RespondError(rw, rq, err, "failed logout", http.StatusUnauthorized)
		return
	}

	err = h.Sessioner.StoreManager.Revoke(rq.Context(), rw)
	if err != nil {
		RespondError(rw, rq, err, "failed session out", http.StatusUnauthorized)
	} else {
		http.Redirect(rw, rq, "/", http.StatusTemporaryRedirect)
	}
}

func welcomeAuthedHandler(rw http.ResponseWriter, rq *http.Request) {
	EncodeBody(rw, rq, "Welcome unknow to minimal web api authenticated")
}

func (s *Handler) usersGetHandler(rw http.ResponseWriter, rq *http.Request) {
	items, err := s.ItemService.GetItems(0, 10)
	if err != nil {
		RespondError(rw, rq, err, "Error on retrive items", http.StatusInternalServerError)
		return
	}

	EncodeBody(rw, rq, items)
}
