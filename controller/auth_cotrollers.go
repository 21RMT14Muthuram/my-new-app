package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	Config "github.com/21RMT14Muthuram/my-new-app/database"
	models "github.com/21RMT14Muthuram/my-new-app/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)



func GetUsers(c *gin.Context) {
    var userList []models.User


    rows, err := Config.DB.Query(context.Background(), `
        SELECT id, email, password_hash, is_verified, 
               otp_code, otp_expires_at, verified_at, 
               created_at, updated_at 
        FROM usermgmt.users
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        fmt.Printf("Query error: %v\n", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User
        err := rows.Scan(
            &user.ID,
            &user.Usermail, 
            &user.Password,
            &user.IsVerified,
            &user.OTPCode,
            &user.OTPExpiresAt,
            &user.VerifiedAt,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user data"})
            fmt.Printf("Scan error: %v\n", err)
            return
        }
        userList = append(userList, user)
    }

    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating through users"})
        fmt.Printf("Rows error: %v\n", err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"users": userList})
}


// VerifyOTPHandler handles OTP verification
func VerifyOTPHandler(c *gin.Context) {
	var verificationRequest struct {
		Email string `json:"email" binding:"required"`
		OTP   string `json:"otp" binding:"required"`
	}

	if err := c.ShouldBindJSON(&verificationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and OTP are required"})
		return
	}

	// Get user's OTP data from database
	var storedOTP string
	var otpExpiresAt time.Time
	var userID int
	var isVerified bool

	err := Config.DB.QueryRow(context.Background(),
		`SELECT id, otp_code, otp_expires_at, is_verified FROM usermgmt.users WHERE email = $1`,
		verificationRequest.Email,
	).Scan(&userID, &storedOTP, &otpExpiresAt, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			fmt.Printf("Database error: %v\n", err)
		}
		return
	}

	// Check if already verified
	if isVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account is already verified"})
		return
	}

	// Validate OTP
	if !ValidateOTP(verificationRequest.OTP, storedOTP, otpExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Mark user as verified and clear OTP data
	now := time.Now()
	_, err = Config.DB.Exec(context.Background(),
		`UPDATE usermgmt.users SET is_verified = TRUE, verified_at = $1, otp_code = NULL, otp_expires_at = NULL WHERE id = $2`,
		now, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify account"})
		fmt.Printf("Update error: %v\n", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

// ResendOTPHandler handles OTP resend requests
func ResendOTPHandler(c *gin.Context) {
	var resendRequest struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&resendRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	// Check if user exists and is not verified
	var userID int
	var isVerified bool
	err := Config.DB.QueryRow(context.Background(),
		`SELECT id, is_verified FROM usermgmt.users WHERE email = $1`,
		resendRequest.Email,
	).Scan(&userID, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if isVerified {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account is already verified"})
		return
	}

	// Generate new OTP
	otp, expiry, err := GenerateOTPWithExpiry()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// Update OTP in database
	_, err = Config.DB.Exec(context.Background(),
		`UPDATE usermgmt.users SET otp_code = $1, otp_expires_at = $2 WHERE id = $3`,
		otp, expiry, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update OTP"})
		fmt.Printf("Update OTP error: %v\n", err)
		return
	}

	// Send new OTP email
	if IsEmailConfigured() {
		err = SendOTPEmail(resendRequest.Email, otp)
		if err != nil {
			fmt.Printf("Failed to send OTP email: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"message": "OTP resent but failed to send email. Please contact support.",
				"otp":     otp, // Include OTP in response for development
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OTP resent to your email"})
	} else {
		// Email not configured
		c.JSON(http.StatusOK, gin.H{
			"message": "OTP regenerated:",
			"otp":     otp, // Remove this in production
		})
	}
}



func SignUpHandler(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Validate input
	if !IsValidEmail(newUser.Usermail) || !IsValidPassword(newUser.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password format"})
		return
	}

	// Check if email already exists (including verified users)
	var existingID int
	var isVerified bool
	err := Config.DB.QueryRow(context.Background(),
		`SELECT id, is_verified FROM usermgmt.users WHERE email = $1`, 
		newUser.Usermail,
	).Scan(&existingID, &isVerified)

	if err == nil {
		if isVerified {
			c.JSON(http.StatusConflict, gin.H{"message": "User email already exists and is verified!"})
		} else {
			// User exists but not verified - allow resending OTP
			c.JSON(http.StatusConflict, gin.H{"message": "Email exists but not verified. Please verify your account."})
		}
		return
	} else if err != pgx.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		fmt.Printf("Database error: %v\n", err)
		return
	}

	// Generate OTP
	otp, expiry, err := GenerateOTPWithExpiry()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate OTP"})
		fmt.Printf("OTP generation error: %v\n", err)
		return
	}

	// Hash password
	hashedPassword, err := Hashing(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		fmt.Printf("Password hashing error: %v\n", err)
		return
	}

	// Insert user with OTP data
	_, err = Config.DB.Exec(context.Background(),
		`INSERT INTO usermgmt.users (email, password_hash, is_verified, otp_code, otp_expires_at) 
		 VALUES ($1, $2, $3, $4, $5)`,
		newUser.Usermail, hashedPassword, false, otp, expiry,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		fmt.Printf("Insert error: %v\n", err)
		return
	}

	// Send OTP email
	if IsEmailConfigured() {
		err = SendOTPEmail(newUser.Usermail, otp)
		if err != nil {
			// Log email error but don't fail the signup
			fmt.Printf("Failed to send OTP email: %v\n", err)
			c.JSON(http.StatusOK, gin.H{
				"message": "User registered successfully but failed to send OTP email. Please contact support.",
				"otp":     otp, // Include OTP in response for development
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully. OTP sent to your email.",
		})
	} else {
		// Email not configured - return OTP in response (for development)
		c.JSON(http.StatusOK, gin.H{
			"message": "User registered successfully. OTP for verification:",
			"otp":     otp, // Remove this in production
		})
	}
}

func LoginHandler(c *gin.Context) {
	var lguser models.User
	if err := c.ShouldBindJSON(&lguser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request"})
		return
	}

	if lguser.Usermail == "" || lguser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
		return
	}

	var checkpass string
	var isVerified bool
	err := Config.DB.QueryRow(context.Background(),
		`SELECT password_hash, is_verified FROM usermgmt.users WHERE email = $1`, 
		lguser.Usermail,
	).Scan(&checkpass, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			fmt.Printf("Database error: %v\n", err)
		}
		return
	}

	// Check if account is verified
	if !isVerified {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Account not verified. Please verify your email first.",
			"verified": false,
		})
		return
	}

	if CheckHashPass(lguser.Password, checkpass) {
		// c.JSON(http.StatusOK, gin.H{"message": "Login Successfully"})

		expirationTime := time.Now().Add(2 * time.Hour)
		claims := &models.Claims{
			Usermail: lguser.Usermail,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(models.JWTKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		if _, err := Config.DB.Exec(context.Background(), `UPDATE usermgmt.users SET jwt_token = $1 where email = $2`, tokenString, lguser.Usermail); err != nil {
			c.JSON(http.StatusLoopDetected, gin.H{"error" : "Unable to store the token"})
			return
		}

		if err := Greeting(lguser.Usermail, "./templates/greeting.html"); err != nil {
			fmt.Printf("Failed to send greeting email: %v\n", err)
		}

		
		c.JSON(http.StatusOK, gin.H{
			"message": "Login Successfully",
			"token": tokenString,
		})
		
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
	}
}

func DeleteUser(c *gin.Context){
	id := c.Param("id")
	result, err := Config.DB.Exec(context.Background(), `DELETE FROM usermgmt.users WHERE id = $1`, id)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message":"DataBase Error"})
		return
	}
	row_Affected := result.RowsAffected()
	if row_Affected == 0{
		c.JSON(http.StatusNotFound, gin.H{"message":"User not Found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully !"})
}
