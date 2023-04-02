package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/acheong08/OpenAIAuth/auth"
)

type Account struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}
type Proxy struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (p Proxy) Socks5URL() string {
	// Returns proxy URL (socks5)
	return fmt.Sprintf("socks5://%s:%s@%s:%s", p.User, p.Pass, p.IP, p.Port)
}

// Read accounts.txt and create a list of accounts
func readAccounts() []Account {
	accounts := []Account{}
	// Read accounts.txt and create a list of accounts
	file, err := os.Open("accounts.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Loop through each line in the file
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
	return accounts
}

// Read proxies from proxies.txt and create a list of proxies
func readProxies() []Proxy {
	proxies := []Proxy{}
	// Read proxies.txt and create a list of proxies
	file, err := os.Open("proxies.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Loop through each line in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Split by :
		line := strings.Split(scanner.Text(), ":")
		// Create a proxy
		proxy := Proxy{
			IP:   line[0],
			Port: line[1],
			User: line[2],
			Pass: line[3],
		}
		// Append to proxies
		proxies = append(proxies, proxy)
	}
	return proxies
}

func main() {
	// Read accounts and proxies
	accounts := readAccounts()
	proxies := readProxies()
	puid := os.Getenv("PUID")

	// Loop through each account
	for _, account := range accounts {
		println(proxies[0].Socks5URL())
		println(account.Email)
		println(account.Password)
		authenticator := auth.NewAuthenticator(account.Email, puid, account.Password, proxies[0].Socks5URL())
		// Push used proxy to the back of the list
		proxies = append(proxies[1:], proxies[0])
		err := authenticator.Begin()
		if err.Error != nil {
			// println("Error: " + err.Details)
			println("Location: " + err.Location)
			println("Status code: " + fmt.Sprint(err.StatusCode))
			println("Embedded error: " + err.Error.Error())
			return
		}
		access_token, err := authenticator.GetAccessToken()
		if err.Error != nil {
			// println("Error: " + err.Details)
			println("Location: " + err.Location)
			println("Status code: " + fmt.Sprint(err.StatusCode))
			println("Embedded error: " + err.Error.Error())
			return
		}
		// Append access token to access_tokens.txt
		f, go_err := os.OpenFile("access_tokens.txt", os.O_APPEND|os.O_WRONLY, 0600)
		if go_err != nil {
			continue
		}
		defer f.Close()
		if _, go_err = f.WriteString(access_token + "\n"); go_err != nil {
			continue
		}
		// Write authenticated account to authenticated_accounts.txt
		f, go_err = os.OpenFile("authenticated_accounts.txt", os.O_APPEND|os.O_WRONLY, 0600)
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
}
