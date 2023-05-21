package main

import (
	"context"
	"encoding/json"
	"freechatgpt/internal/tokens"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var HOST string
var PORT string
var ACCESS_TOKENS tokens.AccessToken

func init() {
	HOST = os.Getenv("SERVER_HOST")
	PORT = os.Getenv("SERVER_PORT")
	if HOST == "" {
		HOST = "127.0.0.1"
	}
	if PORT == "" {
		PORT = "8080"
	}
	accessToken := os.Getenv("ACCESS_TOKENS")
	if accessToken != "" {
		accessTokens := strings.Split(accessToken, ",")
		ACCESS_TOKENS = tokens.NewAccessToken(accessTokens)
	}
	// Check if access_tokens.json exists
	if _, err := os.Stat("access_tokens.json"); os.IsNotExist(err) {
		// Create the file
		file, err := os.Create("access_tokens.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		// Load the tokens
		file, err := os.Open("access_tokens.json")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		var token_list []string
		err = decoder.Decode(&token_list)
		if err != nil {
			return
		}
		ACCESS_TOKENS = tokens.NewAccessToken(token_list)
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
	srv := &http.Server{
		Addr:    HOST + ":" + PORT,
		Handler: router,
	}

	// Receive another goroutine listen error
	serverError := make(chan error, 1)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		serverError <- srv.ListenAndServe()
	}()

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("listen: %s\n", err)
	case <-quit:
		log.Println("Shutting down server...")
		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		select {
		case <-ctx.Done():
			log.Println("timeout of 5 seconds.")
		}
		log.Println("Server exiting")
	}
}
