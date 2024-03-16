package leader

import (
	"net/http"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"

	"gitlab.com/arnaud-web/neli-webservices/test/assert"
)

func TestCreate(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.LeaderRole)
	p := stub.UserParam(u.Email, u.Firstname, u.Lastname)

	expect.LeaderInsert(u, mock)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.User(t, rr, u)
}

func TestCreate_LastnameMissing(t *testing.T) {
	u := stub.User(models.LeaderRole)
	p := stub.UserParam(u.Email, u.Firstname, "")

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_FirstnameMissing(t *testing.T) {
	u := stub.User(models.LeaderRole)
	p := stub.UserParam(u.Email, "", u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_EmailMissing(t *testing.T) {
	u := stub.User(models.LeaderRole)
	p := stub.UserParam("", u.Firstname, u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_EmailInvalid(t *testing.T) {
	u := stub.User(models.LeaderRole)
	p := stub.UserParam(u.Firstname, u.Firstname, u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestList(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByRole(models.LeaderRole, mock)

	rr, r := stub.Http()

	List(rr, r)

	assert.User(t, rr, u)
}
