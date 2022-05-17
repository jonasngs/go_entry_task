package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

const (
	username = "root"
	password = "garena"
	hostname = "localhost:3306"
	dbname   = "test_db"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createDatabase(db *sql.DB) {
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_bin",
		"go_entry_task_db")
	if _, err := db.Exec(createDBQuery); err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}

	if _, err := db.Exec("USE " + "go_entry_task_db"); err != nil {
		log.Printf("Error %s when using DB\n", err)
	}

	query := `CREATE TABLE IF NOT EXISTS user_tab` +
		`(user_id BIGINT UNSIGNED AUTO_INCREMENT NOT NULL,
					username VARCHAR(32) NOT NULL,
					password VARCHAR(64) NOT NULL,
					nickname VARCHAR(32) NOT NULL DEFAULT '', 
					profile_picture VARCHAR(160) NOT NULL DEFAULT '',
					created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
					PRIMARY KEY(user_id)
					) ENGINE = INNODB COLLATE utf8mb4_unicode_ci`
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Error %s when creating Table\n", err)
	}

	query = `INSERT INTO user_tab (username, password, nickname, profile_picture) values ("user123", "password2", "nickname123", "profile_picture123")`
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error %s when inserting user\n", err)
	}

	query = "SELECT user_id, username, password, nickname, profile_picture FROM user_tab WHERE username = ?"
	row := db.QueryRow(query, "user67")
	var user_id uint64
	var password string
	var nickname string
	var username string
	var profile_picture string

	err = row.Scan(&user_id, &username, &password, &nickname, &profile_picture)
	if err != nil {
		log.Printf("Error %s when Fetching user\n", err)
	}
	// fmt.Print(username)
	// fmt.Print(user_id)
	// fmt.Print(password)
	// fmt.Print(nickname)
	// fmt.Print(profile_picture)
}

func configureDatabase(db *sql.DB) {
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)
}

func insert(db *sql.DB) {
	hash1, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
	// hash2, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
	//fmt.Println(len(hash1))
	query := fmt.Sprintf(`INSERT INTO user_tab (username, password, nickname, profile_picture) values ("user30", "%s", "nickname30", "profile_picture30")`, string(hash1))
	_, err = db.Exec(query)
	if err != nil {
		log.Printf("Error %s when inserting user\n", err)
	}
	// fmt.Println(err)
	// test := bcrypt.CompareHashAndPassword(hash1, []byte("akncjkasjk"))
	//fmt.Println(test == nil)
}

func updateNickName(db *sql.DB) error {
	query := `UPDATE user_tab SET nickname=? WHERE username=?`
	_, err := db.Exec(query, "NEW", "user64")
	if err != nil {
		log.Printf("Error %s when updateing nickname\n", err)
	}
	return err
}

func main() {
	db, err := sql.Open("mysql", dsn("mysql"))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	createDatabase(db)
	configureDatabase(db)
	insert(db)
	//updateNickName(db)
}
