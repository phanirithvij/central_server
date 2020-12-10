package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72

// Hash returns the password hash
func Hash(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

// ComparePasswords compare the hash and password
func ComparePasswords(hashedPwd string, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
