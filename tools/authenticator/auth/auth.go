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
	AccessToken        string
	UserAgent          string
	State              string
	URL                string
	Verifier_code      string
	Verifier_challenge string
	AuthDetails        AuthDetails
}

type AuthDetails struct {
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

func NewAuthDetails(challenge string) AuthDetails {
	// Generate state (secrets.token_urlsafe(32))
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return AuthDetails{
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

	auth.AuthDetails = NewAuthDetails(auth.Verifier_challenge)

	return auth
}

func (auth *Authenticator) URLEncode(str string) string {
	return url.QueryEscape(str)
}

func (auth *Authenticator) Begin() Error {
	// Just realized that the client id is hardcoded in the JS file

	return auth.partOne()
}
func (auth *Authenticator) partOne() Error {

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
		"client_id":             {auth.AuthDetails.ClientID},
		"scope":                 {auth.AuthDetails.Scope},
		"response_type":         {auth.AuthDetails.ResponseType},
		"redirect_uri":          {auth.AuthDetails.RedirectURL},
		"audience":              {auth.AuthDetails.Audience},
		"prompt":                {auth.AuthDetails.Prompt},
		"state":                 {auth.AuthDetails.State},
		"code_challenge":        {auth.AuthDetails.CodeChallenge},
		"code_challenge_method": {auth.AuthDetails.CodeChallengeMethod},
	}
	auth_url = auth_url + "?" + payload.Encode()
	req, err := http.NewRequest("GET", auth_url, nil)
	if err != nil {
		return *NewError("part_one", 0, "", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return *NewError("part_one", 0, "", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *NewError("part_one", 0, "", err)
	}

	if resp.StatusCode == 302 {
		return auth.partTwo("https://auth0.openai.com" + resp.Header.Get("Location"))
	} else {
		err := NewError("part_one", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))
		return *err
	}
}

func (auth *Authenticator) partTwo(url string) Error {

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Connection":      "keep-alive",
		"User-Agent":      auth.UserAgent,
		"Accept-Language": "en-US,en;q=0.9",
		"Referer":         "https://ios.chat.openai.com/",
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return *NewError("part_two", 0, "", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return *NewError("part_two", 0, "", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return *NewError("part_two", 0, "", err)
	}

	if resp.StatusCode == 302 || resp.StatusCode == 200 {

		stateRegex := regexp.MustCompile(`state=(.*)`)
		stateMatch := stateRegex.FindStringSubmatch(string(body))
		if len(stateMatch) < 2 {
			return *NewError("part_two", 0, "Could not find state in response", fmt.Errorf("error: Check details"))
		}

		state := strings.Split(stateMatch[1], `"`)[0]
		return auth.partThree(state)
	} else {
		err := NewError("__part_two", resp.StatusCode, string(body), fmt.Errorf("error: Check details"))
		return *err
	}
}
func (auth *Authenticator) partThree(state string) Error {

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

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return *NewError("part_four", 0, "", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return *NewError("part_four", 0, "", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 || resp.StatusCode == 200 {
		return auth.partFive(state)
	} else {
		err := NewError("__part_four", resp.StatusCode, "Your email address is invalid.", fmt.Errorf("error: Check details"))
		return *err
	}

}
func (auth *Authenticator) partFive(state string) Error {

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

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return *NewError("part_five", 0, "", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return *NewError("part_five", 0, "", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		redirectURL := resp.Header.Get("Location")
		return auth.partSix(state, redirectURL)
	} else {
		body := bytes.NewBuffer(nil)
		_, err1 := body.ReadFrom(resp.Body)
		if err1 != nil {
			return *NewError("part_five", 0, "", err1)
		}
		err := NewError("__part_five", resp.StatusCode, body.String(), fmt.Errorf("error: Check details"))
		return *err
	}

}
func (auth *Authenticator) partSix(oldState string, redirectURL string) Error {

	url := "https://auth0.openai.com" + redirectURL

	headers := map[string]string{
		"Host":            "auth0.openai.com",
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Connection":      "keep-alive",
		"User-Agent":      auth.UserAgent,
		"Accept-Language": "en-GB,en-US;q=0.9,en;q=0.8",
		"Referer":         fmt.Sprintf("https://auth0.openai.com/u/login/password?state=%s", oldState),
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return *NewError("part_six", 0, "", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := auth.Session.Do(req)
	if err != nil {
		return *NewError("part_six", 0, "", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		auth.URL = resp.Header.Get("Location")

		return Error{}
	} else {
		err := NewError("__part_six", resp.StatusCode, resp.Status, fmt.Errorf("error: Check details"))
		return *err
	}

}
func (auth *Authenticator) GetAccessToken() (string, Error) {
	code := regexp.MustCompile(`code=(.*)&`).FindStringSubmatch(auth.URL)
	if len(code) == 0 {
		err := NewError("__get_access_token", 0, auth.URL, fmt.Errorf("error: Check details"))
		return "", *err
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
		return "", *NewError("get_access_token", 0, "", err)
	}
	defer resp.Body.Close()
	// Parse response
	body, _ := io.ReadAll(resp.Body)
	// Parse as JSON
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)

	if err != nil {
		return "", *NewError("get_access_token", 0, "", err)
	}

	// Check if access token in data
	if _, ok := data["access_token"]; !ok {
		return "", *NewError("get_access_token", 0, "Missing access token", fmt.Errorf("error: Check details"))
	}

	return data["access_token"].(string), Error{}
}
