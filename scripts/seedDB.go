package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/go-sql-driver/mysql"
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
}

func configureDatabase(db *sql.DB) {
	db.SetMaxOpenConns(151)
	db.SetMaxIdleConns(100)
}


// Generates password of length 16
func passwordGenerator() string {
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
    str := make([]rune, 16)
    for i := range str {
        str[i] = chars[rand.Intn(len(chars))]
    }
    return string(str)
}

func seedDatabase(db *sql.DB) {
	for i := 0; i < 10000; i++ {
		query := `INSERT INTO user_tab (username, password, nickname, profile_picture) values `
		for j := 0; j < 1000; j++ {
			username := fmt.Sprintf("user-%d", i*j + 1)
			nickname := fmt.Sprintf("nickname-%d", i*j + 1)
			password := passwordGenerator()
			password_hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			profile_picture := "resources/default-profile-picture.jpeg"
			row := fmt.Sprintf(`("%s", "%s", "%s", "%s")`, username, password_hash, nickname, profile_picture)
			query = query + row
			if j < 999 {
				query = query + `, `
			}
		}
		fmt.Println(query)
		_, err := db.Exec(query)
			if err != nil {
				log.Printf("Error %s when inserting user\n", err)
			}
	}
}

func main() {
	db, err := sql.Open("mysql", dsn("mysql"))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	createDatabase(db)
	configureDatabase(db)
	seedDatabase(db)
}