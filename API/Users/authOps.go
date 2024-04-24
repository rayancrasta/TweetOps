package main

import (
	"Users/initializers"
	"Users/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

type Profile struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Followerscount int    `json:"followers_count"`
	Followingcount int    `json:"following_count"`
	Verified       bool   `json:"verified"`
	Accountlang    string `json:"account_lang"`
}

func (app *Config) SignUp(w http.ResponseWriter, r *http.Request) {

	var profile Profile
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse signup form: %v", err), http.StatusInternalServerError)
		return
	}

	//hashing the password
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), 10)

	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to hash password: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a user
	user := models.User{UserID: profile.UserID, Username: profile.Username, Password: string(hash), FollowersCount: profile.Followerscount, FollowingCount: profile.Followingcount, Accountlang: profile.Accountlang}
	result := initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&user) // if not exsists insert

	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to create User: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var response jsonResponse
	response.Message = fmt.Sprintf("Success: Created User: %v", profile.Username)
	writeJSON(w, http.StatusCreated, response)
}

func (app *Config) Login(w http.ResponseWriter, r *http.Request) {
	var loginbody Profile
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&loginbody)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse signup form: %v", err), http.StatusInternalServerError)
		return
	}

	var user models.User
	initializers.DB.First(&user, "username = ?", loginbody.Username)

	// log.Println(user.UserID)
	if user.UserID == 0 {
		errorJSON(w, fmt.Errorf("error: Invalid User or Password"), http.StatusInternalServerError)
		return
	}

	//Compare the password after hashing with password in db
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginbody.Password))

	if err != nil {
		errorJSON(w, fmt.Errorf("error: Invalid User or Password"), http.StatusInternalServerError)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": int(user.UserID),
		"exp": time.Now().Add(time.Hour * 5).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to create Token"), http.StatusInternalServerError)
		return
	}

	var response jsonResponse
	response.Message = fmt.Sprintf("Login Succesful")
	response.Token = tokenString
	writeJSON(w, http.StatusCreated, response)
}
