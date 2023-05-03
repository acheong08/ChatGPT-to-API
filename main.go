package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _  = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	user_agent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"
	http_proxy = os.Getenv("http_proxy")
)

func main() {

	if http_proxy != "" {
		client.SetProxy(http_proxy)
		println("Proxy set:" + http_proxy)
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "9090"
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	handler.OPTIONS("/v1/chat/completions", optionsHandler)

	handler.POST("/v1/chat/completions", proxy)

	gin.SetMode(gin.ReleaseMode)
	endless.ListenAndServe(os.Getenv("HOST")+":"+PORT, handler)
}

func proxy(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	// Read body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Parse as json
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if "max_tokens" is set and if set to nil
	if _, ok := jsonBody["max_tokens"]; !ok {

	} else if jsonBody["max_tokens"] == nil {
		// Remove "max_tokens" from json
		delete(jsonBody, "max_tokens")
	}

	request_body, err := json.Marshal(jsonBody)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if _, ok := jsonBody["model"]; !ok {
		c.JSON(400, gin.H{"error": "No model specified"})
		return
	}

	var url string
	var request_method string
	var request *http.Request
	var response *http.Response

	url = "https://api.jeeves.ai/generate/v3/chat"
	request_method = c.Request.Method

	request, err = http.NewRequest(request_method, url, bytes.NewReader(request_body))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	request.Header.Set("Host", "api.jeeves.ai")
	request.Header.Set("Origin", "https://jeeves.ai/")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Keep-Alive", "timeout=360")
	request.Header.Set("Authorization", c.Request.Header.Get("Authorization"))
	request.Header.Set("sec-ch-ua", "\"Chromium\";v=\"112\", \"Brave\";v=\"112\", \"Not:A-Brand\";v=\"99\"")
	request.Header.Set("sec-ch-ua-mobile", "?0")
	request.Header.Set("sec-ch-ua-platform", "\"Linux\"")
	request.Header.Set("sec-fetch-dest", "empty")
	request.Header.Set("sec-fetch-mode", "cors")
	request.Header.Set("sec-fetch-site", "same-origin")
	request.Header.Set("sec-gpc", "1")
	request.Header.Set("user-agent", user_agent)

	response, err = client.Do(request)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()
	// Check if "stream" is set and if set to true
	if _, ok := jsonBody["stream"]; !ok {
		if jsonBody["stream"] == true {
			c.Header("Content-Type", response.Header.Get("Content-Type"))
			// Get status code
			c.Status(response.StatusCode)
			c.Stream(func(w io.Writer) bool {
				// Write data to client
				io.Copy(w, response.Body)
				return false
			})
			return
		}
	}
	// Loop through response
	if response.StatusCode != 200 {
		c.JSON(response.StatusCode, gin.H{"error": "Error"})
		return
	}
	var fulltext string = ""
	for {
		// Stream each line
		line, err := bufio.NewReader(response.Body).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if len(line) < 5 {
			fmt.Println(line)
			break
		}
		// Remove data:
		line = strings.Replace(line, "data: ", "", 1)
		if strings.HasPrefix(line, "[DONE]") {
			break
		}
		if !strings.HasPrefix(line, "{") {
			break
		}
		// Parse as json
		var jsonLine Data
		err = json.Unmarshal([]byte(line), &jsonLine)
		if err != nil {
			fmt.Println(err)
			break
		}
		fulltext += fulltext + jsonLine.Choices[0].Delta.Content
	}
	c.JSON(200, NewFullCompletion(fulltext, jsonBody["model"].(string)))

}

func optionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
