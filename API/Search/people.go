package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"searchquery/initializers"
)

type Username struct {
	Username string `json:"username"`
}

func (app *Config) searchPeople(w http.ResponseWriter, r *http.Request) {

	var userQuery Query
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to query search form: %v", err), http.StatusInternalServerError)
		return
	}

	var usernames []Username
	result := initializers.DB.Table("users").Select("username").Where("username LIKE ?", "%"+userQuery.Query+"%").Scan(&usernames)
	if result.Error != nil {
		errorJSON(w, fmt.Errorf("error: Failed to find Users: %v", result.Error), http.StatusInternalServerError)
		return
	}

	//Response
	responseData := struct {
		Message string     `json:"message"`
		Users   []Username `json:"users"`
	}{}

	if len(usernames) > 0 {
		responseData.Message = "Found Users"
	} else {
		responseData.Message = fmt.Sprintf("No users found matching the search query '%s'.", userQuery.Query)
	}

	responseData.Users = usernames

	// To JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	err = encoder.Encode(responseData)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to encode response to JSON: %v", err), http.StatusInternalServerError)
		return
	}
}
