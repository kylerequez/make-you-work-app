package utils

import (
	"errors"
	"math/rand"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const numberOfChar = 10

func GenerateNewPassword() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{};':,.<>/?")

	s := make([]rune, numberOfChar)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func ValidatePassword(password, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func ComparePassword(password, confirmPassword string) error {
	isPasswordSame := strings.Compare(password, confirmPassword)
	switch isPasswordSame {
	case 0:
		return nil
	default:
		return errors.New("passwords are not the same")
	}
}
