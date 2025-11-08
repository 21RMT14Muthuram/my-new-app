package controller

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
    // Basic email regex pattern
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    match, _ := regexp.MatchString(emailRegex, email)
    return match
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
    // Check minimum length
    if len(password) < 6 {
        return false
    }
	
    // Check for at least one special character
    specialCharRegex := `[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`
    hasSpecialChar, _ := regexp.MatchString(specialCharRegex, password)
    
    // Check for at least one numeric character
    numericRegex := `[0-9]`
    hasNumeric, _ := regexp.MatchString(numericRegex, password)
    
    return hasSpecialChar && hasNumeric
}

func Hashing(pass string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	fmt.Print(hash)
	return string(hash), err
}


func CheckHashPass(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}




