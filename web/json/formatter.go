package json

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/dynastymasra/go-library/web"
)

// SuccessResponse sends a successful JSON response without any additional data.
// It sets the request ID header, service name, and service version in the response headers.
// Then it sets the status of the response and sends a JSON response with a status of "success".
//
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request that we are responding to.
// - status: The HTTP status code to set in the response.
func SuccessResponse(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))
	w.Header().Set(web.XServiceName, fmt.Sprintf("%v", r.Context().Value(web.ServiceName)))
	w.Header().Set(web.XServiceVersion, fmt.Sprintf("%v", r.Context().Value(web.ServiceVersion)))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "success",
	})
}

// DataResponse sends a successful JSON response with data.
// It sets the request ID header, service name, and service version in the response headers.
// Then it sets the status of the response and sends a JSON response with a status of "success" and the provided data.
//
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request that we are responding to.
// - status: The HTTP status code to set in the response.
// - data: The data to include in the response.
func DataResponse(w http.ResponseWriter, r *http.Request, status int, data map[string]any) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))
	w.Header().Set(web.XServiceName, fmt.Sprintf("%v", r.Context().Value(web.ServiceName)))
	w.Header().Set(web.XServiceVersion, fmt.Sprintf("%v", r.Context().Value(web.ServiceVersion)))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "success",
		"data":   data,
	})
}

// FailedResponse sends a failed JSON response with data.
// It sets the request ID header, service name, and service version in the response headers.
// Then it sets the status of the response and sends a JSON response with a status of "failed" and the provided data.
//
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request that we are responding to.
// - status: The HTTP status code to set in the response.
// - data: The data to include in the response.
func FailedResponse(w http.ResponseWriter, r *http.Request, status int, data []map[string]any) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))
	w.Header().Set(web.XServiceName, fmt.Sprintf("%v", r.Context().Value(web.ServiceName)))
	w.Header().Set(web.XServiceVersion, fmt.Sprintf("%v", r.Context().Value(web.ServiceVersion)))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "failed",
		"data":   data,
	})
}

// ErrorResponse sends an error JSON response with a message.
// It sets the request ID header, service name, and service version in the response headers.
// Then it sets the status of the response and sends a JSON response with a status of "error" and the provided message.
//
// Parameters:
// - w: The http.ResponseWriter to write the response to.
// - r: The http.Request that we are responding to.
// - status: The HTTP status code to set in the response.
// - message: The error message to include in the response.
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))
	w.Header().Set(web.XServiceName, fmt.Sprintf("%v", r.Context().Value(web.ServiceName)))
	w.Header().Set(web.XServiceVersion, fmt.Sprintf("%v", r.Context().Value(web.ServiceVersion)))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status":  "error",
		"message": message,
	})
}
