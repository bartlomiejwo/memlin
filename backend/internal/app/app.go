package app

import (
	"net/http"
	"time"

	"backend/internal/api"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers/auth"
	"backend/internal/handlers/middleware"
	"backend/internal/localization"
	"backend/internal/repositories"
	"backend/internal/services"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type App struct {
	Store              *db.Store
	Logger             *zap.Logger
	Repositories       *repositories.Repositories
	Services           *services.Services
	ProdEnv            bool
	RateLimiter        *middleware.RateLimiter
	CSRFSecret         string
	JWTSecret          string
	JWTRefreshSecret   string
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
	GoogleAuth         *auth.GoogleAuth
	CorsAllowedOrigins []string
	Localizer          *localization.Localizer
}

func New(store *db.Store, cfg *config.Config, logger *zap.Logger) *App {
	repos := repositories.InitRepositories(store, logger)
	services := services.InitServices(repos, logger)

	// Initialize Localizer
	localizer := localization.NewLocalizer(logger)

	// Convert requests per minute to rate.Every durations
	// Formula: duration (ms) = 60,000 ms / requests per minute
	ipRateDuration := time.Duration(60000/cfg.IPRateLimit) * time.Millisecond     // e.g., 60,000 / 300 = 200ms
	userRateDuration := time.Duration(60000/cfg.UserRateLimit) * time.Millisecond // e.g., 60,000 / 60 = 1000ms

	// IP: ~300 req/min, burst 30; User: ~60 req/min, burst 15
	rateLimiter := middleware.NewRateLimiter(
		rate.Every(ipRateDuration), cfg.IPBurst, // IP rate and burst
		rate.Every(userRateDuration), cfg.UserBurst, // User rate and burst
		cfg.JWTSecretKey,
		cfg.UseUserRate, // User based limiting
	)

	googleAuth := auth.NewGoogleAuth(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.JWTSecretKey,
		cfg.JWTRefreshSecretKey,
		time.Duration(cfg.JWTExpiryMinutes)*time.Minute,
		time.Duration(cfg.JWTRefreshExpiryDays)*time.Hour*24,
	)

	return &App{
		Store:              store,
		Logger:             logger,
		Repositories:       repos,
		Services:           services,
		ProdEnv:            cfg.ProdEnv,
		RateLimiter:        rateLimiter,
		CSRFSecret:         cfg.CSRFSecretKey,
		JWTSecret:          cfg.JWTSecretKey,
		JWTRefreshSecret:   cfg.JWTRefreshSecretKey,
		TokenExpiry:        googleAuth.TokenExpiry,
		RefreshTokenExpiry: googleAuth.RefreshTokenExpiry,
		GoogleAuth:         googleAuth,
		CorsAllowedOrigins: cfg.CORSAllowedOrigins,
		Localizer:          localizer,
	}
}

func (app *App) RegisterRoutes(mux *http.ServeMux) {
	params := &api.RegisterRoutesParams{
		ProdEnv:            app.ProdEnv,
		Services:           app.Services,
		Logger:             app.Logger,
		RateLimiter:        app.RateLimiter,
		CSRFSecret:         app.CSRFSecret,
		JWTSecret:          app.JWTSecret,
		JWTRefreshSecret:   app.JWTRefreshSecret,
		TokenExpiry:        app.TokenExpiry,
		RefreshTokenExpiry: app.RefreshTokenExpiry,
		CorsAllowedOrigins: app.CorsAllowedOrigins,
		GoogleAuth:         app.GoogleAuth,
		Localizer:          app.Localizer,
	}

	api.RegisterRoutes(mux, params)
}
