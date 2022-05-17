package services

import (
	"encoding/base64"
	"errors"

	//"fmt"
	"io/ioutil"
	"log"
	"os"
)

type ImageInterface interface {
	StoreImage(profilePictureName string, data string, username string) (string, error)
}

type ImageService struct {
}

func InitializeImageService() ImageInterface {
	return ImageService{}
}

func (is ImageService) StoreImage(profilePictureName string, data string, username string) (string, error) {

	decodedImage, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Printf("Error decoding image: %s\n", err)
		return "", err
	}
	imageDirectory := "../resources/" + username
	if _, err = os.Stat(imageDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(imageDirectory, 0777)
		if err != nil {
			log.Printf("Error creating image directory: %s\n", err)
			return "", err
		}
	}
	imagePath := imageDirectory + "/" + profilePictureName
	err = ioutil.WriteFile(imagePath, decodedImage, 0777)
	if err != nil {
		log.Printf("Error writing image file: %s\n", err)
		return "", err
	}

	dbRelativePath := "resources/" + username + "/" + profilePictureName
	return dbRelativePath, nil
}
