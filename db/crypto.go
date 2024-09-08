package db

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordNotMachRequirements = errors.New("crypto: Password dose not match the requirements.")

// hashPassword hashes a plain text password with bcrypt and returns the hashed password.
func hashPassword(password string) (string, error) {
	// Generate a hash of the password using bcrypt
	log.Println(len(password), len([]byte(password)))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func verifyPassword(password string, hashedPassword string) (bool, error) {
	// Generate a hash of the password using bcrypt

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil
}
