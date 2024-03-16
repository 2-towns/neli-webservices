package user

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/db/models"

	"github.com/go-chi/chi"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Resource to be stick with specification
type Resource struct {
	*UserInfo `json:"resource,omitempty"`
}

type changePasswordInfo struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type UserInfo struct {
	UserID int64  `json:"userId,omitempty"`
	Login  string `json:"login,omitempty"`
	Role   string `json:"profile,omitempty"`
}

// Get return user information identified by user id in token
func Get(w http.ResponseWriter, r *http.Request) {
	u := new(models.User)
	uid := api.UserIdFromContext(r)

	if err := u.Find(uid); err != nil {
		api.SendError(w, http.StatusNotFound, messages.UserNotFound)
		return
	}

	api.Send(w, http.StatusOK, u.ToJSON())
}

// Edit user informations.
// Admin can edit leader informations only.
// Leader can edit member his member informations only.
func Edit(w http.ResponseWriter, r *http.Request) {
	// Ignore error because uid is necessarily an int due to match param url
	uid, _ := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	u := models.User{ID: uid}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	j, err := u.Jsonify()

	if err != nil {
		api.Send(w, http.StatusNoContent, nil)
		return
	}

	if err := j.Validate(); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidUserInfo)
		return
	}

	oid := api.UserIdFromContext(r)

	if j.EmailExists(oid) {
		api.SendError(w, http.StatusBadRequest, messages.EmailExists)
		return
	}

	if b := j.IsOwned(oid); !b {
		api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
		return
	}

	sqlr, err := j.Save()

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	if a, _ := sqlr.RowsAffected(); a <= 0 {
		api.Send(w, http.StatusNotModified, nil)
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

// Delete an user
func Delete(w http.ResponseWriter, r *http.Request) {
	// Ignore error because uid is necessarily an int due to match param url
	uid, _ := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
	u := models.User{ID: uid}

	j, err := u.Jsonify()

	if err != nil {
		api.Send(w, http.StatusNoContent, nil)
		return
	}

	oid := api.UserIdFromContext(r)

	if b := j.IsOwned(oid); !b {
		api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
		return
	}

	if err := j.Delete(); err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

// Password change user password.
// Two parameters are required:
// * newPassword has to be not empty
// * oldPassword has to match with previous password
func Password(w http.ResponseWriter, r *http.Request) {
	c := changePasswordInfo{}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if c.NewPassword == "" || c.OldPassword == "" {
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	u := new(models.User)
	uid := api.UserIdFromContext(r)

	if err := u.Find(uid); err != nil {
		api.SendError(w, http.StatusNotFound, messages.UserNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.OldPassword)); err != nil {
		api.SendError(w, http.StatusNoContent, "")
		return
	}

	h, err := bcrypt.GenerateFromPassword([]byte(c.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := u.UpdatePassword(string(h)); err != nil {
		logger.Println(err)
		api.Send(w, http.StatusBadRequest, "")
		return
	}

	api.Send(w, http.StatusNoContent, "")
}
