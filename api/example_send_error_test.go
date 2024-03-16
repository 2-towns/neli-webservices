package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func Example_SendError() {
	w := httptest.NewRecorder()

	SendError(w, http.StatusInternalServerError, "Test")

	fmt.Println(w.Body.String())
	// Output: {"code":500,"message":"Test"}
}
