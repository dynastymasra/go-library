package middleware

import (
	"context"
	"net/http"

	"github.com/dynastymasra/go-library/web"
)

// Service represents a service with a name and version.
// It is used to add service-specific headers to HTTP responses.
type Service struct {
	Name    string
	Version string
}

// AddServiceHeader adds service-specific headers to the HTTP response.
// It sets the service name and version in the request context and then calls the next handler.
//
// Parameters:
// - next: The next http.Handler to call after setting the service headers.
//
// Returns:
// - An http.Handler that sets the service headers and then calls the next handler.
func (s Service) AddServiceHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), web.ServiceName, s.Name)
		ctx = context.WithValue(ctx, web.ServiceVersion, s.Version)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
