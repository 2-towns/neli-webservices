package content

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/mail"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

type contentInfo struct {
	MaxDuration int `json:"maxDuration"`
}

type contentReady struct {
	RandomString string `json:"randomString"`
	Duration     int    `json:"duration"`
}

type videoContentSharingInfo struct {
	ExpirationDate string `json:"expirationDate"`
	Message        string `json:"message"`
}

type singleShareInfo struct {
	ExpirationDate time.Time
	Message        string
	UserID         int64
}

type videoContentSharingInfoMultiple struct {
	Ids            []int64 `json:"ids"`
	ExpirationDate string  `json:"expirationDate"`
	Message        string  `json:"message"`
}

type multipleShareError struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Ids     []int64 `json:"ids"`
}

// Get a content related to the leader identified
func Get(w http.ResponseWriter, r *http.Request) {
	lid := api.UserIdFromContext(r)
	l := models.Leader{User: &models.User{ID: lid}}
	c, err := l.Content()

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, c)
}

// Shares list all share by the leader identified
func Shares(w http.ResponseWriter, r *http.Request) {
	// Ignore error because uid is necessarily an int due to match param url
	cid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)

	s := new(models.Share)
	c, err := s.List(cid)

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, c)
}

// All returns all shares by all leader
func All(w http.ResponseWriter, r *http.Request) {
	aid := api.UserIdFromContext(r)
	a := models.Admin{User: &models.User{ID: aid}}
	c, err := a.Content()

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, c)
}

// New create a new content.
// It loads the max durations from settings.
func New(w http.ResponseWriter, r *http.Request) {
	uid := api.UserIdFromContext(r)
	cnt := models.Content{LeaderID: uid}

	if err := cnt.New(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s := models.Settings{}

	if err := s.Load(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	cnt.LeaderID = 0
	cnt.MaxDuration = s.MaxDuration

	api.Send(w, http.StatusOK, struct {
		VideoContentID int64 `json:"videoContentId"`
		MaxDuration    int   `json:"maxDuration"`
	}{
		cnt.ID,
		s.MaxDuration,
	})
}

// MaxDuration set max duration.
func MaxDuration(w http.ResponseWriter, r *http.Request) {
	c := contentInfo{}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	s := models.Settings{MaxDuration: c.MaxDuration}

	if err := s.Update(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

// GetMaxDuration return the max duration
func GetMaxDuration(w http.ResponseWriter, r *http.Request) {
	s := models.Settings{}

	if err := s.Load(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusOK, s)
}

// Edit a share by a leader.
func Edit(w http.ResponseWriter, r *http.Request) {
	uid := api.UserIdFromContext(r)
	// Ignore error because uid is necessarily an int due to match param url
	cid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)
	c := models.Content{ID: cid, LeaderID: uid}

	if err := c.Find(); err != nil {
		api.SendError(w, http.StatusNotFound, messages.InvalidVideoContentId)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if c.Name == "" || c.Description == "" {
		api.SendError(w, http.StatusBadRequest, errors.New(messages.InvalidParameters).Error())
		return
	}

	if err := c.Update(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

// Share records new share. Parameters have to be valid:
// * expirationDate has to valid est formatted with RFC3339
// *  message hast to be not empty.
// It generates a new share url encoded in base 64 and sent it to user.
// If the share exist already so a new mail is sent again.
func Share(w http.ResponseWriter, r *http.Request) {
	v := videoContentSharingInfo{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	t, err := time.Parse(time.RFC3339, v.ExpirationDate)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	// Ignore error because uid is necessarily an int due to match param url
	uid, _ := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)

	// Ignore error because vid is necessarily an int due to match param url
	vid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)

	lid := api.UserIdFromContext(r)
	c := models.Content{ID: vid, LeaderID: lid}

	if err := c.Find(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusNotFound, messages.InvalidVideoContentIdOrUserId)
		return
	}

	m := models.Member{User: &models.User{ID: uid}}

	if b := m.InTribe(lid); !b {
		logger.Println(err)
		api.SendError(w, http.StatusNotFound, messages.InvalidVideoContentIdOrUserId)
		return
	}

	s := models.Share{
		UserID:         uid,
		ContentID:      vid,
		Message:        v.Message,
		ExpirationDate: models.JSONTime(t),
	}

	err = s.New(&c)
	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if s.IsReady() {
		go sendShareURL(uid, &s, &c)
	}

	api.Send(w, http.StatusNoContent, nil)
}

func singleShare(w http.ResponseWriter, r *http.Request, ssi singleShareInfo, ch chan int64) {
	// Ignore error because vid is necessarily an int due to match param url
	vid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)

	lid := api.UserIdFromContext(r)
	c := models.Content{ID: vid, LeaderID: lid}

	if err := c.Find(); err != nil {
		logger.Println(err)
		ch <- c.ID
		return
	}

	m := models.Member{User: &models.User{ID: ssi.UserID}}

	if b := m.InTribe(lid); !b {
		ch <- ssi.UserID
		return
	}

	s := models.Share{
		UserID:         ssi.UserID,
		ContentID:      vid,
		Message:        ssi.Message,
		ExpirationDate: models.JSONTime(ssi.ExpirationDate),
	}

	err := s.New(&c)
	if err != nil {
		logger.Println(err)
		ch <- ssi.UserID
		return
	}

	if s.IsReady() {
		go sendShareURL(ssi.UserID, &s, &c)
	}

	ch <- 0
}

func MultipleShares(w http.ResponseWriter, r *http.Request) {
	v := videoContentSharingInfoMultiple{}

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	t, err := time.Parse(time.RFC3339, v.ExpirationDate)
	if err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	c := make(chan int64, len(v.Ids))

	for _, i := range v.Ids {
		s := singleShareInfo{
			UserID:         i,
			ExpirationDate: t,
			Message:        v.Message,
		}
		go singleShare(w, r, s, c)
	}

	var ids []int64
	for _ = range v.Ids {
		uid := <-c
		if uid > 0 {
			ids = append(ids, uid)
		}
	}

	if len(ids) > 0 {
		mse := multipleShareError{
			Code:    404,
			Message: messages.InvalidUserIdOrContentId,
			Ids:     ids,
		}
		api.Send(w, http.StatusNotFound, mse)
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

func sendShareURL(uid int64, share *models.Share, c *models.Content) {
	u := new(models.User)

	if err := u.Find(uid); err != nil {
		logger.Println(err)
		return
	}

	l := new(models.User)

	if err := l.Find(c.LeaderID); err != nil {
		logger.Println(err)
		return
	}

	mail.Share(u, share, c, l)
}

// Delete a content.
// The content is not really deleted by linked to zombie leader.
func Delete(w http.ResponseWriter, r *http.Request) {
	lid := api.UserIdFromContext(r)
	l := models.Leader{User: &models.User{ID: lid}}

	// Ignore error because uid is necessarily an int due to match param url
	cid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)

	if b := l.IsOwner(cid); !b {
		api.SendError(w, http.StatusForbidden, messages.ActionForbidden)
		return
	}

	if err := l.DeleteContent(cid); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

func Ready(w http.ResponseWriter, r *http.Request) {
	cr := contentReady{}

	if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	if cr.RandomString == "" {
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	cid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)

	// Hack for hub
	// We return status no content if hub send video content id equals to 0
	if cid == 0 {
		api.Send(w, http.StatusNoContent, nil)
		return
	}

	c := models.Content{
		ID:       cid,
		Path:     cr.RandomString,
		Duration: cr.Duration,
	}

	s := models.Settings{}

	if err := s.Load(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if c.Duration > s.MaxDuration {
		logger.Printf("Duration %d truncated to %d", c.Duration, s.MaxDuration)
		c.Duration = s.MaxDuration
	}

	shares, err := c.Shares()

	if err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, messages.TechnicalError)
		return
	}

	for _, share := range shares {

		err := share.UpdateURL(&c)

		if err != nil {
			logger.Println(err)
			api.SendError(w, http.StatusInternalServerError, messages.TechnicalError)
			return
		}

		go sendShareURL(share.UserID, &share, &c)

	}

	if err := c.UpdatePath(); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusInternalServerError, messages.TechnicalError)
		return
	}

	api.Send(w, http.StatusNoContent, nil)
}

func Cleaning(w http.ResponseWriter, r *http.Request) {
	s := struct {
		Date string `json: "date",omitempty`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}

	defer r.Body.Close()

	var t string

	if s.Date == "" {
		t = time.Now().Format("2006-01-02")
	} else {
		t = s.Date
	}

	if err := models.Clean(t); err != nil {
		log.Println(err)
		api.SendError(w, http.StatusInternalServerError, messages.TechnicalError)
		return
	}

	api.Send(w, http.StatusNoContent, nil)

	return
}
