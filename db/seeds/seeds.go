package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gitlab.com/arnaud-web/neli-webservices/api/auth/bearer"
	"gitlab.com/arnaud-web/neli-webservices/config"
	"gitlab.com/arnaud-web/neli-webservices/db/models"
	"golang.org/x/crypto/bcrypt"
)

type jsonObject struct {
	Users    []models.User
	Tribes   []models.Tribe
	Settings models.Settings `json:"settings"`
	Content  []models.Content
	Shares   []models.Share
}

type postman struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Values []value `json:"values"`
	Scope  string  `json:"_postman_variable_scope"`
}

type value struct {
	Enabled bool   `json:"enabled"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	Type    string `json:"type"`
}

var db *sqlx.DB

func main() {
	var err error
	db, err = sqlx.Connect("mysql", *config.DB)
	if err != nil {
		log.Fatalln(err)
	}

	if *config.Env == "production" {
		// Search for zombie
		u := models.User{}

		if err := u.Find(1); err == nil {
			log.Println("Seeder already initialized")
			return
		}
	}

	clearData()

	data := loadData()

	v := []value{}

	log.Println("Saving users")
	v = saveUsers(data.Users, v)

	log.Println("Saving tribes")
	saveTribes(data.Tribes)

	log.Println("Saving settings")
	saveSettings(data.Settings)

	log.Println("Saving content")
	saveContent(data.Content)

	log.Println("Saving share")
	saveShares(data.Shares)

	savePostman(v)
}

func clearData() {
	if _, err := db.Exec("DELETE FROM tribe"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("DELETE FROM user"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("DELETE FROM settings"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("DELETE FROM video_content"); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("DELETE FROM share"); err != nil {
		log.Fatal(err)
	}
}

func loadData() jsonObject {
	g, _ := os.Getwd()
	raw, err := ioutil.ReadFile(fmt.Sprintf("%s/__resources__/seeds/%s/seeds.json", g, *config.Env))
	if err != nil {
		log.Fatal(err)
	}

	var data jsonObject
	json.Unmarshal(raw, &data)

	return data
}

func saveUsers(users []models.User, v []value) []value {
	roles := ""

	for _, u := range users {
		p := u.Password

		if u.Password != "" {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
			u.Password = string(hashedPassword)
		}

		if _, err := db.NamedExec("INSERT INTO user (id, email, role, password, firstname, lastname) VALUES (:id, :email, :role, :password, :firstname, :lastname)", u); err != nil {
			log.Fatal(err)
		}

		if !strings.Contains(roles, u.Role) {
			r, _ := bearer.Create(u, 0)

			v = append(v, value{Enabled: true, Key: fmt.Sprintf("%s_access_token", u.Role), Value: r.AccessToken, Type: "text"})

			if u.Role == models.SuperAdminRole {
				v = append(v, value{Enabled: true, Key: "login", Value: u.Email, Type: "text"})
				v = append(v, value{Enabled: true, Key: "password", Value: p, Type: "text"})
				v = append(v, value{Enabled: true, Key: "refresh_token", Value: r.RefreshToken, Type: "text"})
			}
			roles += u.Role
		}
	}
	return v
}

func saveTribes(t []models.Tribe) {
	for _, tribe := range t {
		if _, err := db.NamedExec("INSERT INTO tribe (user_id, leader_id) VALUES (:user_id, :leader_id)", tribe); err != nil {
			log.Fatal(err)
		}
	}
}

func saveSettings(s models.Settings) {
	if _, err := db.NamedExec("INSERT INTO settings (max_duration) VALUES (:max_duration)", s); err != nil {
		log.Fatal(err)
	}
}

func saveContent(l []models.Content) {
	for _, c := range l {
		c.Path = "123456"
		if _, err := db.Exec(
			`INSERT INTO video_content (id, name, description, path, duration, leader_id, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, NOW())`,
			c.ID, c.Name, c.Description, c.Path, c.Duration, c.LeaderID); err != nil {
			log.Fatal(err)
		}
	}
}

func saveShares(l []models.Share) {
	for _, s := range l {
		s.ExpirationDate = models.JSONTime(time.Now())
		if _, err := db.Exec(
			`INSERT INTO share (id, url, user_id, content_id, expiration_date, message,created_at) 
			VALUES (?, ?, ?, ?, ?, ?, NOW())`,
			s.ID, s.URL, s.UserID, s.ContentID, time.Time(s.ExpirationDate), s.Message); err != nil {
			log.Fatal(err)
		}
	}
}

func savePostman(v []value) {
	v = append(v, value{Enabled: true, Key: "token_life", Value: fmt.Sprintf("%d", *config.TokenLife), Type: "text"})
	v = append(v, value{Enabled: true, Key: "url", Value: "http://localhost:8082", Type: "text"})
	v = append(v, value{Enabled: true, Key: "mailhog", Value: "http://localhost:8025/api/v2/messages", Type: "text"})
	v = append(v, value{Enabled: true, Key: "expired_token", Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2LCJyb2xlIjoiYWRtaW4iLCJSZWZyZXNoVG9rZW4iOnRydWUsImV4cCI6MTUwNjAwNjU1NH0.FuheX3CDvRTdnyFAr5-qCqrvMBoXblAtLquGgjyq7bQ", Type: "text"})

	p := postman{ID: "29d652b9-531b-cdf9-3d88-c97f9edb99ad", Name: "local", Values: v}
	b, _ := json.Marshal(p)

	g, _ := os.Getwd()
	err := ioutil.WriteFile(fmt.Sprintf("%s/__resources__/postman/local.postman_environment.json", g), b, 0644)
	if err != nil {
		log.Println(err)
	}
}
