package bearer

import (
	"net/http"
	"testing"
	"time"

	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"

	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"

	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

func TestCreate(t *testing.T) {
	u := models.User{}

	auth, _ := Create(u, 0)

	at, _ := time.Parse(
		time.RFC3339,
		auth.AccessTokenExpires)

	et, _ := time.Parse(
		time.RFC3339,
		auth.RefreshTokenExpires)

	b := at.Equal(et)

	if auth.AccessToken == "" || auth.AccessTokenExpires == "" || auth.RefreshToken == "" || auth.RefreshTokenExpires == "" {
		t.Errorf("Access token is not complete")
	}

	if !b {
		t.Errorf("Date %s should be equal to %s", at, et)
	}
}
func TestCreate_WithTTL(t *testing.T) {
	u := models.User{}

	ttl := int64(10)

	auth, _ := Create(u, ttl)

	at, _ := time.Parse(
		time.RFC3339,
		auth.AccessTokenExpires)

	et, _ := time.Parse(
		time.RFC3339,
		auth.RefreshTokenExpires)

	now := at.Add(-time.Second * time.Duration(*config.TokenLife))
	date := now.Add(time.Second * time.Duration(ttl))

	if date.Unix() != et.Unix() {
		t.Errorf("Date %s should be equal to %s", date, et)
	}
}

func TestRefresh(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)

	rr, r := stub.HttpWithContext(u.ID)

	Refresh(rr, r)

	assert.Response(t, rr, http.StatusOK, "")
}

func TestRefresh_UserNotFound(t *testing.T) {
	rr, r := stub.Http()

	Refresh(rr, r)

	assert.ResponseError(t, rr, http.StatusNotFound, messages.UserNotFound)
}
