package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/badoux/checkmail"
	"gitlab.com/arnaud-web/neli-webservices/api/messages"
)

const (
	// SuperAdminRole is super administrator who manages admins
	SuperAdminRole = "super_admin"
	// AdminRole is administrator who manages leaders
	AdminRole = "administrator"
	// LeaderRole is leader who manages members
	LeaderRole = "leader"
	// MemberRole is user who access to content
	MemberRole = "member"
	// HubRole is hub which set a video content to ready
	HubRole = "hub"
)

// User representation
type User struct {
	ID                 int64     `db:"id" json:"userId"`
	Lastname           string    `db:"lastname" json:"lastName,omitempty"`
	Firstname          string    `db:"firstname" json:"firstName,omitempty"`
	Email              string    `db:"email" json:"email,omitempty"`
	Role               string    `db:"role" json:"profile,omitempty"`
	Password           string    `db:"password" json:"password,omitempty"`
	PasswordResetToken string    `db:"password_reset_token" json:"-"`
	CreatedAt          time.Time `db:"created_at" json:"-"`
	UpdatedAt          time.Time `db:"updated_at" json:"-"`
}

// Jsonify is an interface which represents an user, a leader, a member or an admin
type Jsonify interface {
	ToJSON() interface{}
	IsOwned(int64) bool
	Save() (sql.Result, error)
	Validate() error
	Delete() error
	EmailExists(lid int64) bool
}

// Find searches user by id
func (u *User) Find(id int64) error {
	return db.Get(u,
		`SELECT id, lastname, firstname, email, role, password 
		FROM user 
		WHERE id = ?`,
		id)
}

// Jsonify transforms an user to a jsonify object
func (u *User) Jsonify() (Jsonify, error) {
	c := *u

	if err := c.Find(u.ID); err != nil {
		return nil, err
	}

	u.Role = c.Role

	switch c.Role {
	case AdminRole:
		return Admin{u}, nil
	case LeaderRole:
		return Leader{u}, nil
	case MemberRole:
		return Member{u, Tribe{}}, nil
	default:
		return u, nil
	}
}

// FindByLogin search an user by login
func (u *User) FindByLogin(login string) error {
	return db.Get(u, `
		SELECT id, email, password, role 
		FROM user 
		WHERE email = ? AND role != ?`,
		login, MemberRole)
}

// FindByPasswordToken search an user by password token
func (u *User) FindByPasswordToken(t string) error {
	return db.Get(u, `
		SELECT id, email 
		FROM user 
		WHERE password_reset_token = ?`,
		t)
}

// UpdatePasswordToken update user password token
func (u User) UpdatePasswordToken(t string) error {
	_, err := db.Exec(`
		UPDATE user 
		SET password_reset_token = ? 
		WHERE id = ?`,
		t, u.ID)

	return err
}

// UpdatePassword update user password and set password token to null
func (u User) UpdatePassword(pw string) error {
	_, err := db.Exec(`
		UPDATE user 
		SET password = ?, password_reset_token = NULL  
		WHERE id = ?`,
		pw, u.ID)

	return err
}

// Save update user informations
func (u User) Save() (sql.Result, error) {
	return db.NamedExec(`
		UPDATE user 
		SET email = :email, firstname = :firstname, lastname = :lastname
		WHERE id = :id`,
		&u)
}

// New records user informations in database
func (u *User) New() error {
	r, err := db.Exec(
		`INSERT INTO user (lastname, firstname, email, role, password, password_reset_token, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())`,
		u.Lastname, u.Firstname, u.Email, u.Role, u.Password, u.PasswordResetToken)

	if err != nil {
		return err
	}

	u.ID, err = r.LastInsertId()
	return err

}

// ListUsers return list of user depending of role parameter
func ListUsers(role string) ([]User, error) {
	l := []User{}

	err := db.Select(&l, `
		SELECT user.id as id, firstname, lastname, email 
		FROM user WHERE role = ?`,
		role)

	if l == nil {
		l = []User{}
	}

	return l, err
}

// Delete an user
func (u User) Delete() error {
	_, err := db.NamedExec(`
		DELETE FROM user 
		WHERE id = :id`,
		&u)

	return err
}

// ToJSON clean properties for json result
func (u User) ToJSON() interface{} {
	u.Password = ""
	return u
}

// IsOwned checks if the uid in parameter is the owner of the
// current object.
func (u User) IsOwned(uid int64) bool {
	return u.ID == uid
}

// Validate checks user informations:
// * firstname has to be not empty
// * lastname has to be not empty
// * email has to be not empty and valid
func (u *User) Validate() error {
	if err := checkmail.ValidateFormat(u.Email); err != nil || u.Firstname == "" || u.Lastname == "" {
		return errors.New(messages.InvalidUserInfo)
	}

	return nil
}

func (u User) Exists() error {
	 e := User{}

	 err := db.Get(&e, `
		SELECT id 
		FROM user 
		WHERE email = ?`,
		u.Email)

	 if err != nil {
	 	return nil
	 }

	 return errors.New(messages.EmailExists)
}

func (u User) EmailExists(lid int64) bool {
	return false
}