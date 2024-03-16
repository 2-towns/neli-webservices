package mail

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"

	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
)

var sendmail = smtp.SendMail

const (
	resetSubject  = "Subject: Reset your password"
	adminSubject  = "Subject: Your backoffice account"
	leaderSubject = "Subject: Your N7 account"
	shareSubject  = "FIRSTNAME LASTNAME want to share a video with you!"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// ResetPassword sends an email to an user with a token to reset his password.
func ResetPassword(token, recipient string) {
	b, err := ioutil.ReadFile("./__resources__/mails/reset.html") // just pass the file name

	if err != nil {
		logger.Print(err)
	}

	c := strings.Replace(string(b), "TOKEN", string(token), 1)
	c = strings.Replace(c, "URL_CONTENT", *config.URL, 2)
	c = strings.Replace(c, "URL", *config.URL, 1)

	send(c, resetSubject, recipient)
}

// Admin sends an email to new admin with his login and password
func Admin(pw, email string) {
	b, err := ioutil.ReadFile("./__resources__/mails/admin.html") // just pass the file name

	if err != nil {
		logger.Print(err)
	}

	c := strings.Replace(string(b), "LOGIN", email, 1)
	c = strings.Replace(c, "URL_CONTENT", *config.URL, 1)
	c = strings.Replace(c, "PASSWORD", pw, 1)

	send(c, adminSubject, email)
}

// Leader sends an email to new leader with an url to set his password
func Leader(token, email string) {
	b, err := ioutil.ReadFile("./__resources__/mails/leader.html") // just pass the file name

	if err != nil {
		logger.Print(err)
	}

	c := strings.Replace(string(b), "LOGIN", email, 1)
	c = strings.Replace(c, "URL_CONTENT", *config.URL, 1)
	c = strings.Replace(c, "TOKEN", token, 1)
	c = strings.Replace(c, "URL", *config.URL, 1)

	send(c, leaderSubject, email)
}

// Share sends an email to user with url generated
func Share(u *models.User, s *models.Share, c *models.Content, l *models.User) {
	b, err := ioutil.ReadFile("./__resources__/mails/share.html") // just pass the file name

	if err != nil {
		logger.Print(err)
	}

	cnt := strings.Replace(string(b), "URL_CONTENT", *config.URL, 1)
	cnt = strings.Replace(cnt, "URL", s.URL, 1)
	cnt = strings.Replace(cnt, "MESSAGE", s.Message, 1)
	cnt = strings.Replace(cnt, "DESCRIPTION", c.Description, 1)
	cnt = strings.Replace(cnt, "FIRSTNAME", l.Firstname, 2)
	cnt = strings.Replace(cnt, "LASTNAME", l.Lastname, 2)
	cnt = strings.Replace(cnt, "NAME", c.Name, 1)

	if *config.Stub == 1 {
		cnt = strings.Replace(cnt, "IMAGE", fmt.Sprintf("%s/public/%s.jpg", *config.URLContent, config.StubPath), 1)
	} else {
		cnt = strings.Replace(cnt, "IMAGE", fmt.Sprintf("%s/%d.jpg", *config.URLContent, s.ContentID), 1)
	}

	date := time.Time(s.ExpirationDate).Format("01 02, 2006 at 03:04 pm")
	cnt = strings.Replace(cnt, "DATE", date, 1)

	m := c.Duration / 60
	sec := c.Duration % 60

	secPretty := fmt.Sprintf("%d", sec)

	if sec < 10 {
		secPretty = fmt.Sprintf("0%s", secPretty)
	}

	// Add invalid character for template displaying in email client
	cnt = strings.Replace(cnt, "DURATION", fmt.Sprintf("%d:%s", m, secPretty), 1)

	subject := strings.Replace(shareSubject, "FIRSTNAME", l.Firstname, 2)
	subject = strings.Replace(subject, "LASTNAME", l.Lastname, 2)

	send(cnt, subject, u.Email)
}

func send(content, subject, email string) {
	from := mail.Address{Name: "n7", Address: *config.SMTPFrom}
	to := mail.Address{Name: "", Address: email}

	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	header["Date"] = time.Now().Format(time.RFC1123Z)

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(content))

	auth := smtp.PlainAuth(
		"",
		*config.SMTPLogin,
		*config.SMTPPassword,
		*config.SMTPHost,
	)

	err := sendmail(
		fmt.Sprintf("%s:%d", *config.SMTPHost, *config.SMTPPort),
		auth,
		*config.SMTPFrom,
		[]string{email},
		[]byte(message),
	)

	if err != nil {
		logger.Println(err)
	}
}
