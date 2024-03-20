package middleware

import (
	"net/http"
	"strings"

	"github.com/dynastymasra/go-library/web/chi"
)

// ContentTypeJSON is a middleware function that checks if the request's content type is application/json.
// It takes a http.Handler as an argument which represents the next handler to be executed in the middleware chain.
// If the content type of the request is not application/json, it responds with a JSON message and a status of http.StatusUnsupportedMediaType.
// If the content type is application/json, it calls the next handler in the middleware chain.
//
// Parameters:
// next: The next http.Handler to be executed in the middleware chain.
//
// Returns:
// A http.Handler that can be used in the middleware chain.
func ContentTypeJSON(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		contentType := strings.Trim(r.Header.Get("Content-Type"), " ")

		if !strings.Contains(strings.ToLower(contentType), "application/json") {
			messages := []map[string]any{
				{
					"message": "Content-Type is empty or not application/json",
				},
			}
			chi.FailedJSONResponse(w, r, http.StatusUnsupportedMediaType, messages)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ContentTypeUTF8 is a middleware function that checks if the request's content type is charset=utf-8.
// It takes a http.Handler as an argument which represents the next handler to be executed in the middleware chain.
// If the content type of the request is not charset=utf-8, it responds with a JSON message and a status of http.StatusUnsupportedMediaType.
// If the content type is charset=utf-8, it calls the next handler in the middleware chain.
//
// Parameters:
// next: The next http.Handler to be executed in the middleware chain.
//
// Returns:
// A http.Handler that can be used in the middleware chain.
func ContentTypeUTF8(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		contentType := strings.Trim(r.Header.Get("Content-Type"), " ")

		if !strings.Contains(strings.ToLower(contentType), "charset=utf-8") {
			messages := []map[string]any{
				{
					"message": "Content-Type is empty or not charset=utf-8",
				},
			}
			chi.FailedJSONResponse(w, r, http.StatusUnsupportedMediaType, messages)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
