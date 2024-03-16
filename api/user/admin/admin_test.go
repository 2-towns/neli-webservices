package admin

import (
	"net/http"
	"testing"

	"github.com/icrowley/fake"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"

	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func TestCreate(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := stub.User(models.AdminRole)
	p := stub.UserParam(u.Email, u.Firstname, u.Lastname)

	expect.AdminInsert(u, mock)

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.User(t, rr, u)
}

func TestCreate_LastnameMissing(t *testing.T) {
	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), "")

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_FirstnameMissing(t *testing.T) {
	p := stub.UserParam(fake.EmailAddress(), "", fake.LastName())

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_EmailMissing(t *testing.T) {
	p := stub.UserParam("", fake.FirstName(), fake.LastName())

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestCreate_EmailInvalid(t *testing.T) {
	p := stub.UserParam(fake.FirstName(), fake.FirstName(), fake.LastName())

	rr, r := stub.PostHttp(p)

	Create(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestList(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByRole(models.AdminRole, mock)

	rr, r := stub.Http()

	List(rr, r)

	assert.User(t, rr, u)
}
