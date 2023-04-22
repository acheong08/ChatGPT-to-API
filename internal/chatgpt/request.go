package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"strings"

	typings "freegpt4/internal/typings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var proxies []string

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		// Disable SSL verification
		tls_client.WithInsecureSkipVerify(),
	}
	client, _  = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	http_proxy = os.Getenv("http_proxy")
)

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

func random_int(min int, max int) int {
	return min + rand.Intn(max-min)
}

func SendRequest(message typings.ChatMessage) (*http.Response, error) {
	if http_proxy != "" && len(proxies) > 0 {
		client.SetProxy(http_proxy)
	}
	// Take random proxy from proxies.txt
	if len(proxies) > 0 {
		client.SetProxy(proxies[random_int(0, len(proxies)-1)])
	}

	apiUrl := "https://backend.cwumsy.cc/backend-api/v2/conversation"

	// JSONify the body and add it to the request
	body_json, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}

	request, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(body_json))
	if err != nil {
		return &http.Response{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	request.Header.Set("Accept", "text/event-stream")
	if err != nil {
		return &http.Response{}, err
	}
	response, err := client.Do(request)
	return response, err
}
