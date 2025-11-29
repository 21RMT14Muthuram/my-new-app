package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
    ID           int        `json:"id"`
    Usermail     string     `json:"email"`
    Password     string     `json:"password"`
    IsVerified   bool       `json:"is_verified"`
    OTPCode      *string    `json:"otp_code,omitempty"`
    OTPExpiresAt *time.Time `json:"otp_expires_at,omitempty"`
    VerifiedAt   *time.Time `json:"verified_at,omitempty"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

type Claims struct {
	Usermail string `json:"usermail"`
	jwt.RegisteredClaims
}

var JWTKey = []byte("my_secret_key")


type OrganizationType struct{
    ID          int         `json:"id"`
    OrgCode     string      `json:"orgcode"`
    OrgName     string      `json:"orgname"`
}