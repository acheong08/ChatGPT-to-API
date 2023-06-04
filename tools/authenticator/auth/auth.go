package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"crypto/rand"

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
	AuthRequest        AuthRequest
	AuthResult         AuthResult
}

type AuthRequest struct {
	ClientID            string `json:"client_id"`
	Scope               string `json:"scope"`
	ResponseType        string `json:"response_type"`
	RedirectURL         string `json:"redirect_url"`
	Audience            string `json:"audience"`
	Prompt              string `json:"prompt"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
}

type AuthResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	PUID         string `json:"puid"`
}

func NewAuthDetails(challenge string) AuthRequest {
	// Generate state (secrets.token_urlsafe(32))
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return AuthRequest{
		ClientID:            "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
		Scope:               "openid email profile offline_access model.request model.read organization.read",
		ResponseType:        "code",
		RedirectURL:         "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		Audience:            "https://api.openai.com/v1",
		Prompt:              "login",
		State:               state,
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	}
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
		tls_client.WithClientProfile(tls_client.Firefox_102),
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

	auth.AuthRequest = NewAuthDetails(auth.Verifier_challenge)

	return auth
}

func (auth *Authenticator) URLEncode(str string) string {
	return url.QueryEscape(str)
}

func (auth *Authenticator) Begin() *Error {
	// Just realized that the client id is hardcoded in the JS file

	return auth.partOne()
}
func (auth *Authenticator) partOne() *Error {

	auth_url := "https://auth0.openai.com/authorize"
	headers := map[string]string{
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
	payload := url.Values{
		"client_id":             {auth.AuthRequest.ClientID},
		"scope":                 {auth.AuthRequest.Scope},
		"response_type":         {auth.AuthRequest.ResponseType},
		"redirect_uri":          {auth.AuthRequest.RedirectURL},
		"audience":              {auth.AuthRequest.Audience},
		"prompt":                {auth.AuthRequest.Prompt},
		"state":                 {auth.AuthRequest.State},
		"code_challenge":        {auth.AuthRequest.CodeChallenge},
		"code_challenge_method": {auth.AuthRequest.CodeChallengeMethod},
	}
	auth_url = auth_url + "?" + payload.Encode()
	req, _ := http.NewRequest("GET", auth_url, nil)

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

	if resp.StatusCode == 302 {
		return auth.partTwo("https://auth0.openai.com" + resp.Header.Get("Location"))
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
		"Referer":         "https://ios.chat.openai.com/",
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
		return NewError("__part_two", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))

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
		return NewError("part_four", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 200 {
		return auth.partFour(state)
	} else {
		return NewError("__part_four", resp.StatusCode, "Your email address is invalid.", fmt.Errorf("error: Check details"))

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
		return NewError("part_five", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		redirectURL := resp.Header.Get("Location")
		return auth.partFive(state, redirectURL)
	} else {
		body := bytes.NewBuffer(nil)
		body.ReadFrom(resp.Body)
		return NewError("__part_five", resp.StatusCode, body.String(), fmt.Errorf("error: Check details"))

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
		return NewError("part_six", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		auth.URL = resp.Header.Get("Location")
		return auth.partSix()
	} else {
		return NewError("__part_six", resp.StatusCode, resp.Status, fmt.Errorf("error: Check details"))

	}

}
func (auth *Authenticator) partSix() *Error {
	code := regexp.MustCompile(`code=(.*)&`).FindStringSubmatch(auth.URL)
	if len(code) == 0 {
		return NewError("__get_access_token", 0, auth.URL, fmt.Errorf("error: Check details"))
	}
	payload, _ := json.Marshal(map[string]string{
		"redirect_uri":  "com.openai.chat://auth0.openai.com/ios/com.openai.chat/callback",
		"grant_type":    "authorization_code",
		"client_id":     "pdlLIX2Y72MIl2rhLhTE9VV9bN905kBh",
		"code":          code[1],
		"code_verifier": auth.Verifier_code,
		"state":         auth.State,
	})

	req, _ := http.NewRequest("POST", "https://auth0.openai.com/oauth/token", strings.NewReader(string(payload)))
	for k, v := range map[string]string{
		"User-Agent":   auth.UserAgent,
		"content-type": "application/json",
	} {
		req.Header.Set(k, v)
	}
	resp, err := auth.Session.Do(req)
	if err != nil {
		return NewError("get_access_token", 0, "Failed to send request", err)
	}
	defer resp.Body.Close()
	// Parse response
	body, _ := io.ReadAll(resp.Body)
	// Parse as JSON
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return NewError("get_access_token", 0, "Response was not JSON", err)
	}

	// Check if access token in data
	if _, ok := data["access_token"]; !ok {
		return NewError("get_access_token", 0, "Missing access token", fmt.Errorf("error: Check details"))
	}
	auth.AuthResult.AccessToken = data["access_token"].(string)
	auth.AuthResult.RefreshToken = data["refresh_token"].(string)

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
