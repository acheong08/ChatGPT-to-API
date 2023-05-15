package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	tls_client "github.com/bogdanfinn/tls-client"
)

var proxies []string

// Read proxies.txt and check if they work
func init() {
	// Check for proxies.txt
	if _, err := os.Stat("proxies.txt"); err == nil {
		// Each line is a proxy, put in proxies array
		file, _ := os.Open("proxies.txt")
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			// Split line by :
			proxy := scanner.Text()
			proxy_parts := strings.Split(proxy, ":")
			if len(proxy_parts) == 2 {
				proxy = "socks5://" + proxy
			} else if len(proxy_parts) == 4 {
				proxy = "socks5://" + proxy_parts[2] + ":" + proxy_parts[3] + "@" + proxy_parts[0] + ":" + proxy_parts[1]
			} else {
				continue
			}
			proxies = append(proxies, proxy)
		}
	}

}

func main() {
	wg := sync.WaitGroup{}
	for _, proxy := range proxies {
		wg.Add(1)
		go func(proxy string) {
			defer wg.Done()
			jar := tls_client.NewCookieJar()
			options := []tls_client.HttpClientOption{
				tls_client.WithTimeoutSeconds(360),
				tls_client.WithClientProfile(tls_client.Chrome_110),
				tls_client.WithNotFollowRedirects(),
				tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
				// Disable SSL verification
				tls_client.WithInsecureSkipVerify(),
			}
			client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

			client.SetProxy(proxy)
			resp, err := client.Get("https://example.com")
			if err != nil {
				fmt.Println("Error: ", err)
				fmt.Println("Proxy: ", proxy)
				return
			}
			if resp.StatusCode != 200 {
				fmt.Println("Error: ", resp.StatusCode)
				fmt.Println("Proxy: ", proxy)
				return
			} else {
				fmt.Println(".")
			}
		}(proxy)
	}
	wg.Wait()
}
