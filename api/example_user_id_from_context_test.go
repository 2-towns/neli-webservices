package api

import (
	"fmt"

	"gitlab.com/arnaud-web/neli-webservices/test/stub"
)

func Example_UserIdFromContext() {
	uid := int64(123)
	_, r := stub.HttpWithContext(uid)

	fmt.Println(UserIdFromContext(r))
	// Output: 123
}
