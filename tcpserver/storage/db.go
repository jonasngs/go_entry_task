package storage

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/jonasngs/go_entry_task/tcpserver/config"
	_ "github.com/go-sql-driver/mysql"
)

// const (
// 	username = "root"
// 	password = "garena"
// 	hostname = "localhost:3306"
// 	dbname   = "go_entry_task_db"
// )

type Database struct {
	DB *sql.DB
}

var dbConfig config.MySQL

func dsn(dbName string, username string, password string, hostname string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func createDatabase(db *sql.DB) {
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_bin",
		dbConfig.DB_name)
	if _, err := db.Exec(createDBQuery); err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}

	if _, err := db.Exec("USE " + dbConfig.DB_name); err != nil {
		log.Printf("Error %s when using DB\n", err)
	}

}

func configureDatabase(db *sql.DB) {
	db.SetMaxOpenConns(dbConfig.Max_connections)
	db.SetMaxIdleConns(dbConfig.Max_idle_connections)
	// db.SetConnMaxLifetime(time.Minute * 5)
}

func InitializeDatabase() Database {

	dbConfig = config.GetMYSQLServer()

	//db, err := sql.Open("mysql", dsn("mysql"))
	connectionURL := dsn(dbConfig.Driver, dbConfig.Username, dbConfig.Password, dbConfig.Hostname)
	db, err := sql.Open(dbConfig.Driver, connectionURL)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	createDatabase(db)
	configureDatabase(db)
	return Database{DB: db}
}
