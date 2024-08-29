package chi

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// EmptySuccessJSONResponse is a function that sends a successful JSON response with no data.
// It sets the request ID header and the status of the response, and then sends a JSON response with a status of "success".
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// r: The http.Request that we are responding to.
// status: The HTTP status code to set in the response.
func EmptySuccessJSONResponse(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "success",
	})
}

// SuccessJSONResponseWithData is a function that sends a successful JSON response with some data.
// It sets the request ID header and the status of the response, and then sends a JSON response with a status of "success" and the provided data.
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// r: The http.Request that we are responding to.
// status: The HTTP status code to set in the response.
// data: The data to include in the response.
func SuccessJSONResponseWithData(w http.ResponseWriter, r *http.Request, status int, data map[string]any) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "success",
		"data":   data,
	})
}

// FailedJSONResponse is a function that sends a failed JSON response.
// It sets the request ID header and the status of the response, and then sends a JSON response with a status of "failed" and the provided data.
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// r: The http.Request that we are responding to.
// status: The HTTP status code to set in the response.
// data: The data to include in the response.
func FailedJSONResponse(w http.ResponseWriter, r *http.Request, status int, data []map[string]any) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status": "failed",
		"data":   data,
	})
}

// ErrorJSONResponse is a function that sends an error JSON response.
// It sets the request ID header and the status of the response, and then sends a JSON response with a status of "error", a code, and a message.
//
// Parameters:
// w: The http.ResponseWriter to write the response to.
// r: The http.Request that we are responding to.
// status: The HTTP status code to set in the response.
// message: The error message to include in the response.
func ErrorJSONResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set(middleware.RequestIDHeader, middleware.GetReqID(r.Context()))

	render.Status(r, status)
	render.JSON(w, r, map[string]any{
		"status":  "error",
		"message": message,
	})
}
