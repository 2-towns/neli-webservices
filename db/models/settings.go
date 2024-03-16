package models

// Settings represents general application settings
type Settings struct {
	ID          int64 `db:"id" json:"userId,omitempty"`
	MaxDuration int   `db:"max_duration" json:"maxDuration,omitempty"`
}

// Load settings from database
func (s *Settings) Load() error {
	return db.Get(s, "SELECT max_duration FROM settings LIMIT 1")
}

// Update settings in database
func (s *Settings) Update() error {
	_, err := db.NamedExec("UPDATE settings SET max_duration= :max_duration", &s)

	return err
}
