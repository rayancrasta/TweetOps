package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

type Query struct {
	Query       string `json:"query"`
	Safemode    bool   `json:"safemode",optional`
	Page        int    `json:"page",optional`
	Accountlang string `json:"account_lang"` // middleware
}

func (app *Config) searchLatest(w http.ResponseWriter, r *http.Request) {

	var userQuery Query
	userQuery.Page = 0 // default for pagination

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
	query := getQuery(userQuery) // if safe mode true, then check for safe flag
	print(query)

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Printf("Error encoding query: %s", err)
		return
	}

	// Search request
	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex("tweetscombined"), // Index to search within
		esClient.Search.WithBody(&buf),
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

	// Return the JSON response
	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	resultJSON, err := ParseTweets(bodyContent)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)

}

func getQuery(userQuery Query) map[string]interface{} {
	if userQuery.Safemode { //check safe flag
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"text": userQuery.Query,
							},
						},
						{
							"term": map[string]interface{}{
								"lang": userQuery.Accountlang,
							},
						},
						{
							"term": map[string]interface{}{
								"safe": userQuery.Safemode,
							},
						},
					},
				},
			},
			"sort": []map[string]interface{}{
				{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			},
			"from": userQuery.Page,
			"size": os.Getenv("resultsize"),
			"_source": map[string]interface{}{
				"includes": []string{"screen_name", "text"},
			},
		}
		return query
	} else { // normal
		query := map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"match": map[string]interface{}{
								"text": userQuery.Query,
							},
						},
						{
							"term": map[string]interface{}{
								"lang": userQuery.Accountlang,
							},
						},
					},
				},
			},
			"sort": []map[string]interface{}{
				{
					"created_at": map[string]interface{}{
						"order": "desc",
					},
				},
			},
			"from": userQuery.Page,
			"size": os.Getenv("resultsize"),
			"_source": map[string]interface{}{
				"includes": []string{"screen_name", "text"},
			},
		}
		return query
	}
}
