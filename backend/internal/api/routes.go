package api

import (
	"net/http"
	"time"

	"backend/internal/constants"
	"backend/internal/handlers"
	"backend/internal/handlers/auth"
	"backend/internal/handlers/middleware"
	"backend/internal/localization"
	"backend/internal/services"

	"go.uber.org/zap"
)

type RunSecurityMiddlewareParams struct {
	Routes             []Route
	Logger             *zap.Logger
	RateLimiter        *middleware.RateLimiter
	CSRFSecret         string
	JWTSecret          string
	JWTRefreshSecret   string
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
	CorsAllowedOrigins []string
	Services           *services.Services
	ProdEnv            bool
	GoogleAuth         *auth.GoogleAuth
	Localizer          *localization.Localizer
}

// runs all security middleware and then executes route handler, important: it executes in reverse order than applied
func runSecurityMiddleware(mux *http.ServeMux, params *RunSecurityMiddlewareParams) {
	for _, route := range params.Routes {
		handler := route.Handler

		// Innermost (last to execute): Logs handler-level requests
		handler = middleware.DetailedLoggingMiddleware(params.Logger, handler)

		// Localization middleware
		handler = middleware.LocalizerMiddleware(params.Localizer, handler)
		handler = middleware.DetectLanguageMiddleware(handler)

		// Permission middleware (only one permission for a route)
		if route.Options.Permission != "" {
			handler = middleware.PermissionMiddleware(route.Options.Permission, params.Services.UserService, handler)
		}

		// Role middleware (allowing multiple roles)
		if len(route.Options.Roles) > 0 {
			handler = middleware.RoleMiddleware(route.Options.Roles, handler)
		}

		handler = middleware.SecurityHeadersMiddleware(handler)
		handler = middleware.CSRFProtectionMiddleware(params.CSRFSecret, handler)

		if !route.Options.Public {
			handler = middleware.JWTMiddleware(
				params.JWTSecret,
				params.JWTRefreshSecret,
				params.TokenExpiry,
				params.RefreshTokenExpiry,
				params.ProdEnv,
				handler,
			)
		}

		if !route.Options.NoCORS {
			handler = middleware.CORSMiddleware(params.CorsAllowedOrigins, handler)
		}

		handler = middleware.RateLimitMiddleware(params.RateLimiter, handler)
		handler = middleware.SecurityLoggingMiddleware(params.Logger, handler)
		// Outermost (first to execute): Logs rate limit rejections

		mux.Handle(route.Pattern, handler)
	}
}

// RouteOptions defines optional settings for a route
type RouteOptions struct {
	NoCORS     bool
	Public     bool     // If true, this route does not require authentication
	Roles      []string // List of roles allowed to access the route
	Permission string   // Single permission required to access the route
}

// Route struct holds information about a route
type Route struct {
	Pattern string
	Handler http.Handler
	Options RouteOptions
}

// newRoute creates a Route with optional settings
func newRoute(pattern string, handler http.Handler, opts RouteOptions) Route {
	return Route{
		Pattern: pattern,
		Handler: handler,
		Options: opts,
	}
}

type RegisterRoutesParams struct {
	ProdEnv            bool
	Services           *services.Services
	Logger             *zap.Logger
	RateLimiter        *middleware.RateLimiter
	CSRFSecret         string
	JWTSecret          string
	JWTRefreshSecret   string
	TokenExpiry        time.Duration
	RefreshTokenExpiry time.Duration
	CorsAllowedOrigins []string
	GoogleAuth         *auth.GoogleAuth
	Localizer          *localization.Localizer
}

// RegisterRoutes registers all routes with their middleware
func RegisterRoutes(mux *http.ServeMux, params *RegisterRoutesParams) {
	routes := []Route{
		newRoute(constants.WordsRoutes.Words, handlers.WordsHandler(params.Services.WordService), RouteOptions{Public: true}),
		newRoute(constants.AuthRoutes.GoogleLogin, params.GoogleAuth.GoogleLogin(), RouteOptions{Public: true}),
		newRoute(constants.AuthRoutes.GoogleCallback, params.GoogleAuth.GoogleCallback(params.Services.UserService, params.ProdEnv, params.Logger),
			RouteOptions{Public: true}),
		// newRoute("/api/public", handlers.WordsHandler(services.WordService), RouteOptions{Public: true, NoCORS: true}),
		/*newRoute("/api/words/protected", handlers.ProtectedWordsHandler(services.WordService), RouteOptions{
			Roles:      []string{constants.Roles.Admin, constants.Roles.Moderator},
			Permission: "view_protected_words",
		}),*/
	}

	securityMiddlewareParams := &RunSecurityMiddlewareParams{
		Routes:             routes,
		Logger:             params.Logger,
		RateLimiter:        params.RateLimiter,
		CSRFSecret:         params.CSRFSecret,
		JWTSecret:          params.JWTSecret,
		JWTRefreshSecret:   params.JWTRefreshSecret,
		TokenExpiry:        params.TokenExpiry,
		RefreshTokenExpiry: params.RefreshTokenExpiry,
		CorsAllowedOrigins: params.CorsAllowedOrigins,
		Services:           params.Services,
		ProdEnv:            params.ProdEnv,
		GoogleAuth:         params.GoogleAuth,
		Localizer:          params.Localizer,
	}

	runSecurityMiddleware(mux, securityMiddlewareParams)
}
