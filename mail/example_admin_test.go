package mail

import (
	"fmt"
	"net/smtp"

	"github.com/icrowley/fake"
)

func ExampleAdmin() {
	sendmail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error { return nil }

	Admin(fake.SimplePassword(), fake.EmailAddress())

	fmt.Println(adminSubject)
	// Output: Subject: Your backoffice account
}
