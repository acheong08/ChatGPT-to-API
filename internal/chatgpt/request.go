package chatgpt

import (
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

func SendRequest(message typings.ChatGPTRequest, puid *string, access_token string) (*http.Response, error) {
	// JSONify the body and add it to the request
	body_json, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}

	request, err := http.NewRequest(http.MethodPost, "https://chat.openai.com/backend-api/conversation", bytes.NewBuffer(body_json))
	if err != nil {
		return &http.Response{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	request.Header.Set("Accept", "*/*")
	request.AddCookie(&http.Cookie{
		Name:  "_puid",
		Value: *puid,
	})
	if access_token != "" {
		request.Header.Set("Authorization", "Bearer "+access_token)
	}
	if err != nil {
		return &http.Response{}, err
	}
	println(string(body_json))
	response, err := client.Do(request)
	return response, err
}
