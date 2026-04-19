package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/constants"
	"backend/internal/handlers/middleware"
	loc "backend/internal/localization"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuth contains OAuth configuration
type GoogleAuth struct {
	Config             *oauth2.Config
	JWTSecret          string
	JWTRefreshSecret   string
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
}

// NewGoogleAuth initializes Google OAuth configuration
func NewGoogleAuth(clientID, clientSecret string, jwtSecret, jwtRefreshSecret string, tokenExpiry, refreshTokenExpiry time.Duration) *GoogleAuth {
	return &GoogleAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       []string{"email", "profile"},
			Endpoint:     google.Endpoint,
		},
		JWTSecret:          jwtSecret,
		JWTRefreshSecret:   jwtRefreshSecret,
		TokenExpiry:        tokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
	}
}

// GoogleLogin redirects the user to Google's OAuth 2.0 consent page
func (g *GoogleAuth) GoogleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read redirect_uri from query parameters
		redirectURI := r.URL.Query().Get("redirect_uri")
		if redirectURI == "" {
			lang := middleware.GetLanguage(r)
			l := middleware.GetLocalizer(r)
			http.Error(w, l.T(lang, loc.AuthKeys.MissingRedirectURI, nil), http.StatusBadRequest)
			return
		}

		// Set the redirect URI dynamically
		conf := g.Config
		conf.RedirectURL = redirectURI

		// Generate OAuth URL
		url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// GoogleCallback handles the callback from Google after authentication
func (g *GoogleAuth) GoogleCallback(userService *services.UserService, prodEnv bool, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := middleware.GetLanguage(r)
		l := middleware.GetLocalizer(r)

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, l.T(lang, loc.AuthKeys.MissingAuthCode, nil), http.StatusBadRequest)
			return
		}

		// Exchange authorization code for access token
		token, err := g.Config.Exchange(context.Background(), code)
		if err != nil {
			logger.Error("Failed to exchange token", zap.Error(err))
			http.Error(w, l.T(lang, loc.AuthKeys.AuthFail, nil), http.StatusInternalServerError)
			return
		}

		// Fetch user info from Google
		client := g.Config.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			logger.Error("Failed to get user info", zap.Error(err))
			http.Error(w, l.T(lang, loc.AuthKeys.UserInfoRetrieveFail, nil), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var googleUser struct {
			ID      string `json:"id"`
			Email   string `json:"email"`
			Name    string `json:"name"`
			Picture string `json:"picture"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			http.Error(w, l.T(lang, loc.AuthKeys.UserInfoParseFail, nil), http.StatusInternalServerError)
			return
		}

		// Create or update the user in the database
		user := models.User{
			Email:  googleUser.Email,
			Name:   googleUser.Name,
			Avatar: googleUser.Picture,
			Role:   constants.Roles.User,
			LinkedAccounts: []models.LinkedAccount{
				{
					Provider:   constants.AuthProviders.Google,
					ProviderID: googleUser.ID,
				},
			},
		}

		createdUser, err := userService.CreateOrUpdateUser(r.Context(), user)
		if err != nil {
			http.Error(w, l.T(lang, loc.AuthKeys.UserCreateUpdateFail, nil), http.StatusInternalServerError)
			return
		}

		jwtToken, err := utils.GenerateJWT(g.JWTSecret, createdUser.ID, createdUser.Role, g.TokenExpiry)
		if err != nil {
			http.Error(w, l.T(lang, loc.AuthKeys.JWTCreateFail, nil), http.StatusInternalServerError)
			return
		}

		refreshToken, err := utils.GenerateJWTRefresh(g.JWTRefreshSecret, createdUser.ID, g.RefreshTokenExpiry)
		if err != nil {
			http.Error(w, l.T(lang, loc.AuthKeys.RefreshTokenCreateFail, nil), http.StatusInternalServerError)
			return
		}

		// Set the refresh token as HttpOnly cookie
		utils.SetRefreshTokenCookie(w, refreshToken, g.RefreshTokenExpiry, prodEnv)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token": jwtToken,
		})
	}
}
