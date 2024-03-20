package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const requestID = "requestId"

// LogRequestWithZerolog is a middleware function that logs HTTP requests and responses.
// It logs the start and end time of the request, the duration, the request details (address, path, method, headers, queries),
// and the response details (status, bytes written, headers). If the response status is 400 or above, it logs a warning.
// Otherwise, it logs an info message.
//
// The function takes the next http.Handler to call in the middleware chain.
// It returns a new http.Handler that wraps the original handler with logging functionality.
func LogRequestWithZerolog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			span := zerolog.Dict().Time("start", now).Time("end", time.Now().UTC()).
				Str("duration", time.Since(now).String())
			request := zerolog.Dict().Str("address", r.RemoteAddr).Str("path", r.URL.Path).
				Str("method", r.Method).Interface("headers", r.Header).
				Interface("queries", r.URL.Query())
			response := zerolog.Dict().Int("status", ww.Status()).
				Int("byte", ww.BytesWritten()).Interface("headers", ww.Header())

			if ww.Status() >= http.StatusBadRequest {
				log.Warn().Str(requestID, middleware.GetReqID(r.Context())).Dict("span", span).
					Dict("request", request).Dict("response", response).Msg("HTTP message logging")
			} else {
				log.Info().Str(requestID, middleware.GetReqID(r.Context())).Dict("span", span).
					Dict("request", request).Dict("response", response).Msg("HTTP message logging")
			}
		}()
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
