package util

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(in string) (out string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	if err == nil {
		out = string(hash)
	}
	return
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
