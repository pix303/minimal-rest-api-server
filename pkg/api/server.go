package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

// contextKey rappersents specific context key type to override default interface{} type for key in Context
type contextKey struct {
	name string
}

// PersistenceHandler take care of persistence requests
type PersistenceHandler struct {
	UserService persistence.UserPersistencer
}

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
func NewRouter(secretKey string, dbDns string) (*mux.Router, error) {

	authToken.SecretKey = []byte(secretKey)

	s, err := newServer(dbDns)
	if err != nil {
		return nil, err
	}

	// only for debug
	log.Info().Msg(GenerateTokenForDebug(3600))

	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", welcomeHandler)

	subr := r.PathPrefix("/api/v1").Subrouter()
	subr.Use(JWTValidatorMiddleware)
	subr.HandleFunc("/", welcomeAuthedHandler)
	subr.HandleFunc("/users", s.usersGetHandler)

	return r, nil
}

func welcomeHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Welcome to minimal web api"))
}

func JWTValidatorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {

		bearerHead := rq.Header.Get("Authorization")
		authToken.Source = strings.Split(bearerHead, " ")[1]

		claims := &jwt.StandardClaims{}
		parsedToken, err := jwt.ParseWithClaims(authToken.Source, claims, func(t *jwt.Token) (interface{}, error) {

			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
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
	username := rq.Context().Value(contextKeyUsernameKey).(string)
	rw.Write([]byte(fmt.Sprintf("Welcome %s to minimal web api authenticated", username)))
}

func (s *PersistenceHandler) usersGetHandler(rw http.ResponseWriter, rq *http.Request) {
	users, err := s.UserService.GetUsers()
	if err != nil {
		RespondError(rw, rq, err, "Error on retrive users", http.StatusInternalServerError)
		return
	}
	userJson, err := json.Marshal(users)
	if err != nil {
		RespondError(rw, rq, err, "Error on json encoding user", http.StatusInternalServerError)
		return
	}
	rw.Write([]byte(userJson))
}
