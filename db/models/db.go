package models

import (

	// Used to load mysql driver

	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gitlab.com/arnaud-web/neli-webservices/config"
)

var db *sqlx.DB

// init create a database access for models package
func init() {

	if *config.DB != "" {
		var err error
		db, err = sqlx.Connect("mysql", *config.DB)

		if err != nil {
			log.Fatalln(err)
		}
	}

}

// SetDB provide a setter useful for testing.
func SetDB(d *sqlx.DB) {
	db = d
}
