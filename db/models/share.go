package models

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"time"

	"gitlab.com/arnaud-web/neli-webservices/config"
)

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	b := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(b), nil
}

// Share representation.
type Share struct {
	ID             int64     `db:"id" json:"id,omitempty"`
	UserID         int64     `db:"user_id" json:"memberId,omitempty"`
	UserName       string    `db:"name" json:"name,omitempty"`
	ContentID      int64     `db:"content_id" json:"content_id,omitempty"`
	URL            string    `db:"url" json:"url,omitempty"`
	Message        string    `db:"message" json:"message,omitempty"`
	ExpirationDate JSONTime  `db:"expiration_date" json:"expirationDate,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"-"`
}

// New create a new share
func (s *Share) New(c *Content) error {
	err := s.generateURL(c)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
			INSERT INTO share (user_id, content_id, url, expiration_date, message) 
			VALUES (?, ?, ?, ?, ?)`,
		s.UserID, s.ContentID, s.URL, time.Time(s.ExpirationDate), s.Message)

	return err
}

func (s *Share) generateURL(c *Content) error {
	// If content is not ready so don't generate url.
	// It will be set later, when the hub will call the api.
	if c.IsNotReady() {
		return nil
	}

	e, err := s.encodeBase64(c)

	if err != nil {
		return err
	}

	s.URL = fmt.Sprintf("%s/%d/video/%s", *config.URL, s.UserID, e)

	return nil
}

// EncodeBase64 transform a content to a string in base64
func (s *Share) encodeBase64(c *Content) (string, error) {
	r := struct {
		Message        string
		ExpirationDate time.Time
		Path           string
		Name           string
	}{
		s.Message,
		time.Time(s.ExpirationDate),
		c.Path,
		c.Name,
	}

	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(r)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b.Bytes()), nil
}

// IsReady return true if the share url is generated
func (s *Share) IsReady() bool {
	return s.URL != ""
}

// Find share by id
func (s *Share) Find() error {
	return db.Get(s, `
		SELECT url 
		FROM share 
		WHERE user_id = ? AND content_id = ?`,
		s.UserID, s.ContentID)
}

// List get all shares
func (s *Share) List(cid int64) ([]Share, error) {
	var l []Share
	err := db.Select(&l, `
		SELECT user_id, expiration_date, message, CONCAT(firstname, " ", lastname) as name
		FROM share
		JOIN video_content ON share.content_id = video_content.id
		JOIN user ON share.user_id = user.id
		WHERE share.content_id = ?`, cid)

	return l, err
}

// UpdateURL update only url field for share
func (s *Share) UpdateURL(c *Content) error {
	err := s.generateURL(c)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE share
		SET url = ?
		WHERE id = ? `,
		s.URL, s.ID)

	return err
}
