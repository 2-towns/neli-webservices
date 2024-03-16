package assert

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

func Response(t *testing.T, rr *httptest.ResponseRecorder, c int, i interface{}) {
	if rr.Code != c {
		t.Errorf("Error code should be '%d' but is '%d'", c, rr.Code)
	}
}

func ResponseError(t *testing.T, rr *httptest.ResponseRecorder, c int, msg string) {
	if rr.Code != c {
		t.Errorf("Error code should be '%d' but is '%d'", c, rr.Code)
	}

	Contains(t, rr.Body.String(), msg)
}

func Equals(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("'%s' not equals to '%s'", a, b)
	}
}

func NotEmpty(t *testing.T, a string) {
	if a == "" {
		t.Errorf("'%s' is empty", a)
	}
}

func NotZero(t *testing.T, i int64) {
	if i == 0 {
		t.Errorf("'%d' is zero", i)
	}
}

func Contains(t *testing.T, a, b string) {
	if !strings.Contains(a, b) {
		t.Errorf("'%s' not found in '%s'", b, a)
	}
}

func User(t *testing.T, rr *httptest.ResponseRecorder, u models.User) {
	Response(t, rr, http.StatusOK, "")

	Contains(t, rr.Body.String(), u.Role)
	Contains(t, rr.Body.String(), fmt.Sprintf("%d", u.ID))

	if u.Role == models.AdminRole {
		Contains(t, rr.Body.String(), u.Email)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Error should be nil")
	}
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error should nil but got : '%s'", err.Error())
	}
}
