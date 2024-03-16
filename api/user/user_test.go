package user

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"

	"gitlab.com/arnaud-web/neli-webservices/test/assert"

	"github.com/icrowley/fake"
)

func TestGet(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)

	rr, r := stub.HttpWithContext(u.ID)

	Get(rr, r)

	assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", u.ID))
	assert.Contains(t, rr.Body.String(), u.Email)
	assert.Contains(t, rr.Body.String(), u.Firstname)
	assert.Contains(t, rr.Body.String(), u.Lastname)
}
func TestEdit_Admin(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.LeaderRole, mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.AdminOwner(uid, models.AdminRole, mock)
	expect.UserUpdate(1, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_Member(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderOwner(u.ID, uid, mock)
	expect.UserUpdate(1, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_LastnameMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), "")

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderOwner(u.ID, uid, mock)
	expect.UserUpdate(1, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_FirstnameMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	p := stub.UserParam(fake.EmailAddress(), "", fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderOwner(u.ID, uid, mock)
	expect.UserUpdate(1, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_EmailMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)
	p := stub.UserParam("", u.Firstname, u.Lastname)

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestEdit_FirstNameAndLastNameMissing(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)
	p := stub.UserParam(u.Email, "", "")

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestEdit_BadEmail(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)
	p := stub.UserParam(u.Firstname, u.Firstname, u.Lastname)

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidUserInfo)
}

func TestEdit_UserNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.NotFind(mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_AdminBadRole(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.LeaderRole, mock)
	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.AdminOwner(u.ID, models.LeaderRole, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.ActionForbidden)
}

func TestEdit_MemberBadRole(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderNotOwner(u.ID, uid, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.ActionForbidden)
}

func TestEdit_MemberNotModified(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	p := stub.UserParam(fake.EmailAddress(), fake.FirstName(), fake.LastName())

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderOwner(u.ID, uid, mock)
	expect.UserUpdate(0, mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNotModified, "")
}

func TestDelete_Admin(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.LeaderRole, mock)

	uid := int64(rand.Uint64())
	rr, r := stub.HttpWithContext(uid)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.AdminOwner(uid, models.AdminRole, mock)
	expect.ContentUpdate(mock)
	expect.UserDelete(u, mock)

	Delete(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestDelete_Member(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	uid := int64(rand.Uint64())
	rr, r := stub.HttpWithContext(uid)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderOwner(u.ID, uid, mock)
	expect.UserDelete(u, mock)

	Delete(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestDeleted_UserNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.NotFind(mock)

	uid := int64(rand.Uint64())
	rr, r := stub.HttpWithContext(uid)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	Delete(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestDelete_AdminBadRole(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.LeaderRole, mock)

	uid := int64(rand.Uint64())
	rr, r := stub.HttpWithContext(uid)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.AdminOwner(uid, models.LeaderRole, mock)

	Delete(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.ActionForbidden)
}

func TestDelete_MemberBadRole(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindByIdWithRole(models.MemberRole, mock)

	uid := int64(rand.Uint64())
	rr, r := stub.HttpWithContext(uid)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", u.ID))

	expect.LeaderNotOwner(u.ID, uid, mock)

	Delete(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.ActionForbidden)
}

func TestPassword(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)
	expect.UserUpdate(1, mock)

	p := stub.PasswordParam(fake.SimplePassword(), u.Password)

	rr, r := stub.PostHttpWithContext(u.ID, p)

	Password(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestPassword_NewPasswordEmpty(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)
	expect.UserUpdate(1, mock)

	p := stub.PasswordParam("", u.Password)

	rr, r := stub.PostHttpWithContext(u.ID, p)

	Password(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestPassword_OldPasswordEmpty(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)
	expect.UserUpdate(1, mock)

	p := stub.PasswordParam(u.Password, "")

	rr, r := stub.PostHttpWithContext(u.ID, p)

	Password(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestPassword_UserNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)
	expect.UserUpdate(1, mock)

	p := stub.PasswordParam(fake.SimplePassword(), u.Password)

	rr, r := stub.PostHttpWithContext(int64(rand.Uint64()), p)

	Password(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.UserNotFound)
}

func TestPassword_BadPassword(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	u := expect.FindById(mock)
	expect.UserUpdate(1, mock)

	p := stub.PasswordParam(fake.SimplePassword(), fake.SimplePassword())

	rr, r := stub.PostHttpWithContext(u.ID, p)

	Password(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}
