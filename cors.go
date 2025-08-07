// Package cors implements Cross-Origin Resource Sharing (CORS) middleware for HTTP servers.
//
// It provides options to customize allowed origins, methods, headers, and the cache duration
// for preflight responses. This middleware simplifies enabling CORS in web applications while
// adhering to the CORS specification. For production use, ensure that the default configuration
// is properly customized to restrict access to trusted origins.
//
// For more information about CORS, visit:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
package cors

import (
	"cmp"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// Options defines the configuration for CORS (Cross-Origin Resource Sharing) middleware.
// CORS allows web applications running at one domain to access resources from another domain.
type Options struct {
	// AllowOrigins specifies which origins (domains) are allowed to make cross-origin requests.
	// Each origin should be in the format "https://example.com" or "http://localhost:3000".
	// Use "*" to allow all origins (not recommended for production with credentials).
	// Use WithDefaultOptions() to customize origins while keeping other defaults:
	// WithDefaultOptions(Options{AllowOrigins: []string{"https://yourdomain.com"}})
	AllowOrigins []string

	// AllowMethods specifies which HTTP methods are allowed for cross-origin requests.
	// Common methods include GET, POST, PUT, PATCH, DELETE, OPTIONS.
	// OPTIONS is typically required for preflight requests.
	AllowMethods []string

	// AllowHeaders specifies which headers the client is allowed to use during the actual request.
	// This is used in response to preflight requests that include Access-Control-Request-Headers.
	// Common headers include "Content-Type", "Authorization", "X-Requested-With".
	AllowHeaders []string

	// MaxAge specifies how long (in seconds) the results of a preflight request can be cached.
	// This reduces the number of preflight requests for subsequent requests from the same origin.
	// A value of 0 means no caching, negative values are invalid.
	MaxAge int
}

const (
	// defaultMaxAge is the default cache duration for preflight requests (60 minutes).
	// This is a reasonable default that balances performance with security.
	defaultMaxAge = 3600
	wildCard      = "*"
)

// DefaultOptions provides a permissive CORS configuration suitable for development.
// For production use, you should customize these settings based on your security requirements.
//
// Default Values:
//   - AllowOrigins: "*" permits requests from any origin
//     ⚠️  WARNING: This is permissive and should be restricted in production to specific domains
//     Use WithDefaultOptions() to customize origins while keeping other defaults:
//     WithDefaultOptions(Options{AllowOrigins: []string{"https://yourdomain.com"}})
//   - AllowMethods: Includes all common HTTP methods for RESTful APIs:
//     GET (retrieve data), POST (create resources), PATCH (partial updates),
//     DELETE (remove resources), PUT (replace resources), OPTIONS (preflight requests)
//   - AllowHeaders: Permits commonly used headers:
//     "Content-Type" (for JSON/XML/form data), "Authorization" (for auth tokens)
//     You may need to add custom headers like "X-API-Key", "X-Requested-With", etc.
//   - MaxAge: Cache preflight responses for 50 minutes (3000 seconds)
//     This reduces network overhead by allowing browsers to reuse preflight responses
var DefaultOptions = Options{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodPut, http.MethodOptions},
	AllowHeaders: []string{"Content-Type", "Authorization"},
	MaxAge:       defaultMaxAge,
}

// WithDefaultOptions populates empty fields in the provided Options with values from DefaultOptions.
// This function is useful when you want to customize only specific CORS settings while keeping
// sensible defaults for the rest.
//
// Fields are considered "empty" and will be populated with defaults if:
//   - AllowOrigins: slice is nil (empty slice [] is valid and preserved)
//   - AllowMethods: slice is nil (empty slice [] is valid and preserved)
//   - AllowHeaders: slice is nil (empty slice [] is valid and preserved)
//   - MaxAge: value is 0 or negative
//
// Parameters:
//   - opt: The Options struct to populate with defaults
//
// Returns:
//   - Options: A new Options struct with empty fields filled from DefaultOptions
//
// Usage Example:
//
//	customOptions := Options{
//		AllowOrigins: []string{"https://myapp.com"},
//		// Other fields will be populated with defaults
//	}
//	finalOptions := WithDefaultOptions(customOptions)
func WithDefaultOptions(opt Options) Options {
	result := opt

	// Populate AllowOrigins if nil (empty slice is valid)
	if result.AllowOrigins == nil {
		result.AllowOrigins = DefaultOptions.AllowOrigins
	}

	// Populate AllowMethods if nil (empty slice is valid)
	if result.AllowMethods == nil {
		result.AllowMethods = DefaultOptions.AllowMethods
	}

	// Populate AllowHeaders if nil (empty slice is valid)
	if result.AllowHeaders == nil {
		result.AllowHeaders = DefaultOptions.AllowHeaders
	}

	// Populate MaxAge if zero or negative
	if result.MaxAge <= 0 {
		result.MaxAge = DefaultOptions.MaxAge
	}

	return result
}

// setAllowOriginHeader sets the Access-Control-Allow-Origin header correctly according to CORS spec.
// The CORS specification only allows a single origin or "*", not comma-separated multiple origins.
// For multiple allowed origins, we check the request's Origin header and return that specific
// origin if it's in our allowed list.
func setAllowOriginHeader(w http.ResponseWriter, r *http.Request, allowedOrigins []string) {
	requestOrigin := r.Header.Get("Origin")
	if len(allowedOrigins) == 0 || requestOrigin == "" {
		return
	}

	if slices.Contains(allowedOrigins, wildCard) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	// Check if the request origin is in our allowed list
	if slices.Contains(allowedOrigins, requestOrigin) {
		w.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		return
	}

	// Request origin not in allowed list, don't set the header (request will be blocked)
}

// CorsMiddleware creates HTTP middleware that handles Cross-Origin Resource Sharing (CORS).
// CORS is a security feature implemented by web browsers that blocks requests from one domain
// to another unless the server explicitly allows it. This middleware adds the necessary
// headers to responses to enable cross-origin requests based on the provided options.
//
// Parameters:
//   - next: The next HTTP handler in the middleware chain
//   - opt: CORS configuration options (use DefaultOptions for development)
//
// Returns:
//   - http.Handler: Middleware that can be used with any HTTP router/mux
//
// Usage Example:
//
//	handler := CorsMiddleware(yourHandler, DefaultOptions)
//	http.ListenAndServe(":8080", handler)
func CorsMiddleware(next http.Handler, opt Options) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Convert MaxAge to string, using default if not specified or zero
		maxAge := strconv.Itoa(cmp.Or(opt.MaxAge, defaultMaxAge))

		// Set CORS headers that apply to all requests
		// Access-Control-Allow-Origin: Must be a single origin or "*" per CORS specification
		// Handle multiple allowed origins by checking the request's Origin header
		setAllowOriginHeader(w, r, opt.AllowOrigins)

		// Vary: Origin tells caches that the response varies based on the Origin header
		// This prevents cached responses from being served to different origins
		w.Header().Set("Vary", "Origin")

		// Access-Control-Allow-Methods: Lists methods allowed for cross-origin requests
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(opt.AllowMethods, ", "))

		// Access-Control-Allow-Headers: Lists headers allowed in cross-origin requests
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(opt.AllowHeaders, ", "))

		// Access-Control-Max-Age: How long browsers can cache preflight response (in seconds)
		w.Header().Set("Access-Control-Max-Age", maxAge)

		// Handle preflight requests
		// Preflight requests are automatically sent by browsers for "non-simple" requests
		// (e.g., requests with custom headers, non-GET/POST methods, etc.)
		if r.Method == "OPTIONS" {
			// Return 200 OK for preflight requests without processing further
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to the next handler for actual requests (not preflight)
		next.ServeHTTP(w, r)
	})
}
