package controller

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)


func Hashing(pass string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	fmt.Print(hash)
	return string(hash), err
}


func CheckHashPass(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}




