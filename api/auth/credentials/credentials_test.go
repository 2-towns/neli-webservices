package credentials

import (
	"net/http"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"

	"github.com/icrowley/fake"
	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func loginParam(l, p string, ttl int64) map[string]interface{} {
	return map[string]interface{}{
		"login":    l,
		"password": p,
		"ttl":      ttl,
	}
}

func passwordParam(p, pt string) map[string]interface{} {
	return map[string]interface{}{
		"newPassword":      p,
		"setPasswordToken": pt,
	}
}

func TestLogin(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := loginParam(fake.UserName(), fake.SimplePassword(), 0)

	expect.FindByLogin(p["login"].(string), p["password"].(string), mock)

	rr, r := stub.PostHttp(p)

	Login(rr, r)

	assert.Response(t, rr, http.StatusOK, "")
}

func TestLogin_LoginMissing(t *testing.T) {
	p := loginParam("", fake.SimplePassword(), 0)

	rr, r := stub.PostHttp(p)

	Login(rr, r)

	assert.Response(t, rr, http.StatusUnauthorized, messages.InvalidParameters)
}

func TestLogin_PasswordMissing(t *testing.T) {
	p := loginParam(fake.UserName(), "", 0)

	rr, r := stub.PostHttp(p)

	Login(rr, r)

	assert.Response(t, rr, http.StatusUnauthorized, messages.InvalidParameters)
}

func TestLogin_BadPassword(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := loginParam(fake.UserName(), fake.SimplePassword(), 0)

	expect.FindByLogin(p["login"].(string), p["password"].(string), mock)

	p["password"] = fake.SimplePassword()

	rr, r := stub.PostHttp(p)

	Login(rr, r)

	assert.Response(t, rr, http.StatusUnauthorized, messages.InvalidParameters)
}

func TestLogin_BadTtl(t *testing.T) {
	p := map[string]interface{}{
		"login":    fake.UserName(),
		"password": fake.SimplePassword(),
		"ttl":      fake.Sentence(),
	}

	rr, r := stub.PostHttp(p)

	Login(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestReset(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := loginParam(fake.UserName(), "", 0)

	u := expect.FindByLogin(p["login"].(string), p["password"].(string), mock)
	expect.UserUpdate(u.ID, mock)

	rr, r := stub.PostHttp(p)

	Reset(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestReset_NoLogin(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := loginParam("", "", 0)

	expect.FindByLogin(p["login"].(string), p["password"].(string), mock)

	rr, r := stub.PostHttp(p)

	Reset(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestReset_LoginNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := loginParam(fake.UserName(), "", 0)

	expect.NotFind(mock)

	rr, r := stub.PostHttp(p)

	Reset(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestSet(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := passwordParam(fake.SimplePassword(), fake.Sentence())

	u := expect.FindByPasswordToken(p["setPasswordToken"].(string), mock)
	expect.UserUpdate(u.ID, mock)

	rr, r := stub.PostHttp(p)

	Set(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestSetWithoutNewPassword(t *testing.T) {
	p := passwordParam("", fake.Sentence())

	rr, r := stub.PostHttp(p)

	Set(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestSetWithoutPasswordToken(t *testing.T) {
	p := passwordParam(fake.SimplePassword(), "")

	rr, r := stub.PostHttp(p)

	Set(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestSetBadPasswordToken(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := passwordParam(fake.SimplePassword(), fake.Sentence())

	expect.FindByPasswordToken(p["setPasswordToken"].(string), mock)

	p["setPasswordToken"] = fake.Sentence()

	rr, r := stub.PostHttp(p)

	Set(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}
