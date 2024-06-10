package utils

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolBytes = "!@#$%^&*()-_=+[{]};:'\",<.>/? "
	numberBytes = "0123456789"
	allCharSet  = letterBytes + symbolBytes + numberBytes
)

func GenerateNewPassword() ([]byte, error) {
	numberOfChar := 10
	charSet := allCharSet

	b := make([]byte, numberOfChar)
	var err error

	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charSet[rand.Intn(len(charSet))]
	}

	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 8)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func ValidatePassword(password []byte, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func ComparePassword(password string, confirmPassword string) error {
	isPasswordSame := strings.Compare(password, confirmPassword)
	switch isPasswordSame {
	case 0:
		return nil
	default:
		return errors.New("passwords are not the same")
	}
}
