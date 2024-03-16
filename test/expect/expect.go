package expect

import (
	"errors"
	"math/rand"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"gitlab.com/arnaud-web/neli-webservices/test/stub"
	"golang.org/x/crypto/bcrypt"
)

func AdminInsert(u models.User, mock sqlmock.Sqlmock) {
	userInsert(u, mock)
}

func LeaderInsert(u models.User, mock sqlmock.Sqlmock) {
	userInsert(u, mock)
}

func MemberInsert(u models.User, mock sqlmock.Sqlmock) {
	userInsert(u, mock)
}

func AdminOwner(uid int64, r string, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	if r == models.AdminRole {
		mock.ExpectQuery("SELECT (.+) FROM user (.+)").WithArgs(uid, r).WillReturnRows(rows)
	} else {
		mock.ExpectQuery("SELECT (.+) FROM user (.+)").WithArgs(uid, r).WillReturnError(errors.New(""))
	}
}

func LeaderOwner(uid, lid int64, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery("SELECT (.+) FROM tribe (.+)").WithArgs(uid, lid).WillReturnRows(rows)
}

func LeaderNotOwner(uid, lid int64, mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT (.+) FROM tribe (.+)").WithArgs(uid, lid).WillReturnError(errors.New("Pas gentil"))
}

func TribeInsert(u models.User, mock sqlmock.Sqlmock) {
	mock.ExpectExec("INSERT INTO tribe (.+)").
		WithArgs(sqlmock.AnyArg(), u.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func UserUpdate(a int64, mock sqlmock.Sqlmock) {
	mock.ExpectExec("UPDATE user SET (.+)").WillReturnResult(sqlmock.NewResult(1, a))
}

func UserDelete(u models.User, mock sqlmock.Sqlmock) {
	mock.ExpectExec("DELETE FROM user WHERE (.+)").WithArgs(u.ID).WillReturnResult(sqlmock.NewResult(1, 1))
}

func userInsert(u models.User, mock sqlmock.Sqlmock) {
	mock.ExpectExec("INSERT INTO user (.+)").
		WithArgs(u.Lastname, u.Firstname, u.Email, u.Role, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(u.ID, 1))
}

func FindByLogin(l, p string, mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.LeaderRole)
	pw, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	u.Password = string(pw)
	find(&u, mock).WithArgs(l, sqlmock.AnyArg())
	return u
}

func FindById(mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.LeaderRole)
	find(&u, mock).WithArgs(u.ID)
	return u
}

func FindByIdInParameter(uid int64, mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.LeaderRole)
	find(&u, mock).WithArgs(uid)
	return u
}

func FindMember(uid int64, mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.MemberRole)
	find(&u, mock).WithArgs(models.MemberRole, uid)
	return u
}

func FindByIdWithRole(r string, mock sqlmock.Sqlmock) models.User {
	u := stub.User(r)
	find(&u, mock).WithArgs(u.ID)
	return u
}

func NotFind(mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.LeaderRole)
	mock.ExpectQuery("SELECT (.+) FROM user WHERE (.+)").WillReturnRows(sqlmock.NewRows([]string{}))
	return u
}

func FindByRole(r string, mock sqlmock.Sqlmock) models.User {
	u := stub.User(r)
	find(&u, mock).WithArgs(u.Role)
	return u
}

func FindByPasswordToken(t string, mock sqlmock.Sqlmock) models.User {
	u := stub.User(models.LeaderRole)
	find(&u, mock).WithArgs(t)
	return u
}

func find(u *models.User, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	id := int64(rand.Uint64())
	rows := sqlmock.NewRows([]string{"id", "email", "password", "role", "firstname", "lastname"}).AddRow(id, u.Email, u.Password, u.Role, u.Firstname, u.Lastname)
	u.ID = id
	return mock.ExpectQuery("SELECT (.+) FROM user (.+)").WillReturnRows(rows)
}

func Tribe(t *models.Tribe, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	id := int64(rand.Uint64())
	rows := sqlmock.NewRows([]string{"leader_id"}).AddRow(id)
	t.LeaderID = id
	return mock.ExpectQuery("SELECT (.+) FROM tribe WHERE (.+)").WillReturnRows(rows)
}

func Settings(mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	rows := sqlmock.NewRows([]string{"max_duration"}).AddRow(3600)
	return mock.ExpectQuery("SELECT (.+) FROM settings LIMIT 1").WillReturnRows(rows)
}

func Content(lid int64, mock sqlmock.Sqlmock) int64 {
	id := int64(rand.Uint64())
	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery("SELECT (.+) FROM video_content (.+)").WithArgs(id, lid).WillReturnRows(rows)
	return id
}

func ContentLeader(lid int64, c models.Content, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "duration", "leader_id", "creation_date", "sharing_status"}).
		AddRow(c.ID, c.Name, c.Description, c.CreatedAt, c.Duration, c.LeaderID, time.Now().Format("20060102"), 1)
	mock.ExpectQuery("SELECT (.+) FROM video_content (.+)").WithArgs(lid).WillReturnRows(rows)
}

func SelectContent(lid int64, c models.Content, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "duration", "leader_id", "creation_date", "sharing_status"}).
		AddRow(c.ID, c.Name, c.Description, c.CreatedAt, c.Duration, c.LeaderID, time.Now().Format("20060102"), 1)
	mock.ExpectQuery("SELECT (.+) FROM video_content").WillReturnRows(rows)
}

func UpdateShareByleaderIdAndContentId(lid, cid int64, mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	return UpdateShare(mock).WithArgs(lid, cid)
}

func SelectShareById(lid int64, s models.Share, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	return SelectShare(s, mock).WithArgs(lid)
}

func SelectShareVideoContentId(cid int64, s models.Share, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	return SelectShare(s, mock).WithArgs(cid)
}

func SelectShare(s models.Share, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	rows := sqlmock.NewRows([]string{"user_id", "expiration_date", "message"}).
		AddRow(s.UserID, s.ExpirationDate, s.Message)
	return mock.ExpectQuery("SELECT (.+) FROM share").WillReturnRows(rows)
}

func UpdateShare(mock sqlmock.Sqlmock) *sqlmock.ExpectedExec {
	return mock.ExpectExec("UPDATE share (.+)").WillReturnResult(sqlmock.NewResult(int64(rand.Uint64()), 1))
}

func ContentNotFound(lid int64, mock sqlmock.Sqlmock) {
	rows := sqlmock.NewRows([]string{})
	mock.ExpectQuery("SELECT (.+) FROM video_content (.+)").WithArgs(sqlmock.AnyArg(), lid).WillReturnRows(rows)
}

func ContentInsert(uid int64, mock sqlmock.Sqlmock) int64 {
	id := int64(rand.Uint64())
	mock.ExpectExec("INSERT INTO video_content (.+)").
		WithArgs(uid).
		WillReturnResult(sqlmock.NewResult(id, 1))
	return id
}

func ContentUpdate(mock sqlmock.Sqlmock) int64 {
	max := int64(rand.Uint64())
	mock.ExpectExec("UPDATE video_content SET (.+)").WillReturnResult(sqlmock.NewResult(1, 1))
	return max
}

func SettingsUpdate(mock sqlmock.Sqlmock) int64 {
	max := int64(rand.Uint64())
	mock.ExpectExec("UPDATE settings SET (.+)").
		WithArgs(max).
		WillReturnResult(sqlmock.NewResult(1, 1))
	return max
}

func ShareInsert(mock sqlmock.Sqlmock) {
	mock.ExpectExec("INSERT INTO share (.+)").
		WillReturnResult(sqlmock.NewResult(1, 1))
}

func IsInTribe(lid, uid int64, mock sqlmock.Sqlmock) *sqlmock.ExpectedQuery {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	return mock.ExpectQuery("SELECT (.+) FROM tribe WHERE (.+)").WithArgs(lid, uid).WillReturnRows(rows)
}
