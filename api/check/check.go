package check

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/jwtauth"
	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Authenticator validate JWT token
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			api.SendError(w, http.StatusUnauthorized, messages.InvalidToken)
			return
		}

		if token == nil || !token.Valid {
			api.SendError(w, http.StatusUnauthorized, messages.InvalidToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// SuperAdmin checks super admin role
func SuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !granted(r, models.SuperAdminRole) {
			api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SuperAdmin checks Admin role
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !granted(r, models.AdminRole) {
			api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Leader checks Leader role
func Leader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !granted(r, models.LeaderRole) {
			api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Hub checks Hub role
func Hub(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !granted(r, models.HubRole) {
			api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func granted(r *http.Request, grant string) bool {
	_, claims, _ := jwtauth.FromContext(r.Context())
	role := claims["role"]

	return strings.Contains(grant, role.(string))
}
