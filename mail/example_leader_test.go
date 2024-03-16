package mail

import (
	"fmt"
	"net/smtp"

	"github.com/icrowley/fake"
)

func ExampleLeader() {
	sendmail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error { return nil }

	Leader(fake.SimplePassword(), fake.EmailAddress())

	fmt.Println(leaderSubject)
	// Output: Subject: Your N7 account
}
