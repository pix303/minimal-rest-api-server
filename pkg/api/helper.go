package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// decodeBody helper and abstract func to parse payload request
func DecodeBody(rq *http.Request, payload interface{}) error {
	defer rq.Body.Close()
	return json.NewDecoder(rq.Body).Decode(payload)
}

// encodeBody helper and abstract func to responde with data
func EncodeBody(rw http.ResponseWriter, rq *http.Request, data interface{}) error {
	return json.NewEncoder(rw).Encode(data)
}

// respond helper func for build a response
func Respond(rw http.ResponseWriter, rq *http.Request, data interface{}, status int) {
	rw.WriteHeader(status)
	if data != nil {
		EncodeBody(rw, rq, data)
	}
}

// respondError helper func for build a custom message error response
func RespondError(rw http.ResponseWriter, rq *http.Request, err any, errMessage string, status int) {
	msg := fmt.Sprintf("%d - %s: %v", status, errMessage, err)
	Respond(rw, rq, msg, status)
}

// respondHTTPErr helper func for build a generic error response
func RespondHTTPErr(rw http.ResponseWriter, rq *http.Request, status int) {
	Respond(rw, rq, fmt.Sprintf("%d - %s", status, http.StatusText(status)), status)
}

// func GenerateToken(subject string, duration int) string {
// 	claims := &jwt.StandardClaims{
// 		ExpiresAt: time.Now().Add(time.Second * time.Duration(duration)).Unix(),
// 		Issuer:    "minimal-rest-api",
// 		Subject:   subject,

// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedTokenString, _ := token.SignedString([]byte(authToken.SecretKey))

// 	return signedTokenString
// }
