package models

import (
	"time"

	"gitlab.com/arnaud-web/neli-webservices/config"
)

// Content represents a video content
type Content struct {
	ID            int64     `db:"id" json:"videoContentId"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	Duration      int       `db:"duration" json:"duration"`
	LeaderID      int64     `db:"leader_id" json:"leaderId,omitempty"`
	MaxDuration   int       `json:"maxDuration,omitempty"`
	Path          string    `db:"path" json:"-"`
	Ready         bool      `db:"ready" json:"ready"`
	SharingStatus bool      `db:"sharing_status" json:"sharingStatus,omitempty"`
	CreationDate  string    `db:"creation_date" json:"creationDate,omitempty"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`
}

// New create a new content in database
func (c *Content) New() error {
	r, err := db.Exec(`
		INSERT INTO video_content (leader_id, created_at)
		VALUES (?, NOW())`,
		c.LeaderID)

	if err != nil {
		return err
	}

	c.ID, err = r.LastInsertId()

	if *config.Stub == 1 {
		_, err = db.Exec(`
			UPDATE video_content
			SET path = ?
			WHERE id = ? `,
			config.StubPath, c.ID)

		if err != nil {
			return err
		}
	}

	return err
}

// Find a content by id
func (c *Content) Find() error {
	return db.Get(c, `
		SELECT 
			id, IFNULL(name,'') AS name, 
			IFNULL(description,'') AS description, 
			IFNULL(path,'') AS path,
			IFNULL(duration, 0) AS duration 
		FROM video_content
		WHERE id = ? AND leader_id = ?`,
		c.ID, c.LeaderID)
}

// Update a content
func (c *Content) Update() error {
	_, err := db.NamedExec(`
		UPDATE video_content
		SET name = :name, description = :description
		WHERE id = :id `,
		&c)

	return err
}

// IsNotReady return true if the video path has not been updated.
// This update can be done from hub or by stube mode.
func (c *Content) IsNotReady() bool {
	return c.Path == ""
}

// UpdatePath update only path field for content
func (c Content) UpdatePath() error {
	_, err := db.Exec(`
		UPDATE video_content
		SET path = ?, duration = ?
		WHERE id = ? `,
		c.Path, c.Duration, c.ID)

	return err
}

func (c *Content) Shares() ([]Share, error) {
	var l []Share
	err := db.Select(&l, `
			SELECT share.id, user_id, expiration_date, message, CONCAT(firstname, " ", lastname) as name
			FROM share
			JOIN user ON user.id = user_id
			WHERE share.content_id = ? and url = ""`, c.ID)

	return l, err
}

func Clean(d string) error {
	_, err := db.Exec(`INSERT INTO cleaning (date, done) VALUES (?, 0)`, d)
	return err
}
