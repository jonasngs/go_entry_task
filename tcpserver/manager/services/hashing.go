package services

import "golang.org/x/crypto/bcrypt"

type HashingInterface interface {
	Hash(password string) ([]byte, error)
	ComparePasswords(hash []byte, plaintext string) bool
}

type HashService struct {
}

func CreateHashFunction() HashingInterface {
	return HashService{}
}

func (hs HashService) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (hs HashService) ComparePasswords(hash []byte, plaintext string) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
	return err == nil
}
