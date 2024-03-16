package member

import (
	"net/http"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"

	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func TestCreate(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.MemberRole)

	expect.MemberInsert(u, mock)
	expect.TribeInsert(u, mock)

	p := stub.UserParam(u.Email, u.Firstname, u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.User(t, rr, u)
}

func TestCreate_LastnameMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.MemberRole)
	u.Lastname = ""

	expect.MemberInsert(u, mock)
	expect.TribeInsert(u, mock)

	p := stub.UserParam(u.Email, u.Firstname, "")

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.User(t, rr, u)
}

func TestCreate_FirstnameMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.MemberRole)
	u.Firstname = ""

	expect.MemberInsert(u, mock)
	expect.TribeInsert(u, mock)

	p := stub.UserParam(u.Email, "", u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.User(t, rr, u)
}

func TestCreate_EmailMissing(t *testing.T) {
	u := stub.User(models.MemberRole)
	p := stub.UserParam("", u.Firstname, u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_FirstnameAndLastnameMissing(t *testing.T) {
	u := stub.User(models.MemberRole)
	p := stub.UserParam(u.Email, "", "")

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_EmailInvalid(t *testing.T) {
	u := stub.User(models.MemberRole)
	p := stub.UserParam(u.Firstname, u.Firstname, u.Lastname)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestList(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.LeaderRole)

	m := expect.FindMember(u.ID, mock)

	rr, r := stub.HttpWithContext(u.ID)

	List(rr, r)

	assert.User(t, rr, m)
}
