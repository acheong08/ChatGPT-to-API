package main

import (
	"bufio"
	"encoding/json"
	"freechatgpt/internal/tokens"
	"os"
	"strings"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var HOST string
var PORT string
var ACCESS_TOKENS tokens.AccessToken
var proxies []string

func checkProxy() {
	// Check for proxies.txt
	proxies = []string{}
	if _, err := os.Stat("proxies.txt"); err == nil {
		// Each line is a proxy, put in proxies array
		file, _ := os.Open("proxies.txt")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Split line by :
			proxy := scanner.Text()
			proxy_parts := strings.Split(proxy, ":")
			if len(proxy_parts) > 0 {
				proxies = append(proxies, proxy)
			} else {
				continue
			}
		}
	}
}

func init() {
	HOST = os.Getenv("SERVER_HOST")
	PORT = os.Getenv("SERVER_PORT")
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	if PORT == "" {
		PORT = "8080"
	}
	checkProxy()
	// Check if access_tokens.json exists
	if stat, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create("access_tokens.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		updateToken()
	} else {
		nowTime := time.Now()
		usedTime := nowTime.Sub(stat.ModTime())
		// update access token 20 days after last modify token file
		toExpire := 1.728e15 - usedTime
		if toExpire > 0 {
			file, err := os.Open("access_tokens.json")
			if err != nil {
				panic(err)
			}
			defer file.Close()
			decoder := json.NewDecoder(file)
			var token_list []string
			err = decoder.Decode(&token_list)
			if err != nil {
				updateToken()
				return
			}
			if len(token_list) == 0 {
				updateToken()
			} else {
				ACCESS_TOKENS = tokens.NewAccessToken(token_list, false)
				time.AfterFunc(toExpire, updateToken)
			}
		} else {
			updateToken()
		}
	}
}

func main() {
	router := gin.Default()

	router.Use(cors)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	admin_routes := router.Group("/admin")
	admin_routes.Use(adminCheck)

	/// Admin routes
	admin_routes.PATCH("/password", passwordHandler)
	admin_routes.PATCH("/tokens", adminCheck, tokensHandler)
	/// Public routes
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", nightmare)
	endless.ListenAndServe(HOST+":"+PORT, router)
}
