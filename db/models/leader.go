package models

import (
	"log"
	"os"

	"gitlab.com/arnaud-web/neli-webservices/config"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

// Leader represents a leader user
type Leader struct {
	*User
}

// IsOwned checks if user id passed in parameter is the leader's owner
func (l Leader) IsOwned(uid int64) bool {
	return l.ID == uid || db.Get(&User{}, `
		SELECT id
		FROM user 
		WHERE id = ? AND role = ?`,
		uid, AdminRole) == nil
}

// ToJSON clean properties for json result
func (l Leader) ToJSON() interface{} {
	l.Firstname = ""
	l.Lastname = ""
	l.Email = ""
	l.Password = ""
	return l
}

// Content get all content created by user
func (l Leader) Content() ([]Content, error) {
	c := []Content{}
	err := db.Select(&c, `
		SELECT 
			id, 
			IFNULL (name, '') AS name, 
			IFNULL (description, '') AS description, 
			DATE_FORMAT(created_at,'%Y-%m-%d') as creation_date,
			IFNULL (duration, 0) AS duration, 
			(SELECT COUNT(*) > 0 FROM share WHERE video_content.id = content_id) as sharing_status,
			IFNULL (duration, 0) != 0 AS ready
		FROM video_content 
		WHERE leader_id = ?
		ORDER BY created_at DESC`, l.ID)
	return c, err
}

// DeleteContent attach video content to zombie
func (l Leader) DeleteContent(cid int64) error {
	_, err := db.Exec("UPDATE video_content SET leader_id = ? WHERE id = ?", *config.ZombieID, cid)
	return err
}

// DeleteContent a leader and attach all videos to zombie
func (l Leader) Delete() error {
	if _, err := db.Exec("UPDATE video_content SET leader_id = ? WHERE leader_id = ?", *config.ZombieID, l.ID); err != nil {
		return err
	}

	_, err := db.NamedExec(`
		DELETE FROM user 
		WHERE id = :id`,
		&l)

	return err
}

func (l Leader) IsOwner(cid int64) bool {
	return db.Get(&Content{}, `
		SELECT id
		FROM video_content 
		WHERE id = ? AND leader_id = ?`,
		cid, l.ID) == nil
}

func (l Leader) EmailExists(lid int64) bool {
	err := db.Get(&Member{}, `
		SELECT user.id 
		FROM user 
		WHERE user.id != ? 
		AND email = ?`,
		l.ID,
		l.Email)

	if err == nil {
		return true
	}

	return false
}
