package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	r.HandleFunc("/ping", welcomeHandler)
	r.HandleFunc("/auth/{action}/{provider}", handler.loginHandler).Methods("GET")
	r.HandleFunc("/auth/logout", handler.logoutHandler).Methods("GET")

	subr := r.PathPrefix("/api/v1").Subrouter()
	subr.Use(handler.Sessioner.StoreManager.Auth)
	subr.HandleFunc("/", welcomeAuthedHandler).Methods("GET")
	subr.HandleFunc("/items", handler.getItemsHandler).Methods("GET")
	subr.HandleFunc("/items", handler.postItemHandler).Methods("POST")
	subr.HandleFunc("/items", handler.putItemHandler).Methods("PUT")
	subr.HandleFunc("/items/{id}", handler.getItemHandler).Methods("GET")

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

func (s *Handler) getItemsHandler(rw http.ResponseWriter, rq *http.Request) {
	items, err := s.ItemService.GetItems(0, 10)
	if err != nil {
		RespondError(rw, rq, err, "Error on retrive items", http.StatusInternalServerError)
		return
	}

	EncodeBody(rw, rq, items)
}

func (s *Handler) getItemHandler(rw http.ResponseWriter, rq *http.Request) {
	id := mux.Vars(rq)["id"]
	item, err := s.ItemService.GetItem(id)
	if err != nil {
		RespondError(rw, rq, err, fmt.Sprintf("Error on retrive item for id %s", id), http.StatusNotFound)
		return
	}

	EncodeBody(rw, rq, item)
}

func (s *Handler) postItemHandler(rw http.ResponseWriter, rq *http.Request) {
	// for {
	// 	var chunk []byte = make([]byte, 10)
	// 	_, err := rq.Body.Read(chunk)
	// 	candidateSource = append(candidateSource, chunk...)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			break
	// 		} else {
	// 			if err != nil {
	// 				RespondError(rw, rq, err, "error on read body to post", http.StatusBadRequest)
	// 				return
	// 			}
	// 		}
	// 	}
	// }

	candidateSource, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		RespondError(rw, rq, err, "error on read body to post", http.StatusBadRequest)
		return
	}

	var candidateItem persistence.Item
	err = json.Unmarshal(candidateSource, &candidateItem)
	if err != nil {
		RespondError(rw, rq, err, "error on read json to post", http.StatusBadRequest)
		return
	}

	id, err := s.ItemService.PostItem(candidateItem)
	if err != nil {
		RespondError(rw, rq, err, "error on persist body to post", http.StatusInternalServerError)
		return
	}
	EncodeBody(rw, rq, id)

}

func (s *Handler) putItemHandler(rw http.ResponseWriter, rq *http.Request) {

	candidateSource, err := ioutil.ReadAll(rq.Body)
	if err != nil {
		RespondError(rw, rq, err, "error on read body to post", http.StatusBadRequest)
		return
	}

	var candidateItem persistence.Item
	err = json.Unmarshal(candidateSource, &candidateItem)
	if err != nil {
		RespondError(rw, rq, err, "error on read json to post", http.StatusBadRequest)
		return
	}

	updatedItem, err := s.ItemService.PutItem(candidateItem)
	if err != nil {
		RespondError(rw, rq, err, "error on persist body to post", http.StatusInternalServerError)
		return
	}
	EncodeBody(rw, rq, updatedItem)

}
