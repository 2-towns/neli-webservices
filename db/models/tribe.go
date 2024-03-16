package models

import "database/sql"

// Tribe representation.
type Tribe struct {
	LeaderID int64 `db:"leader_id" json:"leader_id,omitempty"`
	UserID   int64 `db:"user_id" json:"user_id,omitempty"`
	ID       int64 `db:"id" json:"id"`
}

// New create a new tribe
func (t *Tribe) New() (sql.Result, error) {
	return db.NamedExec("INSERT INTO tribe (leader_id, user_id) VALUES (:leader_id, :user_id)", &t)
}
