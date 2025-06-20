package utils

import (
	"backend/config"
	"crypto/rand"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars     = "0123456789"
	specialChars   = "!@#$%^&*"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.AppConfig.BcryptRounds)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateSecurePassword generates a cryptographically secure random password
func GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 8 // Minimum secure password length
	}

	// Ensure password contains at least one character from each character set
	allChars := lowercaseChars + uppercaseChars + digitChars + specialChars
	password := make([]byte, length)

	// Add at least one character from each required set
	charSets := []string{lowercaseChars, uppercaseChars, digitChars, specialChars}
	for i := 0; i < 4 && i < length; i++ {
		char, err := getRandomChar(charSets[i])
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Fill the rest with random characters from all sets
	for i := 4; i < length; i++ {
		char, err := getRandomChar(allChars)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Shuffle the password to avoid predictable patterns
	for i := range password {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(len(password))))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

// getRandomChar returns a random character from the given charset
func getRandomChar(charset string) (byte, error) {
	randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
	if err != nil {
		return 0, err
	}
	return charset[randomIndex.Int64()], nil
}
