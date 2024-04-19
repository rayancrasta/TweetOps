package main

import (
	"Users/initializers"
	"Users/models"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm/clause"
)

type FollowRequest struct {
	UserID   int `json:"user_id"`
	TargetID int `json:"followed_id"` // token user follows this ID
}

func (app *Config) Follow(w http.ResponseWriter, r *http.Request) {
	
	var followreq FollowRequest
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&followreq)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse followed form: %v", err), http.StatusInternalServerError)
		return
	}

	// Make an entry in Followers table
	follower := models.Follower{UserID: followreq.TargetID, FollowerID: followreq.UserID}
	result := initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&follower)

	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to create entry in Followers TB: %v", result.Error), http.StatusInternalServerError)
		return
	}
	// Make an entry in Following table //Userid follows TargetID
	following := models.Following{UserID: followreq.UserID, FollowingID: followreq.TargetID}
	result = initializers.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&following)

	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to create entry in Following TB: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var response jsonResponse
	response.Message = fmt.Sprintf("Success: %v Follows %v", followreq.UserID, followreq.TargetID)
	writeJSON(w, http.StatusOK, response)

}

func (app *Config) UnFollow(w http.ResponseWriter, r *http.Request) {

	var followed FollowRequest
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&followed)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse followed form: %v", err), http.StatusInternalServerError)
		return
	}

	// Remove an entry from Followers table

	// Remove an entry in Following table

}
