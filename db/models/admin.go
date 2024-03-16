package models

// Admin represents an admin user
type Admin struct {
	*User
}

// ToJSON clean properties for json result
func (a Admin) ToJSON() interface{} {
	a.Firstname = ""
	a.Lastname = ""
	a.Password = ""
	return a
}

// Content returns all content created.
// Content needs to have a name and description.
func (a Admin) Content() ([]Content, error) {
	var c []Content
	err := db.Select(&c, `
		SELECT 
			id, 
			IFNULL (name, '') AS name, 
			IFNULL (description, '') AS description, 
			DATE_FORMAT(created_at,'%Y-%m-%d') as creation_date,
			IFNULL (duration, 0) AS duration,
			leader_id,
			(SELECT COUNT(*) > 0 FROM share WHERE video_content.id = content_id) as sharing_status,
			IFNULL (duration, 0) != 0 AS ready
		FROM video_content`)

	return c, err
}

// IsOwned checks if user id passed in parameter is super admin
func (a Admin) IsOwned(uid int64) bool {
	return a.ID == uid || db.Get(&User{}, `
		SELECT id
		FROM user 
		WHERE id = ? AND role = ?`,
		uid, SuperAdminRole) == nil
}
