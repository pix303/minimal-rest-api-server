package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

const apiSessionKey = "api-session-key"
const apiSessionUserKey = "api-session-key-user"

var sessionStore *sessions.CookieStore

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
func NewRouter(secretKey, authSessionSecretKey, sessionSecretKey, dbDns string) (*mux.Router, error) {

	authToken.SecretKey = []byte(secretKey)

	sessionStore = sessions.NewCookieStore([]byte(sessionSecretKey))
	sessionStore.Options.MaxAge = 20
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Path = "/"

	authProviderSessionStore := sessions.NewCookieStore([]byte(authSessionSecretKey))
	authProviderSessionStore.Options.HttpOnly = true
	authProviderSessionStore.Options.MaxAge = 60
	gothic.Store = authProviderSessionStore

	s, err := newServer(dbDns)
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", welcomeHandler)
	r.HandleFunc("/auth/{action}/{provider}", loginHandler).Methods("GET")

	subr := r.PathPrefix("/api/v1").Subrouter()
	//subr.Use(JWTValidatorMiddleware)
	subr.Use(authMiddleware)
	subr.HandleFunc("/", welcomeAuthedHandler).Methods("GET")
	subr.HandleFunc("/users", s.usersGetHandler).Methods("GET")

	return r, nil
}

func welcomeHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Welcome to minimal web api, you are not logged in"))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		userSession, err := sessionStore.Get(rq, apiSessionKey)
		if err != nil {
			RespondHTTPErr(rw, rq, http.StatusUnauthorized)
			return
		}
		if userSession.IsNew {
			RespondHTTPErr(rw, rq, http.StatusUnauthorized)
			return
		}

		userEmail := userSession.Values[apiSessionUserKey].(string)
		ctx := context.WithValue(rq.Context(), contextKeyUsernameKey, userEmail)
		next.ServeHTTP(rw, rq.WithContext(ctx))
	})
}

func loginHandler(rw http.ResponseWriter, rq *http.Request) {
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

		userSession, err := sessionStore.Get(rq, apiSessionKey)
		if err != nil {
			RespondError(rw, rq, err, "fail on create session", http.StatusInternalServerError)
		}
		userSession.Values[apiSessionUserKey] = u.Email
		userSession.Save(rq, rw)
		fmt.Println(u.Email)
		http.Redirect(rw, rq, "/api/v1/", http.StatusSeeOther)

	case "logout":
		err := gothic.Logout(rw, rq)
		if err != nil {
			RespondError(rw, rq, err, fmt.Sprintf("failed logout for %s", provider), http.StatusUnauthorized)
		} else {
			http.Redirect(rw, rq, "/", http.StatusPermanentRedirect)
		}
	}
}

func JWTValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {

		bearerHead := rq.Header.Get("Authorization")
		if bearerHead == "" {
			RespondHTTPErr(rw, rq, http.StatusUnauthorized)
			return
		}
		authToken.Source = strings.Split(bearerHead, " ")[1]

		claims := &jwt.StandardClaims{}
		parsedToken, err := jwt.ParseWithClaims(authToken.Source, claims, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return authToken.SecretKey, nil
		})

		if err != nil {
			RespondError(rw, rq, err, "Error on parse JWT", http.StatusUnauthorized)
			return
		}

		if !parsedToken.Valid {
			RespondError(rw, rq, err, "Error on valid JWT", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(rq.Context(), contextKeyUsernameKey, claims.Subject)
		next.ServeHTTP(rw, rq.WithContext(ctx))
	})
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
