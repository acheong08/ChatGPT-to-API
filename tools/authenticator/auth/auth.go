package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	pkce "github.com/nirasan/go-oauth-pkce-code-verifier"
)

type Error struct {
	Location   string
	StatusCode int
	Details    string
	Error      error
}

func NewError(location string, statusCode int, details string, err error) *Error {
	return &Error{
		Location:   location,
		StatusCode: statusCode,
		Details:    details,
		Error:      err,
	}
}

type Authenticator struct {
	EmailAddress       string
	Password           string
	Proxy              string
	Session            tls_client.HttpClient
	UserAgent          string
	State              string
	URL                string
	Verifier_code      string
	Verifier_challenge string
	AuthResult         AuthResult
}

type AuthResult struct {
	AccessToken string `json:"access_token"`
	PUID        string `json:"puid"`
}

func NewAuthenticator(emailAddress, password, proxy string) *Authenticator {
	auth := &Authenticator{
		EmailAddress: emailAddress,
		Password:     password,
		Proxy:        proxy,
		UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	}
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(20),
		tls_client.WithClientProfile(tls_client.Okhttp4Android13),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		// Proxy
		tls_client.WithProxyUrl(proxy),
	}
	auth.Session, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	// PKCE
	verifier, _ := pkce.CreateCodeVerifier()
	auth.Verifier_code = verifier.String()
	auth.Verifier_challenge = verifier.CodeChallengeS256()

	return auth
}

func (auth *Authenticator) URLEncode(str string) string {
	return url.QueryEscape(str)
}

func (auth *Authenticator) Begin() *Error {

	url := "https://chat.openai.com/api/auth/csrf"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NewError("begin", 0, "", err)
	}

	req.Header.Set("Host", "chat.openai.com")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", auth.UserAgent)
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://chat.openai.com/auth/login")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("begin", 0, "", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewError("begin", 0, "", err)
	}

	if resp.StatusCode == 200 && strings.Contains(resp.Header.Get("Content-Type"), "json") {

		var csrfTokenResponse struct {
			CsrfToken string `json:"csrfToken"`
		}
		err = json.Unmarshal(body, &csrfTokenResponse)
		if err != nil {
			return NewError("begin", 0, "", err)
		}

		csrfToken := csrfTokenResponse.CsrfToken
		return auth.partOne(csrfToken)
	} else {
		err := NewError("begin", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))
		return err
	}
}

func (auth *Authenticator) partOne(csrfToken string) *Error {

	auth_url := "https://chat.openai.com/api/auth/signin/auth0?prompt=login"
	headers := map[string]string{
		"Host":            "chat.openai.com",
		"User-Agent":      auth.UserAgent,
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept":          "*/*",
		"Sec-Gpc":         "1",
		"Accept-Language": "en-US,en;q=0.8",
		"Origin":          "https://chat.openai.com",
		"Sec-Fetch-Site":  "same-origin",
		"Sec-Fetch-Mode":  "cors",
		"Sec-Fetch-Dest":  "empty",
		"Referer":         "https://chat.openai.com/auth/login",
		"Accept-Encoding": "gzip, deflate",
	}

	// Construct payload
	payload := fmt.Sprintf("callbackUrl=%%2F&csrfToken=%s&json=true", csrfToken)
	req, _ := http.NewRequest("POST", auth_url, strings.NewReader(payload))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_one", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewError("part_one", 0, "Failed to read body", err)
	}

	if resp.StatusCode == 200 && strings.Contains(resp.Header.Get("Content-Type"), "json") {
		var urlResponse struct {
			URL string `json:"url"`
		}
		err = json.Unmarshal(body, &urlResponse)
		if err != nil {
			return NewError("part_one", 0, "Failed to decode JSON", err)
		}
		if urlResponse.URL == "https://chat.openai.com/api/auth/error?error=OAuthSignin" || strings.Contains(urlResponse.URL, "error") {
			err := NewError("part_one", resp.StatusCode, "You have been rate limited. Please try again later.", fmt.Errorf("error: Check details"))
			return err
		}
		return auth.partTwo(urlResponse.URL)
	} else {
		return NewError("part_one", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))
	}
}

func (auth *Authenticator) partTwo(url string) *Error {

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Connection":      "keep-alive",
		"User-Agent":      auth.UserAgent,
		"Accept-Language": "en-US,en;q=0.9",
	}

	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_two", 0, "Failed to make request", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 302 || resp.StatusCode == 200 {

		stateRegex := regexp.MustCompile(`state=(.*)`)
		stateMatch := stateRegex.FindStringSubmatch(string(body))
		if len(stateMatch) < 2 {
			return NewError("part_two", 0, "Could not find state in response", fmt.Errorf("error: Check details"))
		}

		state := strings.Split(stateMatch[1], `"`)[0]
		return auth.partThree(state)
	} else {
		return NewError("part_two", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))

	}
}
func (auth *Authenticator) partThree(state string) *Error {

	url := fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state)
	emailURLEncoded := auth.URLEncode(auth.EmailAddress)

	payload := fmt.Sprintf(
		"state=%s&username=%s&js-available=false&webauthn-available=true&is-brave=false&webauthn-platform-available=true&action=default",
		state, emailURLEncoded,
	)

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Origin":          "https://auth0.openai.com",
		"Connection":      "keep-alive",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"User-Agent":      auth.UserAgent,
		"Referer":         fmt.Sprintf("https://auth0.openai.com/u/login/identifier?state=%s", state),
		"Accept-Language": "en-US,en;q=0.9",
		"Content-Type":    "application/x-www-form-urlencoded",
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_three", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 200 {
		return auth.partFour(state)
	} else {
		return NewError("part_three", resp.StatusCode, "Your email address is invalid.", fmt.Errorf("error: Check details"))

	}

}
func (auth *Authenticator) partFour(state string) *Error {

	url := fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state)
	emailURLEncoded := auth.URLEncode(auth.EmailAddress)
	passwordURLEncoded := auth.URLEncode(auth.Password)
	payload := fmt.Sprintf("state=%s&username=%s&password=%s&action=default", state, emailURLEncoded, passwordURLEncoded)

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Origin":          "https://auth0.openai.com",
		"Connection":      "keep-alive",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"User-Agent":      auth.UserAgent,
		"Referer":         fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", state),
		"Accept-Language": "en-US,en;q=0.9",
		"Content-Type":    "application/x-www-form-urlencoded",
	}

	req, _ := http.NewRequest("POST", url, strings.NewReader(payload))

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_four", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		redirectURL := resp.Header.Get("Location")
		return auth.partFive(state, redirectURL)
	} else {
		body := bytes.NewBuffer(nil)
		body.ReadFrom(resp.Body)
		return NewError("part_four", resp.StatusCode, body.String(), fmt.Errorf("error: Check details"))

	}

}
func (auth *Authenticator) partFive(oldState string, redirectURL string) *Error {

	url := "https://auth0.openai.com" + redirectURL

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Connection":      "keep-alive",
		"User-Agent":      auth.UserAgent,
		"Accept-Language": "en-GB,en-US;q=0.9,en;q=0.8",
		"Referer":         fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState),
	}

	req, _ := http.NewRequest("GET", url, nil)

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_five", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		return auth.partSix(resp.Header.Get("Location"), url)
	} else {
		return NewError("part_five", resp.StatusCode, resp.Status, fmt.Errorf("error: Check details"))

	}

}
func (auth *Authenticator) partSix(url, redirect_url string) *Error {
	req, _ := http.NewRequest("GET", url, nil)
	for k, v := range map[string]string{
		"Host":            "chat.openai.com",
		"Accept":          "application/json",
		"Connection":      "keep-alive",
		"User-Agent":      auth.UserAgent,
		"Accept-Language": "en-GB,en-US;q=0.9,en;q=0.8",
		"Referer":         redirect_url,
	} {
		req.Header.Set(k, v)
	}
	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("part_six", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()
	if err != nil {
		return NewError("part_six", 0, "Response was not JSON", err)
	}
	if resp.StatusCode != 302 {
		return NewError("part_six", resp.StatusCode, url, fmt.Errorf("incorrect response code"))
	}
	// Check location header
	if location := resp.Header.Get("Location"); location != "https://chat.openai.com/" {
		return NewError("part_six", resp.StatusCode, location, fmt.Errorf("incorrect redirect"))
	}

	url = "https://chat.openai.com/api/auth/session"

	req, _ = http.NewRequest("GET", url, nil)

	// Set user agent
	req.Header.Set("User-Agent", auth.UserAgent)

	resp, err = auth.Session.Do(req)
	if err != nil {
		return NewError("get_access_token", 0, "Failed to send request", err)
	}

	if resp.StatusCode != 200 {
		return NewError("get_access_token", resp.StatusCode, "Incorrect response code", fmt.Errorf("error: Check details"))
	}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return NewError("get_access_token", 0, "", err)
	}

	// Check if access token in data
	if _, ok := result["accessToken"]; !ok {
		result_string := fmt.Sprintf("%v", result)
		return NewError("part_six", 0, result_string, fmt.Errorf("missing access token"))
	}
	auth.AuthResult.AccessToken = result["accessToken"].(string)

	return nil
}

func (auth *Authenticator) GetAccessToken() string {
	return auth.AuthResult.AccessToken
}

func (auth *Authenticator) GetPUID() (string, *Error) {
	// Check if user has access token
	if auth.AuthResult.AccessToken == "" {
		return "", NewError("get_puid", 0, "Missing access token", fmt.Errorf("error: Check details"))
	}
	// Make request to https://chat.openai.com/backend-api/models
	req, _ := http.NewRequest("GET", "https://chat.openai.com/backend-api/models", nil)
	// Add headers
	req.Header.Add("Authorization", "Bearer "+auth.AuthResult.AccessToken)
	req.Header.Add("User-Agent", auth.UserAgent)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Referer", "https://chat.openai.com/")
	req.Header.Add("Origin", "https://chat.openai.com")
	req.Header.Add("Connection", "keep-alive")

	resp, err := auth.Session.Do(req)
	if err != nil {
		return "", NewError("get_puid", 0, "Failed to make request", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", NewError("get_puid", resp.StatusCode, "Failed to make request", fmt.Errorf("error: Check details"))
	}
	// Find `_puid` cookie in response
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "_puid" {
			auth.AuthResult.PUID = cookie.Value
			return cookie.Value, nil
		}
	}
	// If cookie not found, return error
	return "", NewError("get_puid", 0, "PUID cookie not found", fmt.Errorf("error: Check details"))
}

func (auth *Authenticator) GetAuthResult() AuthResult {
	return auth.AuthResult
}
