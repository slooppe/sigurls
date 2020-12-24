package session

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Keys contains the current API Keys we have in store
type Keys struct {
	GitHub []string `json:"github"`
}

// Scope is the scope contorl structure
type Scope struct {
	Domain      string
	IncludeSubs bool
}

// Session is the option passed to the source, an option is created
// uniquely for eac source.
type Session struct {
	Scope Scope
	// Client is the current http client
	Client *http.Client
	// Keys is the API keys for the application
	Keys Keys
}

// New creates a new session object for a domain
func New(domain string, includeSubs bool, timeout int, keys Keys) (*Session, error) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(timeout) * time.Second,
	}

	return &Session{
		Scope: Scope{
			Domain:      domain,
			IncludeSubs: includeSubs,
		},
		Client: client,
		Keys:   keys,
	}, nil
}

// SimpleGet makes a simple GET request to a URL
func (session *Session) SimpleGet(getURL string) (*http.Response, error) {
	return session.HTTPRequest(http.MethodGet, getURL, map[string]string{}, nil)
}

// Get makes a GET request to a URL with extended parameters
func (session *Session) Get(getURL string, headers map[string]string) (*http.Response, error) {
	return session.HTTPRequest(http.MethodGet, getURL, headers, nil)
}

// HTTPRequest makes any HTTP request to a URL with extended parameters
func (session *Session) HTTPRequest(method, requestURL string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Connection", "close")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return httpRequestWrapper(session.Client, req)
}

// DiscardHTTPResponse discards the response content by demand
func (session *Session) DiscardHTTPResponse(response *http.Response) {
	if response != nil {
		_, err := io.Copy(ioutil.Discard, response.Body)
		if err != nil {
			return
		}
		response.Body.Close()
	}
}

func httpRequestWrapper(client *http.Client, request *http.Request) (*http.Response, error) {
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		requestURL, _ := url.QueryUnescape(request.URL.String())
		return res, fmt.Errorf("unexpected status code %d received from %s", res.StatusCode, requestURL)
	}
	return res, nil
}
