package services

import (
	"errors"
	"fmt"
	"log"
)

type AuthenticationInterface interface {
	AuthenticateLogin(username string, password string) error
}

type AuthenticationService struct {
	ds DAOInterface
	hs HashService
}

func InitializeAuthService(dao DAOInterface) AuthenticationInterface {
	return AuthenticationService{ds: dao}
}

func (as AuthenticationService) AuthenticateLogin(username string, password string) error {
	user, err := as.ds.fetchUser(username)
	if err != nil {
		log.Printf("User does not exist: %s\n", err)
		// return &pb.LoginResponse{
		// 	Rc: pb.LoginResponse_Failed,
		// }, err
		return err
	}

	if as.hs.ComparePasswords([]byte(user.GetPassword()), password) {
		fmt.Println("Correct password")
		// return &pb.LoginResponse{
		// 	Rc:     pb.LoginResponse_Success,
		// 	//UserId: user.UserId,
		// }, nil
		return nil
	}

	fmt.Println("Wrong password")
	// return &pb.LoginResponse{
	// 	Rc: pb.LoginResponse_Failed,
	// }, err
	return errors.New("incorrect password") //need to modify this error here
}
