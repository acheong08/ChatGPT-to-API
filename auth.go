package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"freechatgpt/internal/tokens"

	"github.com/acheong08/OpenAIAuth/auth"
)

var accounts []Account

type Account struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

// Read accounts.txt and create a list of accounts
func readAccounts() {
	accounts = []Account{}
	// Read accounts.txt and create a list of accounts
	if _, err := os.Stat("accounts.txt"); err == nil {
		// Each line is a proxy, put in proxies array
		file, _ := os.Open("accounts.txt")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Split by :
			line := strings.Split(scanner.Text(), ":")
			// Create an account
			account := Account{
				Email:    line[0],
				Password: line[1],
			}
			// Append to accounts
			accounts = append(accounts, account)
		}
	}
}
func scheduleToken() {
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

func updateToken() {
	token_list := []string{}
	// Loop through each account
	for _, account := range accounts {
		if os.Getenv("CF_PROXY") != "" {
			// exec warp-cli disconnect and connect
			exec.Command("warp-cli", "disconnect").Run()
			exec.Command("warp-cli", "connect").Run()
			time.Sleep(5 * time.Second)
		}
		println("Updating access token for " + account.Email)
		var proxy_url string
		if len(proxies) == 0 {
			proxy_url = ""
		} else {
			proxy_url = proxies[0]
			// Push used proxy to the back of the list
			proxies = append(proxies[1:], proxies[0])
		}
		authenticator := auth.NewAuthenticator(account.Email, account.Password, proxy_url)
		err := authenticator.Begin()
		if err != nil {
			// println("Error: " + err.Details)
			println("Location: " + err.Location)
			println("Status code: " + fmt.Sprint(err.StatusCode))
			println("Details: " + err.Details)
			println("Embedded error: " + err.Error.Error())
			return
		}
		access_token := authenticator.GetAccessToken()
		token_list = append(token_list, access_token)
		println("Success!")
		// Write authenticated account to authenticated_accounts.txt
		f, go_err := os.OpenFile("authenticated_accounts.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if go_err != nil {
			continue
		}
		defer f.Close()
		if _, go_err = f.WriteString(account.Email + ":" + account.Password + "\n"); go_err != nil {
			continue
		}
		// Remove accounts.txt
		os.Remove("accounts.txt")
		// Create accounts.txt
		f, go_err = os.Create("accounts.txt")
		if go_err != nil {
			continue
		}
		defer f.Close()
		// Remove account from accounts
		accounts = accounts[1:]
		// Write unauthenticated accounts to accounts.txt
		for _, acc := range accounts {
			// Check if account is authenticated
			if acc.Email == account.Email {
				continue
			}
			if _, go_err = f.WriteString(acc.Email + ":" + acc.Password + "\n"); go_err != nil {
				continue
			}
		}
	}
	// Append access token to access_tokens.json
	ACCESS_TOKENS = tokens.NewAccessToken(token_list, true)
	time.AfterFunc(1.728e15, updateToken)
}
