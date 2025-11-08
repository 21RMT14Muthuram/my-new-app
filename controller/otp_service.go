package controller

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// OTPConfig holds OTP configuration
type OTPConfig struct {
	Length          int
	ExpiryMinutes   int
	MaxAttempts     int
}

// Default OTP configuration
var DefaultOTPConfig = OTPConfig{
	Length:        6,
	ExpiryMinutes: 10,
	MaxAttempts:   3,
}

// GenerateOTP generates a random numeric OTP
func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		length = DefaultOTPConfig.Length
	}

	otp := ""
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %v", err)
		}
		otp += fmt.Sprintf("%d", num)
	}

	return otp, nil
}

// GenerateOTPWithExpiry generates OTP and its expiry time
func GenerateOTPWithExpiry() (string, time.Time, error) {
	otp, err := GenerateOTP(DefaultOTPConfig.Length)
	if err != nil {
		return "", time.Time{}, err
	}

	expiry := time.Now().Add(time.Duration(DefaultOTPConfig.ExpiryMinutes) * time.Minute)
	return otp, expiry, nil
}

// IsOTPExpired checks if OTP has expired
func IsOTPExpired(expiryTime time.Time) bool {
	return time.Now().After(expiryTime)
}

// ValidateOTP validates the OTP against stored OTP and expiry
func ValidateOTP(inputOTP, storedOTP string, storedExpiry time.Time) bool {
	// Check if OTP has expired
	if IsOTPExpired(storedExpiry) {
		return false
	}

	// Check if OTP matches
	return inputOTP == storedOTP
}