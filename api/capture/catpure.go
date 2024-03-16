package capture

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/arnaud-web/neli-webservices/api"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

const ERROR = -1

const PLAYSTARTED = 1
const PLAYNOVIDEOINPUT = 2
const PLAYNOSERVERCONNECTION = 3
const PLAYNOBTCONNECTION = 4
const PLAYALREADYINPROGRESS = 5
const PLAYNOPROXYCONNECTION = 6

const STOPDONE = 1
const STOPNOBTCONNECTION = 2
const STOPNOCAPTUREINPROGRESS = 3
const STOPAUTHFAILURE = 4
const STOPNOPROXYCONNECTION = 5

const STATUSNOCAPTURE = 0
const STATUSCAPTUREINPROGRESS = 1
const STATUSNOPROXYCONNECTION = 2

type recordStatus struct {
	VideoContentId int64 `json:"videoContentId"`
	MaxDuration    int64 `json:"maxDuration"`
	EllapsedTime   int64 `json:"ellapsedTime"`
}

var playFunction = playRecord
var stopFunction = stopRecord
var statusFunction = getStatusRecord

type Request struct {
	MemberId int64
	Order    string
	Arg1     int64
	Arg2     int64
}

type Response struct {
	Result string
	Arg1   int64
	Arg2   int64
	Arg3   int64
}

//
//   SYNOPSIS:
//	 Order the N5 device to start a video transmission
//
//   INPUTS:
// 	memberId:		Id of the leader ordering the transmission
//	videoContentId:   	Id of the video to be created
//	maxDuration: 		Max duration time of the video
//
//   RETURN CODE:
//   	     One of the PLAY* const value

func playRecord(memberId int64, videoContentId int64, maxDuration int64) int16 {

	// Data
	jsonObject := Request{
		MemberId: memberId,
		Order:    "play",
		Arg1:     videoContentId,
		Arg2:     maxDuration,
	}

	// Request
	response := sendRequest(jsonObject)

	switch response.Result {
	// Results
	case "playstarted":
		return PLAYSTARTED
	case "nosource":
		return PLAYNOVIDEOINPUT
	case "noserver":
		return PLAYNOSERVERCONNECTION
	case "nobt":
		return PLAYNOBTCONNECTION
	case "busy":
		return PLAYALREADYINPROGRESS

		// Errors
	case "conn_error":
	case "json_error":
	case "timeout":
		return PLAYNOPROXYCONNECTION

	case "unknown":
		fmt.Println("Unknown code received from N7:", response.Arg1)
		return ERROR

	default:
		fmt.Println("Unhandled result received from Proxy:", response.Result)
		return ERROR
	}

	return ERROR
}

//
//   SYNOPSIS:
//	Order the N5 device to stop a video transmission
//
//   INPUT:
//	memberId:              Id of the leader ordering to stop the transmission
//
//   RETURN VALUES:
//   	     The duration of the stopped record.
//   	     One of the STOP* const value

func stopRecord(memberId int64) (int64, int16) {

	// Data
	jsonObject := Request{
		MemberId: memberId,
		Order:    "stop",
		Arg1:     0,
		Arg2:     0,
	}

	// Request
	response := sendRequest(jsonObject)

	switch response.Result {
	// Results
	case "stopdone":
		return response.Arg1, STOPDONE

	case "nobt":
		return 0, STOPNOBTCONNECTION
	case "idle":
		return 0, STOPNOCAPTUREINPROGRESS
	case "authfailure":
		return 0, STOPAUTHFAILURE

		// Errors
	case "conn_error":
	case "json_error":
	case "timeout":
		return 0, STOPNOPROXYCONNECTION

	case "unknown":
		fmt.Println("Unknown code received from N5:", response.Arg1)
		return 0, ERROR

	default:
		fmt.Println("Unhandled result received from Proxy:", response.Result)
		return 0, ERROR
	}

	return 0, ERROR
}

//
//   SYNOPSIS:
//	Get, from the N5 device, the status of the current video transmission
//
//   INPUT:
//      memberId:              Id of the leader asking for the transmission status
//
//   OUTPUT:
//	status:	               status of the current transmission
//
//   RETURN VALUE:
//           One of the STATUS* const value
func getStatusRecord(memberId int64, status *recordStatus) int16 {
	jsonObject := Request{
		MemberId: memberId,
		Order:    "status",
		Arg1:     0,
		Arg2:     0,
	}

	// Request
	response := sendRequest(jsonObject)

	switch response.Result {
	// Results
	case "status":
		status.VideoContentId = response.Arg1
		status.MaxDuration = response.Arg2
		status.EllapsedTime = response.Arg3
		return STATUSCAPTUREINPROGRESS

	case "nostatus":
		status.VideoContentId = 0
		status.MaxDuration = 0
		status.EllapsedTime = 0
		return STATUSNOCAPTURE

		// Errors
	case "conn_error":
	case "json_error":
	case "timeout":
		return STATUSNOPROXYCONNECTION

	case "unknown":
		fmt.Println("Unknown code received from N7:", response.Arg1)
		return ERROR

	default:
		fmt.Println("Unhandled result received from Proxy:", response.Result)
		return ERROR
	}

	return ERROR
}

func Play(w http.ResponseWriter, r *http.Request) {
	cid, _ := strconv.ParseInt(chi.URLParam(r, "videoContentId"), 10, 64)
	lid := api.UserIdFromContext(r)
	c := models.Content{ID: cid, LeaderID: lid}
	if err := c.Find(); err != nil {
		api.SendError(w, http.StatusNotFound, messages.InvalidVideoContentId)
		return
	}

	rs := recordStatus{}
	if err := json.NewDecoder(r.Body).Decode(&rs); err != nil {
		logger.Println(err)
		api.SendError(w, http.StatusBadRequest, messages.InvalidParameters)
		return
	}
	defer r.Body.Close()

	result := playFunction(lid, cid, rs.MaxDuration)
	switch result {
	case PLAYSTARTED:
		api.Send(w, http.StatusNoContent, "")
	case PLAYALREADYINPROGRESS:
		api.SendError(w, http.StatusForbidden, messages.CaptureAlreadyInProgress)
	case PLAYNOBTCONNECTION:
		api.SendError(w, http.StatusBadGateway, messages.NoBluetoothConnection)
	case PLAYNOVIDEOINPUT:
		api.SendError(w, http.StatusServiceUnavailable, messages.NoVideoInput)
	case PLAYNOSERVERCONNECTION:
		api.SendError(w, http.StatusGatewayTimeout, messages.Timeout)
	case PLAYNOPROXYCONNECTION:
		api.SendError(w, http.StatusGatewayTimeout, messages.NoProxyConnection)
	default:
		api.SendError(w, http.StatusInternalServerError, "")
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	lid := api.UserIdFromContext(r)
	_, result := stopFunction(lid)
	switch result {
	case STOPDONE:
		api.Send(w, http.StatusNoContent, nil)
	case STOPAUTHFAILURE:
		api.SendError(w, http.StatusUnauthorized, messages.InvalidToken)
	case STOPNOCAPTUREINPROGRESS:
		api.SendError(w, http.StatusForbidden, messages.CaptureAlreadyDone)
	case STOPNOBTCONNECTION:
		api.SendError(w, http.StatusBadGateway, messages.NoBluetoothConnection)
	case STOPNOPROXYCONNECTION:
		api.SendError(w, http.StatusGatewayTimeout, messages.NoProxyConnectionStopOrStatus)
	default:
		api.SendError(w, http.StatusInternalServerError, "")
	}
}

func Status(w http.ResponseWriter, r *http.Request) {
	lid := api.UserIdFromContext(r)
	rs := recordStatus{}
	result := statusFunction(lid, &rs)
	switch result {
	case STATUSCAPTUREINPROGRESS:
		api.Send(w, http.StatusOK, rs)
		return
	case STATUSNOCAPTURE:
		api.SendError(w, http.StatusNotFound, messages.NoCapture)
		return
	case STATUSNOPROXYCONNECTION:
		api.SendError(w, http.StatusGatewayTimeout, messages.NoProxyConnectionStopOrStatus)
		return
	default:
		api.SendError(w, http.StatusInternalServerError, "")
	}
}

func sendRequest(jsonObject Request) Response {

	jsonData, err := json.Marshal(jsonObject)
	if err != nil {
		fmt.Println("Json Marshal error:", err)
		return Response{Result: "json_error"}
	}

	// Request
	connection, err := net.Dial("udp", *config.ProxyBluetooth)
	if err != nil {
		fmt.Println("Could not resolve udp address or connect to it on", *config.ProxyBluetooth)
		fmt.Println(err)
		return Response{Result: "conn_error"}
	}

	n, err := connection.Write(jsonData)

	if err != nil {
		fmt.Println("Error writing data to server", *config.ProxyBluetooth)
		fmt.Println(err)
		return Response{Result: "conn_error"}
	}

	// Response
	recvBuf := make([]byte, 1024)
	connection.SetReadDeadline(time.Now().Add(1 * time.Second))
	n, err = connection.Read(recvBuf)
	if err != nil {
		fmt.Println(err)
		return Response{Result: "timeout"}
	}

	//fmt.Printf("Received data: %s\n", string(recvBuf[:n]))
	connection.Close()

	var response Response
	err = json.Unmarshal(recvBuf[:n], &response)

	if err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return Response{Result: "json_error"}
	}

	return response
}
