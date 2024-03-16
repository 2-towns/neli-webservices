package leader

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/api/user"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/mail"
	"gitlab.com/arnaud-web/neli-webservices/random"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Create creates a new user with role leader.
// Only an administrator can do that.
// Email has to be valid, firstname and lastname not empty.
func Create(w http.ResponseWriter, r *http.Request) {
	u := models.Leader{User: new(models.User)}

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if err := u.Validate(); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidUserInfo)
		return
	}

	if err := u.Exists(); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.EmailExists)
		return
	}

	u.Role = models.LeaderRole
	u.PasswordResetToken = random.String(*config.TokenSize)

	if err := u.New(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go mail.Leader(u.PasswordResetToken, u.Email)

	res := user.Resource{UserInfo: &user.UserInfo{}}
	res.UserID = u.ID
	res.Role = models.LeaderRole

	api.Send(w, http.StatusOK, res)
}

// List return the list of all leader, zombie included
func List(w http.ResponseWriter, r *http.Request) {
	l, err := models.ListUsers(models.LeaderRole)

	if err != nil {
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, l)
}
