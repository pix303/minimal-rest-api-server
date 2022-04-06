package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"

	"github.com/pix303/minimal-rest-api-server/pkg/persistence"
)

// contextKey rappersents specific context key type to override default interface{} type for key in Context
type contextKey struct {
	name string
}

// Server rappresents components that compose app as Service
type Server struct {
	Service *persistence.PersistenceService
}

var contextKeyUsernameKey = &contextKey{"username"}

var authToken *jwtauth.JWTAuth

func newServer(dbdns string) (*Server, error) {
	ps, err := persistence.NewPersistenceService(dbdns)
	if err != nil {
		return nil, err
	}
	return &Server{Service: ps}, nil
}

// NewRouter return new Router/Multiplex to handler api request endpoint
// secretKey is needed to sign auth token, dbDns is the url for connect dbrms
func NewRouter(secretKey string, dbDns string) (*chi.Mux, error) {

	s, err := newServer(dbDns)
	if err != nil {
		return nil, err
	}

	authToken = jwtauth.New("HS256", []byte(secretKey), nil)

	// only for debug-------------------------------------------------------
	_, ts, err := authToken.Encode(map[string]interface{}{
		"iss":  "minimal-api",
		"sub":  "pix303",
		"name": "Paolo Carraro",
		"exp":  time.Now().Add(time.Second * time.Duration(120)).Unix(),
	})
	log.Println(ts)
	//----------------------------------------------------------------------

	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", welcomeHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(jwtauth.Verifier(authToken))
		r.Use(JWTSubjectExtractorMiddelware)
		r.Get("/", welcomeAuthedHandler)
		r.Get("/users", s.usersGetHandler)
	})

	return r, nil
}

func welcomeHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Write([]byte("Welcome to minimal web api"))
}

func JWTSubjectExtractorMiddelware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {

		err := rq.Context().Value(jwtauth.ErrorCtxKey)
		if err != nil {
			RespondError(rw, rq, err, "Error on verify token", http.StatusUnauthorized)
			return
		}

		token := rq.Context().Value(jwtauth.TokenCtxKey).(jwt.Token)

		name, ok := token.Get("sub")
		if ok {
			ctx := context.WithValue(rq.Context(), contextKeyUsernameKey, name)
			next.ServeHTTP(rw, rq.WithContext(ctx))
			return
		}
		next.ServeHTTP(rw, rq)
	})
}

func welcomeAuthedHandler(rw http.ResponseWriter, rq *http.Request) {
	username := rq.Context().Value(contextKeyUsernameKey).(string)
	rw.Write([]byte(fmt.Sprintf("Welcome %s to minimal web api authenticated", username)))
}

func (s *Server) usersGetHandler(rw http.ResponseWriter, rq *http.Request) {
	ctx := rq.Context()
	users, err := s.Service.GetUsers(ctx)
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
