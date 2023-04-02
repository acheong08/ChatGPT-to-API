package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Account struct {
	Username string `json:"username"`
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
	return fmt.Sprintf("socks5h://%s:%s@%s:%s", p.User, p.Pass, p.IP, p.Port)
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
			Username: line[0],
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
	// Print accounts and proxies as test
	fmt.Println(accounts)
	fmt.Println(proxies)
}
