package models

import (
	"errors"

	"github.com/badoux/checkmail"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
)

// Member is a representation of an user in tribe.
type Member struct {
	*User
	Tribe Tribe `json:"-"`
}

// ToJSON clean properties for json result
func (m Member) ToJSON() interface{} {
	m.Firstname = ""
	m.Lastname = ""
	m.Email = ""
	m.Password = ""
	return m
}

// IsOwned checks if user id passed in parameter is the member's owner
func (m Member) IsOwned(uid int64) bool {
	return m.ID == uid || db.Get(&User{}, `
		SELECT id
		FROM tribe
		WHERE user_id = ? AND leader_id = ?`,
		m.ID, uid) == nil
}

// ListMembers return member list linked to leaderId parameter.
func ListMembers(leaderID int64) ([]Member, error) {
	var l []Member
	err := db.Select(
		&l,
		`SELECT user.id as id, firstname, lastname, email
		FROM user
		JOIN tribe on user.id = user_id WHERE role = ? AND leader_id = ?`,
		MemberRole,
		leaderID,
	)
	return l, err
}

// Find a member with his tribe
func (m *Member) Find(id int64) error {
	return db.Get(m, `
		SELECT user.id, role, IFNULL(leader_id, 0) AS 'tribe.leader_id'
		FROM user LEFT JOIN tribe as t ON t.user_id = user.id
		WHERE user.id = ?`, id)
}

// InTribe checks if member is in the tribe of leader passed in parameter
func (m Member) InTribe(lid int64) bool {
	if m.ID == lid {
		return true
	}

	return db.Get(&User{}, `
		SELECT id FROM tribe
		WHERE leader_id = ? AND user_id = ?`,
		lid, m.ID) == nil
}

// Validate checks user informations:
// * firstname has to be not empty
// * lastname has to be not empty
// * email has to be not empty and valid
func (m Member) Validate() error {
	if err := checkmail.ValidateFormat(m.Email); err != nil || (m.Firstname == "" && m.Lastname == "") {
		return errors.New(messages.InvalidUserInfo)
	}

	return nil
}

func (m Member) ExistsInTribe(lid int64) error {
	e := User{}

	err := db.Get(&e, `
		SELECT user.id 
		FROM user 
		JOIN tribe ON tribe.user_id = user.id AND tribe.leader_id = ? 
		WHERE email = ?`,
			lid,
		m.Email)


	if err != nil {
		return nil
	}

	return errors.New(messages.EmailExists)
}

func (m Member) EmailExists(lid int64) bool {
	err := db.Get(&Member{}, `
		SELECT user.id 
		FROM user 
		JOIN tribe ON tribe.user_id = user.id AND tribe.leader_id = ? 
		WHERE user.id != ? 
		AND email = ?`,
		lid,
		m.ID,
		m.Email)

	if err == nil {
		return true
	}

	return false
}