package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

type ProviderKeys struct {
	ClientID     string
	ClientSecret string
	Callback     string
}

func InitOauth(providersMap map[string]ProviderKeys, sessionSecret string) {

	for k, pkeys := range providersMap {
		switch k {
		case "github":
			goth.UseProviders(github.New(pkeys.ClientID, pkeys.ClientSecret, pkeys.Callback))
		}
	}

	authProviderSessionStore := sessions.NewCookieStore([]byte(sessionSecret))
	authProviderSessionStore.Options.HttpOnly = true
	authProviderSessionStore.Options.MaxAge = 60
	gothic.Store = authProviderSessionStore

}
