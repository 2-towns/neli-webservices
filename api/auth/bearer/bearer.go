package bearer

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth"
	"gitlab.com/arnaud-web/neli-webservices/api"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Algorithm used to generated JWT token
const TokenAlgorithm = "HS256"

type result struct {
	AccessToken         string `json:"accessToken"`
	AccessTokenExpires  string `json:"accessTokenExpires"`
	RefreshToken        string `json:"refreshToken"`
	RefreshTokenExpires string `json:"refreshTokenExpires"`
}

// Create an 'authResult' by creating access and refresh token.
// Access token expires date is calculated by adding token life configuration time to 'time.Now()'.
// Refresh token expires date is calculated by subtracting 'ttl' to access token expires date.
// Expires date are string formatted by using 'time.RFC3339' format.
func Create(u models.User, ttl int64) (result, error) {
	now := time.Now().Unix()
	ts := now + *config.TokenLife
	atexp := time.Unix(ts, 0).Format(time.RFC3339)
	at, err := tokenize(u, *config.TokenSecret, ts)

	if err != nil {
		return result{}, err
	}

	// If ttl exist so it's the time life of refresh token
	if ttl > 0 {
		ts = now + ttl
	}

	rtexp := time.Unix(ts, 0).Format(time.RFC3339)
	rt, err := tokenize(u, *config.RefreshSecret, ts)

	if err != nil {
		return result{}, err
	}

	return result{
		at,
		atexp,
		rt,
		rtexp,
	}, nil
}

// Refresh generates new access and refresh token.
// User information are extracted from refresh token.
func Refresh(w http.ResponseWriter, r *http.Request) {
	u := models.User{}
	uid := api.UserIdFromContext(r)

	if err := u.Find(uid); err != nil {
		api.SendError(w, http.StatusNotFound, messages.UserNotFound)
		return
	}

	ar, err := Create(u, 0)

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, messages.TechnicalError)
		return
	}

	api.Send(w, http.StatusOK, ar)
}

// Create a JWT token with user data inside.
// The token is signed and returned as a string value.
func tokenize(u models.User, secret string, exp int64) (string, error) {
	tokenAuth := jwtauth.New(TokenAlgorithm, []byte(secret), nil)
	c := jwtauth.Claims{"user": u.ID, "role": u.Role}.SetExpiry(time.Unix(exp, 0))
	_, tokenString, err := tokenAuth.Encode(c)
	return tokenString, err
}
