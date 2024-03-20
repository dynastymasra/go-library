package middleware

import "net/http"

// AddSecurityHeaders is a middleware function that adds security headers to the HTTP response.
// It takes a http.Handler as an argument which represents the next handler to be executed in the middleware chain.
// The function adds several security headers to the response, including:
// - Strict-Transport-Security: This header is used to enforce secure (HTTP over SSL/TLS) connections to the server.
// - X-Content-Type-Options: This header is used to protect against MIME type confusion attacks.
// - X-Frame-Options: This header is used to indicate whether a browser should be allowed to render a page in a <frame>, <iframe>, <embed> or <object>.
// - Content-Security-Policy: This header is used to prevent a wide range of attacks, including Cross-site scripting and other cross-site injections.
// - X-XSS-Protection: This header is used to configure the XSS Auditor in Chrome, Internet Explorer and Safari (though it's being deprecated).
// - Cache-Control: This header is used to specify directives for caching mechanisms in both requests and responses.
//
// Parameters:
// next: The next http.Handler to be executed in the middleware chain.
//
// Returns:
// A http.Handler that can be used in the middleware chain.
func AddSecurityHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none';")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
