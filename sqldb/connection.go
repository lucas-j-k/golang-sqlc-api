package sqldb

import (
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
)

// ////
// Export a pointer to a DB connection
// ////
var DB *sql.DB

// Connect to the DB, set connection in exported var
func Connect() error {

	sqlPass := viper.Get("SQL_PASS")
	sqlHost := viper.Get("SQL_HOST")
	sqlUser := viper.Get("SQL_USER")

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:3306)/questions_db?parseTime=True", sqlUser, sqlPass, sqlHost)
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		return err
	}

	DB = db
	return nil
}
