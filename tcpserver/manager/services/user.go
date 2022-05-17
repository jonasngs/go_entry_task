package services

import (
	//"fmt"
	"log"
)

type UserInterface interface {
	UpdateNickname(username string, sessionToken string, nickname string) error
	UpdateProfilePicture(username string, sessionToken string, data string, file_ext string) error
}

type UserService struct {
	ds DAOInterface
	cs CacheInterface
	is ImageInterface
}

func InitializeUserService(dao DAOInterface, cs CacheInterface, is ImageInterface) UserInterface {
	return UserService{ds: dao, cs: cs, is: is}
}

func (us UserService) UpdateProfilePicture(username string, sessionToken string, data string, file_ext string) error {
	
	if data == "" && file_ext == "" {
		log.Printf("No image provided by user")
		return nil
	}
	
	profilePictureName := username + "-profile-picture" + file_ext
	imagePath, err := us.is.StoreImage(profilePictureName, data, username)

	if err != nil {
		log.Printf("Error %s when storing profile picture: ", err)
		return err
	}
	// Update database with image path
	err = us.ds.updateProfilePicture(username, imagePath)
	if err != nil {
		log.Printf("Error %s when updating profile picture in database: ", err)
		return err
	}

	// Update cache with image path
	user, err := us.cs.checkCache(sessionToken)
	if err != nil {
		log.Printf("User not found in cache: %s\n", err)
		return err
	}

	updatedUser := user
	updatedUser.ProfilePicture = imagePath
	err = us.cs.updateCache(sessionToken, updatedUser)
	if err != nil {
		log.Printf("Error %s when updating profile picture in cache: ", err)
		return err
	}

	return nil
}

func (us UserService) UpdateNickname(username string, sessionToken string, nickname string) error {
	// Update database
	err := us.ds.updateNickName(username, nickname)
	if err != nil {
		log.Printf("Error %s when updating nickname in database: ", err)
		return err
	}

	// Update Cache
	user, err := us.cs.checkCache(sessionToken)
	if err != nil {
		log.Printf("User not found in cache: %s\n", err)
		return err
	}

	updatedUser := user
	updatedUser.Nickname = nickname

	err = us.cs.updateCache(sessionToken, updatedUser)
	if err != nil {
		log.Printf("Error %s when updating nickname in cache: ", err)
		return err
	}

	return nil
}
