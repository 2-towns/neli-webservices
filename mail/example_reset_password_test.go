package mail

import (
	"fmt"
	"net/smtp"

	"github.com/icrowley/fake"
)

func ExampleResetPassword() {
	sendmail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error { return nil }

	ResetPassword(fake.Sentence(), fake.EmailAddress())

	fmt.Println(resetSubject)
	// Output: Subject: Reset your password
}
