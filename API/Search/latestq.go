package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
)

type Query struct {
	Query    string `json:"query"`
	Safemode bool   `json:"safemode",optional`
}

func (app *Config) searchLatest(w http.ResponseWriter, r *http.Request) {

	var userQuery Query
	//Decode the payload
	err := json.NewDecoder(r.Body).Decode(&userQuery)
	if err != nil {
		errorJSON(w, fmt.Errorf("error: Failed to query search form: %v", err), http.StatusInternalServerError)
		return
	}

	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Elasticsearch server address
	}

	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Printf("Error creating the Elasticsearch client: %s", err)
		return
	}

	// Search query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"text": userQuery.Query,
			},
		},
		"sort": []map[string]interface{}{
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Printf("Error encoding query: %s", err)
		return
	}

	// Search request
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("tweetscombined"), // Index to search within
		esClient.Search.WithBody(&buf),
		esClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		fmt.Printf("Error performing search request: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("Error response: %s", res.Status())
		return
	}

	// Parse the response
	var resp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		fmt.Printf("Error parsing the response body: %s", err)
		return
	}

	// Marshal response into readable JSON
	jsonResponse, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling response: %s", err)
		return
	}

	// Return the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
