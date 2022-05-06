package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

// ProviderKeys rappresents config to connect to authentication provier
type ProviderKeys struct {
	ClientID     string
	ClientSecret string
	Callback     string
}

// InitOauth inits oauth manager
func InitOauth(providersMap map[string]ProviderKeys, sessionSecret string) {

	for key, pkeys := range providersMap {
		switch key {
		case "github":
			goth.UseProviders(github.New(pkeys.ClientID, pkeys.ClientSecret, pkeys.Callback))
		}
	}

	authProviderSessionStore := sessions.NewCookieStore([]byte(sessionSecret))
	authProviderSessionStore.Options.HttpOnly = true
	authProviderSessionStore.Options.MaxAge = 60
	gothic.Store = authProviderSessionStore
}
