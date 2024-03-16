package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth"
)

// JSONError is the default struct sent by api
type JSONError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func header(w http.ResponseWriter, c int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(c)
}

// Send a success response
func Send(w http.ResponseWriter, c int, i interface{}) {
	header(w, c)
	json.NewEncoder(w).Encode(i)
}

// SendError an error response
func SendError(w http.ResponseWriter, c int, msg string) {
	header(w, c)
	json.NewEncoder(w).Encode(JSONError{c, msg})
}

// UserIdFromContext extracts id from JWT token.
// If checks all numeric types to be sure to match with the numeric format.
// But the format returned by JWT token library should be float64.
func UserIdFromContext(r *http.Request) int64 {
	// Error can be ignored because it's the result of error casting.
	// See FromContext implementation for more details.
	_, c, _ := jwtauth.FromContext(r.Context())

	u := c["user"]

	switch n := u.(type) {
	case int:
		return int64(n)
	case int8:
		return int64(n)
	case int16:
		return int64(n)
	case int32:
		return int64(n)
	case int64:
		return int64(n)
	case float32:
		return int64(n)
	case float64:
		return int64(n)
	}

	return 0
}
