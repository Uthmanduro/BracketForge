package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(password, hashedPassword string) bool {
	// compare password with hashed password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}