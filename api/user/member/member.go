package member

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/api/user"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Create create an user and attach it to leader tribe
// Email has to be valid, firstname and lastname not empty.
func Create(w http.ResponseWriter, r *http.Request) {
	m := models.Member{User: &models.User{}}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if err := m.Validate(); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidUserInfo)
		return
	}

	lid := api.UserIdFromContext(r)

	if err := m.ExistsInTribe(lid); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.EmailExists)
		return
	}

	m.Role = models.MemberRole

	if err := m.New(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	t := models.Tribe{LeaderID: lid, UserID: m.ID}

	if _, err := t.New(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u := models.User{ID: m.ID, Role: m.Role}

	res := user.Resource{UserInfo: &user.UserInfo{}}
	res.UserID = u.ID
	res.Role = models.MemberRole

	api.Send(w, http.StatusOK, res)
}

// List return the list of all members
func List(w http.ResponseWriter, r *http.Request) {
	l := []models.Member{}

	uid := api.UserIdFromContext(r)
	l, err := models.ListMembers(uid)

	if err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, l)
}
