package credentials

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/auth/bearer"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/mail"
	"gitlab.com/arnaud-web/neli-webservices/random"
	"golang.org/x/crypto/bcrypt"
)

type loginInfo struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	TTL      int64  `json:"ttl"`
}

type setPasswordInfo struct {
	Password string `json:"newPassword"`
	Token    string `json:"setPasswordToken"`
}

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Login checks credentials then creates access and refresh tokens.
// Parameters required are 'login' (string) and 'password' (string).
// Parameter 'ttl' is optional (int).
// If the credentials check fails a 401 error is returned.
// If ttl is not a number or is greater than token life a 400 error is returned
// The success return is a 'authResult' struct
func Login(w http.ResponseWriter, r *http.Request) {
	var l loginInfo

	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	u, err := validate(l)

	if err != nil {
		api.SendError(w, http.StatusUnauthorized, messages.InvalidLoginOrPassword)
		return
	}

	b, err := bearer.Create(u, l.TTL)

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusUnauthorized, messages.TechnicalError)
		return
	}

	api.Send(w, http.StatusOK, b)
}

// Reset generates a reinitialization token and send email to user.
// The token is a random string.
// The email function is called as a go routine.
func Reset(w http.ResponseWriter, r *http.Request) {
	var u models.User
	var l loginInfo

	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if l.Login == "" {
		api.Send(w, http.StatusNoContent, "")
		return
	}

	if err := u.FindByLogin(l.Login); err != nil {
		api.Send(w, http.StatusNoContent, "")
		return
	}

	t := random.String(*config.TokenSize)

	if err := u.UpdatePasswordToken(t); err != nil {
		log.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go mail.ResetPassword(t, u.Email)

	api.Send(w, http.StatusNoContent, "")
}

// Set tries to retrieve user with  "setPasswordToken".
// If he exists, his password is updated.
func Set(w http.ResponseWriter, r *http.Request) {
	var u models.User
	var s setPasswordInfo

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if s.Password == "" || s.Token == "" {
		logger.Println("Password or Token empty")
		api.Send(w, http.StatusNoContent, "")
		return
	}

	if err := u.FindByPasswordToken(s.Token); err != nil {
		logger.Println("Token not found")
		api.Send(w, http.StatusNoContent, "")
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = u.UpdatePassword(string(h)); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Println("Password done")
	api.Send(w, http.StatusNoContent, "")
}

// validate checks if login is in database then compares password with bcrypt hash.
func validate(l loginInfo) (models.User, error) {
	var u models.User

	if l.Login == "" || l.Password == "" {
		return u, errors.New(messages.InvalidParameters)
	}

	if err := u.FindByLogin(l.Login); err != nil {
		return u, err
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Password))
	return u, err
}
