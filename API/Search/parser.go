package main

import "encoding/json"

// BASE
// 	hits
// 	  hits
// 		_source
//			screen_name
//			text

func ParseTweets(jsonData []byte) ([]byte, error) {
	var data map[string]interface{} //take into base
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	hitsData := data["hits"].(map[string]interface{}) // first hits
	hits := hitsData["hits"].([]interface{})          // inner hits

	var results []map[string]string
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{}) // _source
		screenName := source["screen_name"].(string)                               // screen_name
		text := source["text"].(string)                                            //text
		result := map[string]string{
			"screen_name": screenName,
			"text":        text,
		}
		results = append(results, result)
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	return resultJSON, nil
}
