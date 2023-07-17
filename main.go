package main

import (
	"bufio"
	"freechatgpt/internal/tokens"
	"log"
	"os"
	"strings"
	"time"

	"github.com/acheong08/OpenAIAuth/auth"
	"github.com/acheong08/endless"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
			if len(proxy_parts) > 1 {
				proxies = append(proxies, proxy)
			} else {
				continue
			}
		}
	}
}

func init() {
	_ = godotenv.Load(".env")
	go func() {
		for {
			if os.Getenv("OPENAI_EMAIL") == "" || os.Getenv("OPENAI_PASSWORD") == "" {
				time.Sleep(24 * time.Hour * 7)
				continue
			}
			authenticator := auth.NewAuthenticator(os.Getenv("OPENAI_EMAIL"), os.Getenv("OPENAI_PASSWORD"), os.Getenv("http_proxy"))
			err := authenticator.Begin()
			if err != nil {
				log.Println(err)
				break
			}
			puid, err := authenticator.GetPUID()
			if err != nil {
				break
			}
			os.Setenv("PUID", puid)
			println(puid)
			time.Sleep(24 * time.Hour * 7)
		}
	}()

	HOST = os.Getenv("SERVER_HOST")
	PORT = os.Getenv("SERVER_PORT")
	if PORT == "" {
		PORT = os.Getenv("PORT")
	}
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	if PORT == "" {
		PORT = "8080"
	}
	checkProxy()
	readAccounts()
	scheduleToken()
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
	admin_routes.PATCH("/tokens", tokensHandler)
	admin_routes.PATCH("/puid", puidHandler)
	admin_routes.PATCH("/openai", openaiHandler)
	/// Public routes
	router.OPTIONS("/v1/chat/completions", optionsHandler)
	router.POST("/v1/chat/completions", Authorization, nightmare)
	router.GET("/v1/engines", Authorization, engines_handler)
	router.GET("/v1/models", Authorization, engines_handler)
	endless.ListenAndServe(HOST+":"+PORT, router)
}
