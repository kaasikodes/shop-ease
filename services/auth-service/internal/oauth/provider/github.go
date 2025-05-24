package provider

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kaasikodes/shop-ease/services/auth-service/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubOauthProvider struct {
	config *oauth2.Config
}

func NewGithubOauthProvider(clientId, clientSecret, redirectUrl string) *GithubOauthProvider {
	return &GithubOauthProvider{
		config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     github.Endpoint,
			RedirectURL:  redirectUrl,
			Scopes:       []string{"user:email", "read:user"},
		},
	}

}
func (h *GithubOauthProvider) Login(w http.ResponseWriter, r *http.Request, opt *LoginOption) {
	// e.g., encode roleId into state like base64("roleId:randomstring")
	roleId := strconv.Itoa(int(opt.RoleId))

	state := base64.URLEncoding.EncodeToString([]byte(roleId + ":" + generateRandomString(32)))
	storeStateInCookie(w, OauthProviderTypeGithub, state)
	url := h.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}
func (h *GithubOauthProvider) Callback(w http.ResponseWriter, r *http.Request) (*UserInfo, error) {
	ctx := r.Context()
	stateFromQuery := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	stateFromCookie, err := getStateFromCookie(r, OauthProviderTypeGithub)
	if err != nil || stateFromQuery != stateFromCookie {
		return nil, errors.New("query state does not match cookie state")
	}
	if code == "" {
		return nil, errors.New("missing code")
	}

	token, err := h.config.Exchange(r.Context(), code)
	if err != nil {
		return nil, errors.Join(err, errors.New("token exchange failed"))

	}
	client := h.config.Client(ctx, token)
	respU, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to get user info"))

	}
	defer respU.Body.Close()
	var wg sync.WaitGroup
	var userEmail string
	var errs []error

	wg.Add(1)
	go func() {
		defer wg.Done()

		var localUserEmail string
		var localErrs []error

		respE, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			localErrs = append(localErrs, err, errors.New("failed to get user emails"))
			errs = append(errs, localErrs...)

			return
		}
		defer respE.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(respE.Body).Decode(&emails); err != nil {
			localErrs = append(localErrs, err, errors.New("failed to parse user emails"))
		} else {
			for _, e := range emails {
				if e.Primary && e.Verified {
					localUserEmail = e.Email
					break
				}
			}
		}

		// Assign to outer variables in a thread-safe manner
		userEmail = localUserEmail
		errs = append(errs, localErrs...)
	}()
	wg.Wait()

	wg.Wait()

	body, _ := io.ReadAll(respU.Body)

	var userInfo GithubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil || len(errs) > 0 {
		errs = append(errs, err, errors.New("failed to parse user info"))
		return nil, errors.Join(errs...)

	}
	log.Println(errs, err, userInfo, userEmail, respU)
	decodedState, _ := base64.URLEncoding.DecodeString(stateFromQuery)
	parts := strings.SplitN(string(decodedState), ":", 2)
	_roleId := parts[0] // "roleId"
	roleId, err := strconv.Atoi(_roleId)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Name:   userInfo.Name,
		Email:  userEmail,
		RoleId: store.DefaultRoleID(roleId),
	}, nil
}

type GithubUserInfo struct {
	AvatarURL               string     `json:"avatar_url"`
	Bio                     string     `json:"bio"`
	Blog                    string     `json:"blog"`
	Collaborators           int        `json:"collaborators"`
	Company                 string     `json:"company"`
	CreatedAt               time.Time  `json:"created_at"`
	DiskUsage               int        `json:"disk_usage"`
	Email                   *string    `json:"email"`
	EventsURL               string     `json:"events_url"`
	Followers               int        `json:"followers"`
	FollowersURL            string     `json:"followers_url"`
	Following               int        `json:"following"`
	FollowingURL            string     `json:"following_url"`
	GistsURL                string     `json:"gists_url"`
	GravatarID              string     `json:"gravatar_id"`
	Hireable                bool       `json:"hireable"`
	HTMLURL                 string     `json:"html_url"`
	ID                      int        `json:"id"`
	Location                string     `json:"location"`
	Login                   string     `json:"login"`
	Name                    string     `json:"name"`
	NodeID                  string     `json:"node_id"`
	NotificationEmail       *string    `json:"notification_email"`
	OrganizationsURL        string     `json:"organizations_url"`
	OwnedPrivateRepos       int        `json:"owned_private_repos"`
	Plan                    GithubPlan `json:"plan"`
	PrivateGists            int        `json:"private_gists"`
	PublicGists             int        `json:"public_gists"`
	PublicRepos             int        `json:"public_repos"`
	ReceivedEventsURL       string     `json:"received_events_url"`
	ReposURL                string     `json:"repos_url"`
	SiteAdmin               bool       `json:"site_admin"`
	StarredURL              string     `json:"starred_url"`
	SubscriptionsURL        string     `json:"subscriptions_url"`
	TotalPrivateRepos       int        `json:"total_private_repos"`
	TwitterUsername         *string    `json:"twitter_username"`
	TwoFactorAuthentication bool       `json:"two_factor_authentication"`
	Type                    string     `json:"type"`
	UpdatedAt               time.Time  `json:"updated_at"`
	URL                     string     `json:"url"`
	UserViewType            string     `json:"user_view_type"`
}

type GithubPlan struct {
	Collaborators int    `json:"collaborators"`
	Name          string `json:"name"`
	PrivateRepos  int    `json:"private_repos"`
	Space         int    `json:"space"`
}
