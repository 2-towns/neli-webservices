package capture

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/test/assert"
	"gitlab.com/arnaud-web/neli-webservices/test/expect"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func TestPlay(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYSTARTED
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestPlayNotFound(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	expect.ContentNotFound(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", 123))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYNOVIDEOINPUT
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.InvalidVideoContentId)
}

func TestPlayNoInputVideo(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYNOVIDEOINPUT
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusServiceUnavailable, messages.NoVideoInput)
}

func TestPlayInProgress(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYALREADYINPROGRESS
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.CaptureAlreadyInProgress)
}

func TestPlayNotBluetoothConnection(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYNOBTCONNECTION
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusBadGateway, messages.NoBluetoothConnection)
}

func TestPlayNoServerConnection(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return PLAYNOSERVERCONNECTION
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusGatewayTimeout, messages.Timeout)
}

func TestPlayInternalError(t *testing.T) {
	mock, mockDB := stub.DB()
	defer mockDB.Close()

	uid := int64(rand.Uint64())
	cid := expect.Content(uid, mock)

	p := stub.PlayRecord(3600)

	rr, r := stub.PostHttpWithContext(uid, p)
	r = stub.AddUrlParam(r, "videoContentId", fmt.Sprintf("%d", cid))

	playFunction = func(memberId int64, videoContentId int64, maxDuration int64) int16 {
		return 0
	}

	Play(rr, r)

	assert.Response(t, rr, http.StatusInternalServerError, "")
}

func TestStop(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, map[string]interface{}{})

	stopFunction = func(memberId int64) (int64, int16) {
		return 9, STOPDONE
	}

	Stop(rr, r)

	assert.Response(t, rr, http.StatusNoContent, "")
}

func TestStopAuthenticationFailure(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, map[string]interface{}{})

	stopFunction = func(memberId int64) (int64, int16) {
		return 0, STOPAUTHFAILURE
	}

	Stop(rr, r)

	assert.Response(t, rr, http.StatusUnauthorized, messages.InvalidToken)
}

func TestStopForbidden(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, map[string]interface{}{})

	stopFunction = func(memberId int64) (int64, int16) {
		return 0, STOPNOCAPTUREINPROGRESS
	}

	Stop(rr, r)

	assert.Response(t, rr, http.StatusForbidden, messages.CaptureAlreadyInProgress)
}

func TestStopNoBluetoothConnection(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, map[string]interface{}{})

	stopFunction = func(memberId int64) (int64, int16) {
		return 0, STOPNOBTCONNECTION
	}

	Stop(rr, r)

	assert.Response(t, rr, http.StatusBadGateway, messages.NoBluetoothConnection)
}

func TestStopInternalError(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.PostHttpWithContext(uid, map[string]interface{}{})

	stopFunction = func(memberId int64) (int64, int16) {
		return 0, -1
	}

	Stop(rr, r)

	assert.Response(t, rr, http.StatusInternalServerError, "")
}

func TestStatus(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.HttpWithContext(uid)

	statusFunction = func(memberId int64, status *recordStatus) int16 {
		return STATUSCAPTUREINPROGRESS
	}

	Status(rr, r)

	assert.Response(t, rr, http.StatusOK, "")
	assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", 0))
	assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", 0))
}

func TestStatusNoCapture(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.HttpWithContext(uid)

	statusFunction = func(memberId int64, status *recordStatus) int16 {
		return STATUSNOCAPTURE
	}

	Status(rr, r)

	assert.Response(t, rr, http.StatusNotFound, messages.NoCapture)
}

func TestStatusInternalError(t *testing.T) {
	uid := int64(rand.Uint64())

	rr, r := stub.HttpWithContext(uid)

	statusFunction = func(memberId int64, status *recordStatus) int16 {
		return 3
	}

	Status(rr, r)

	assert.Response(t, rr, http.StatusInternalServerError, "")
}
