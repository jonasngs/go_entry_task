package services

import (
	"log"

	pb "github.com/jonasngs/go_entry_task/grpc"
	"github.com/jonasngs/go_entry_task/tcpserver/storage"
)

type DAOInterface interface {
	fetchUser(username string) (*pb.User, error)
	updateNickName(username string, nickname string) error
	updateProfilePicture(username string, profile_picture string) error
}

type DAOservice struct {
	db storage.Database
}

func InitializeDAOservice(db storage.Database) DAOInterface {
	createTable(db)
	return DAOservice{db: db}
}

func createTable(db storage.Database) {
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
	_, err := db.DB.Exec(query)
	if err != nil {
		log.Printf("Error %s when creating Table\n", err)
	}
}

func (ds DAOservice) fetchUser(username string) (*pb.User, error) {

	query := "SELECT user_id, username, password, nickname, profile_picture FROM user_tab WHERE username = ?"
	row := ds.db.DB.QueryRow(query, username)
	var user_id uint64
	var password string
	var nickname string
	var profile_picture string

	err := row.Scan(&user_id, &username, &password, &nickname, &profile_picture)
	if err != nil {
		log.Printf("Error %s when Fetching user\n", err)
		return nil, err
	}

	return &pb.User{
		UserId:         user_id,
		Username:       username,
		Password:       password,
		Nickname:       nickname,
		ProfilePicture: profile_picture,
	}, nil
}

func (ds DAOservice) updateNickName(username string, nickname string) error {
	query := "UPDATE user_tab SET nickname=? WHERE username=?"
	_, err := ds.db.DB.Exec(query, nickname, username)
	if err != nil {
		log.Printf("Error %s when updating nickname\n", err)
	}
	return err
}

func (ds DAOservice) updateProfilePicture(username string, profile_picture string) error {
	query := "UPDATE user_tab SET profile_picture=? WHERE username=?"
	_, err := ds.db.DB.Exec(query, profile_picture, username)
	if err != nil {
		log.Printf("Error %s when updating nickname\n", err)
	}
	return err
}
