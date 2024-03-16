package stub

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/icrowley/fake"
	"github.com/jmoiron/sqlx"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

func Http() (*httptest.ResponseRecorder, *http.Request) {
	return httptest.NewRecorder(), httptest.NewRequest("GET", "/fake/url", nil)
}

func PostHttp(p map[string]interface{}) (*httptest.ResponseRecorder, *http.Request) {
	b, _ := json.Marshal(p)
	r := httptest.NewRequest("GET", "/fake/url", bytes.NewReader(b))
	return httptest.NewRecorder(), r
}

func PostHttpWithContext(uid int64, p map[string]interface{}) (*httptest.ResponseRecorder, *http.Request) {
	b, _ := json.Marshal(p)
	r := httptest.NewRequest("GET", "/fake/url", bytes.NewReader(b))

	c := jwt.MapClaims{"user": int64(uid)}
	t := &jwt.Token{Claims: c}
	ctx := jwtauth.NewContext(r.Context(), t, nil)

	return httptest.NewRecorder(), r.WithContext(ctx)
}

func AddUrlParam(r *http.Request, k, v string) *http.Request {
	var ctrx *chi.Context

	rtx := r.Context().Value(chi.RouteCtxKey)

	if rtx == nil {
		ctrx = chi.NewRouteContext()
		ctrx.URLParams = chi.RouteParams{}
	} else {
		ctrx = rtx.(*chi.Context)
	}

	ctrx.URLParams.Add(k, v)

	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, ctrx)
	return r.WithContext(ctx)
}

func HttpWithContext(uid int64) (*httptest.ResponseRecorder, *http.Request) {
	r := httptest.NewRequest("GET", "/fake/url", nil)

	c := jwt.MapClaims{"user": int64(uid)}
	t := &jwt.Token{Claims: c}

	ctx := jwtauth.NewContext(r.Context(), t, nil)

	return httptest.NewRecorder(), r.WithContext(ctx)
}

func User(r string) models.User {
	return models.User{
		ID:        int64(rand.Uint64()),
		Firstname: fake.FirstName(),
		Lastname:  fake.FirstName(),
		Email:     fake.EmailAddress(),
		Role:      r,
		Password:  fake.SimplePassword(),
	}
}

func Tribe(u *models.User) models.Tribe {
	return models.Tribe{
		LeaderID: int64(rand.Uint64()),
		UserID:   u.ID,
	}
}

func UserParam(e, f, l string) map[string]interface{} {
	return map[string]interface{}{
		"email":     e,
		"firstName": f,
		"lastName":  l,
	}
}

func VideoParam(n, des string) map[string]interface{} {
	return map[string]interface{}{
		"name":        n,
		"description": des,
	}
}

func PasswordParam(n, o string) map[string]interface{} {
	return map[string]interface{}{
		"newPassword": n,
		"oldPassword": o,
	}
}

func ShareParam(msg, exp string) map[string]interface{} {
	return map[string]interface{}{
		"message":        msg,
		"expirationDate": exp,
	}
}

func Content() models.Content {
	return models.Content{
		ID:            int64(rand.Uint64()),
		Name:          fake.UserName(),
		Description:   fake.Sentence(),
		Duration:      rand.Int(),
		LeaderID:      int64(rand.Uint64()),
		CreatedAt:     time.Now(),
		CreationDate:  time.Now().Format("20060102"),
		SharingStatus: true,
	}
}

func Share() models.Share {
	return models.Share{
		UserID:         int64(rand.Uint64()),
		ExpirationDate: models.JSONTime(time.Now()),
		Message:        fake.Sentence(),
	}
}

func PlayRecord(dur int) map[string]interface{} {
	return map[string]interface{}{
		"maxDuration": dur,
	}
}

func DB() (sqlmock.Sqlmock, *sql.DB) {
	mockDB, mock, _ := sqlmock.New()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	models.SetDB(sqlxDB)

	return mock, mockDB
}
