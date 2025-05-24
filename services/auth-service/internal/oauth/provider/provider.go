package provider

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
)

type CbResponse struct {
	Email, Name string
}
type OauthProviderType string

var (
	OauthProviderTypeGithub OauthProviderType = "github"
)
var OauthProviderRegistry map[OauthProviderType]OauthProvider

type UserInfo struct {
	Name, Email string
	RoleId      store.DefaultRoleID
}
type LoginOption struct {
	RoleId store.DefaultRoleID
}
type OauthProvider interface {
	Login(w http.ResponseWriter, r *http.Request, opt *LoginOption)
	Callback(w http.ResponseWriter, r *http.Request) (*UserInfo, error)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)
}
func storeStateInCookie(w http.ResponseWriter, provider OauthProviderType, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     fmt.Sprintf("oauth_state_%s", provider),
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   300,
	})
}
func getStateFromCookie(r *http.Request, provider OauthProviderType) (string, error) {
	cookie, err := r.Cookie(fmt.Sprintf("oauth_state_%s", provider))
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Flow
// implement oauth with github
// create your own oauth provider = client secret, id, and then test flow
// create diagram and proceed with post
// share post
// create 2 more posts
// apply
// http://localhost:3000/auth/github/callback
