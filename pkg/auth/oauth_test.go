package auth

import (
	"testing"

	"github.com/markbates/goth"
	"github.com/stretchr/testify/assert"
)

func TestInitOauth(t *testing.T) {
	providers := map[string]ProviderKeys{
		"github": {ClientID: "123", ClientSecret: "123", Callback: "http://localhost:8080/api/callback"},
	}

	InitOauth(providers, "123")

	p := goth.GetProviders()["github"]
	assert.NotNil(t, p)

	p = goth.GetProviders()["facebook"]
	assert.Nil(t, p)

}
