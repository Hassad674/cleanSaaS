package hash

import "golang.org/x/crypto/bcrypt"

// cost is the bcrypt work factor. 12 is the recommended minimum for production.
// DefaultCost (10) is too weak for modern hardware.
const cost = 12

func Password(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func Check(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
