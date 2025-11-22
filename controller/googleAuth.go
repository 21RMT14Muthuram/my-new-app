package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig *oauth2.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}


// func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
// 	url := googleOauthConfig.AuthCodeURL("randomstate")
// 	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// }
func HandleGoogleLogin(c *gin.Context){
	url := googleOauthConfig.AuthCodeURL("randomstate")
	c.Redirect(http.StatusTemporaryRedirect, url)
}


// func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
// 	if r.FormValue("state") != "randomstate" {
// 		http.Error(w, "State mismatch", http.StatusBadRequest)
// 		return
// 	}

// 	token, err := googleOauthConfig.Exchange(context.Background(), r.FormValue("code"))
// 	if err != nil {
// 		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
// 	if err != nil {
// 		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	userInfo := map[string]interface{}{}
// 	json.NewDecoder(resp.Body).Decode(&userInfo)

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(userInfo)
// }

func HandleGoogleCallback(c *gin.Context) {

	// 1. Validate state to prevent CSRF
	if c.Query("state") != "randomstate" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "State mismatch"})
		return
	}

	// 2. Get the "code" from Google
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Code not found"})
		return
	}

	// 3. Exchange the code for an access token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to exchange token",
			"error":   err.Error(),
		})
		return
	}

	// 4. Get user info from Google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// 5. Decode user info JSON
	userInfo := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user info"})
		return
	}

	// 6. Return user info to client
	c.JSON(http.StatusOK, gin.H{
		"message":  "Google login successful",
		"userInfo": userInfo,
	})
}
