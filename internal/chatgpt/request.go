package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"

	typings "freechatgpt/internal/typings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Chrome_110),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
)

func constructRequest(body typings.ChatGPTRequest, puid string) (*http.Request, error) {
	// JSONify the body and add it to the request
	body_json, err := json.Marshal(body)
	if err != nil {
		return &http.Request{}, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://chat.openai.com/backend-api/conversation", bufio.NewReader(bytes.NewBuffer(body_json)))
	request.Header.Set("Host", "chat.openai.com")
	request.Header.Set("Origin", "https://chat.openai.com/chat")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("Keep-Alive", "timeout=360")
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	request.AddCookie(&http.Cookie{
		Name:  "_puid",
		Value: puid,
	})
	return request, err
}

func SendRequest(message typings.ChatGPTRequest, puid string) (*http.Response, error) {
	request, err := constructRequest(message, puid)
	if err != nil {
		return &http.Response{}, err
	}
	response, err := client.Do(request)
	return response, err
}
