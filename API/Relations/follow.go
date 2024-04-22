package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"relations/initializers"
	"relations/models"

	"gorm.io/gorm/clause"
)

type FollowRequest struct {
	UserID   int `json:"user_id"`
	TargetID int `json:"target_id"` // token user follows this ID
}

func (app *Config) Follow(w http.ResponseWriter, r *http.Request) {

	var followreq FollowRequest
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&followreq)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse followed form: %v", err), http.StatusInternalServerError)
		return
	}

	if followreq.TargetID == followreq.UserID {
		errorJSON(w, fmt.Errorf("UserID and TargetID cant be same"), http.StatusInternalServerError)
		return
	}
	// Make an entry in Followers table // targetID is being followed by UserID
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

	var followreq FollowRequest
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&followreq)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to parse followed form: %v", err), http.StatusInternalServerError)
		return
	}
	if followreq.TargetID == followreq.UserID {
		errorJSON(w, fmt.Errorf("UserID and TargetID cant be same"), http.StatusInternalServerError)
		return
	}
	// Remove an entry from Followers table
	// Make an entry in Followers table
	follower := models.Follower{UserID: followreq.TargetID, FollowerID: followreq.UserID}
	result := initializers.DB.Where(&follower).Delete(&models.Follower{})
	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to delete in Followers TB: %v", result.Error), http.StatusInternalServerError)
		return
	}

	// Remove an entry in Following table
	// Make an entry in Following table //Userid follows TargetID
	following := models.Following{UserID: followreq.UserID, FollowingID: followreq.TargetID}
	result = initializers.DB.Where(&following).Delete(&models.Following{})

	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to delete  in Following TB: %v", result.Error), http.StatusInternalServerError)
		return
	}

	var response jsonResponse
	response.Message = fmt.Sprintf("Success: %v UnFollowed %v", followreq.UserID, followreq.TargetID)
	writeJSON(w, http.StatusOK, response)
}

