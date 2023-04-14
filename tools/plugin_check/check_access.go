package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	var access_tokens []string
	// Read access_tokens.txt and split by new line
	file, err := os.Open("access_tokens.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		access_tokens = append(access_tokens, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	// Go routine to check access for each token (limit to 20 simultaneous)
	sem := make(chan bool, 20)
	for _, token := range access_tokens {
		sem <- true
		go func(token string) {
			defer func() { <-sem }()
			if check_access(token) {
				println(token)
			}
		}(token)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func check_access(token string) bool {
	print(".")
	req, _ := http.NewRequest("GET", "https://chat.openai.com/backend-api/accounts/check", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	// Set _puid cookie
	req.AddCookie(&http.Cookie{Name: "_puid", Value: os.Getenv("PUID")})
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		// Parse response body as JSON
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		// Check if "tool1", "tool2", or "tool3" is in the features array
		for _, feature := range result["features"].([]interface{}) {
			if feature == "tool1" || feature == "tool2" || feature == "tool3" {
				return true
			}
		}
		return false
	}
	println(resp.StatusCode)
	return false
}
