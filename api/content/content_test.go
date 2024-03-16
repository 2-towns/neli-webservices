package content

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func TestNew(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())

	cid := expect.ContentInsert(uid, mock)
	expect.Settings(mock)

	rr, r := stub.HttpWithContext(uid)

	New(rr, r)

	assert.Response(t, rr, http.StatusOK, "")
	assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", cid))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", 3600))
}

func TestMaxDuration(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	m := expect.SettingsUpdate(mock)

	p := map[string]interface{}{
		"maxDuration": m,
	}

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)

	MaxDuration(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestMaxDuration_MaxDurationInvalid(t *testing.T) {
	p := map[string]interface{}{
		"maxDuration": fake.Sentence(),
	}

	uid := int64(rand.Uint64())
	rr, r := stub.PostHttpWithContext(uid, p)

	MaxDuration(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestEdit(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())
	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.VideoParam(fake.Sentence(), fake.Sentence())

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", lid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	expect.ContentUpdate(mock)

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestEdit_FirstnameEmpty(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())
	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.VideoParam("", fake.Sentence())

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", lid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestEdit_LastnameEmpty(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())
	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.VideoParam(fake.Sentence(), "")

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", lid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestEdit_ContentNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())
	uid := int64(rand.Uint64())
	cid := expect.Content(lid, mock)

	p := stub.VideoParam(fake.Sentence(), fake.Sentence())

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", lid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Edit(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.InvalidVideoContentIdOrUserId)
}

func TestShare(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	lid := int64(rand.Uint64())
	cid := expect.Content(lid, mock)

	p := stub.ShareParam(fake.Sentence(), time.Now().Format(time.RFC3339))

	rr, r := stub.PostHttpWithContext(lid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", uid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	expect.IsInTribe(lid, uid, mock)
	expect.ShareInsert(mock)
	expect.FindByIdInParameter(uid, mock)

	Share(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestShare_NotIntTribe(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	lid := int64(rand.Uint64())
	cid := expect.Content(lid, mock)

	p := stub.ShareParam(fake.Sentence(), time.Now().Format(time.RFC3339))

	rr, r := stub.PostHttpWithContext(lid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", uid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	expect.IsInTribe(uid, uid, mock)

	Share(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.InvalidVideoContentIdOrUserId)
}

func TestShare_MessageEmpty(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	lid := int64(rand.Uint64())
	cid := expect.Content(lid, mock)

	p := stub.ShareParam("", time.Now().Format(time.RFC3339))

	rr, r := stub.PostHttpWithContext(lid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", uid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	expect.IsInTribe(lid, uid, mock)
	expect.ShareInsert(mock)
	expect.FindByIdInParameter(uid, mock)

	Share(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestShare_BadDateFormat(t *testing.T) {
	uid := int64(rand.Uint64())
	lid := int64(rand.Uint64())

	p := stub.ShareParam("", time.Now().Format(time.RFC1123Z))

	rr, r := stub.PostHttpWithContext(lid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", uid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", uid))

	Share(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}

func TestShare_ContentNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	lid := int64(rand.Uint64())
	expect.ContentNotFound(lid, mock)

	p := stub.ShareParam(fake.Sentence(), time.Now().Format(time.RFC3339))

	rr, r := stub.PostHttpWithContext(lid, p)
	r = stub.AddUrlParam(r, "userId", fmt.Sprintf("%d", uid))
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", uid))

	Share(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.InvalidVideoContentIdOrUserId)
}

func TestList(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())

	c := stub.Content()
	expect.ContentLeader(lid, c, mock)

	rr, r := stub.HttpWithContext(lid)

	Get(rr, r)

	b := rr.Body.String()

	assert.Response(t, rr, http.StatusOK, "")

	assert.Contains(t, b, fmt.Sprintf("%d", c.ID))
	assert.Contains(t, b, c.Name)
	assert.Contains(t, b, c.Description)
	assert.Contains(t, b, c.CreationDate)
	assert.Contains(t, b, fmt.Sprintf("%d", c.Duration))
	assert.Contains(t, b, fmt.Sprintf("%v", c.SharingStatus))
}

func TestAll(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())

	c := stub.Content()
	expect.SelectContent(lid, c, mock)

	rr, r := stub.HttpWithContext(lid)

	All(rr, r)

	b := rr.Body.String()

	assert.Response(t, rr, http.StatusOK, "")

	assert.Contains(t, b, fmt.Sprintf("%d", c.ID))
	assert.Contains(t, b, c.Name)
	assert.Contains(t, b, c.Description)
	assert.Contains(t, b, fmt.Sprintf("%d", c.Duration))
	assert.Contains(t, b, fmt.Sprintf("%d", c.LeaderID))
	assert.Contains(t, b, c.CreationDate)
	assert.Contains(t, b, fmt.Sprintf("%v", c.SharingStatus))
}

func TestShareList(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	cid := int64(rand.Uint64())

	s := stub.Share()
	expect.SelectShareVideoContentId(cid, s, mock)

	rr, r := stub.Http()
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Shares(rr, r)

	b := rr.Body.String()

	assert.Response(t, rr, http.StatusOK, "")

	assert.Contains(t, b, fmt.Sprintf("%d", s.UserID))
	assert.Contains(t, b, s.Message)
}

func TestDelete(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	lid := int64(rand.Uint64())
	cid := expect.Content(lid, mock)

	expect.ContentUpdate(mock)

	rr, r := stub.HttpWithContext(lid)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Delete(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestReady(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	p := map[string]interface{}{
		"randomString": fake.Sentence(),
		"duration":     10,
	}

	uid := int64(rand.Uint64())
	cid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, p)

	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	expect.Settings(mock)
	expect.SelectShare(models.Share{
		UserID: uid,
	}, mock)
	expect.UpdateShare(mock)
	expect.ContentUpdate(mock)

	Ready(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestReady_EmptyRandomString(t *testing.T) {
	p := map[string]interface{}{
		"randomString": "",
	}

	uid := int64(rand.Uint64())
	cid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, p)

	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	Ready(rr, r)

	assert.Response(t, rr, http.StatusBadRequest, messages.InvalidParameters)
}
