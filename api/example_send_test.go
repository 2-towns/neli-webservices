package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func Example_Send() {
	msg := struct {
		Message string
	}{
		"Test",
	}

	w := httptest.NewRecorder()

	Send(w, http.StatusOK, msg)

	fmt.Println(w.Body.String())
	// Output: {"Message":"Test"}
}
