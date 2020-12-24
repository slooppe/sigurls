package github

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/tomnomnom/linkheader"
	"github.com/valyala/fasthttp"
)

// Source is the passive sources agent
type Source struct{}

type textMatch struct {
	Fragment string `json:"fragment"`
}

type item struct {
	Name        string      `json:"name"`
	HTMLURL     string      `json:"html_url"`
	TextMatches []textMatch `json:"text_matches"`
}

type response struct {
	TotalCount int    `json:"total_count"`
	Items      []item `json:"items"`
}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, ses *session.Session, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		if len(ses.Keys.GitHub) == 0 {
			return
		}

		tokens := NewTokenManager(ses.Keys.GitHub)

		searchURL := fmt.Sprintf("https://api.github.com/search/code?per_page=100&q=%s&sort=created&order=asc", domain)
		source.Enumerate(searchURL, domainRegexp(domain, includeSubs), tokens, ses, URLs)
	}()

	return URLs
}

// Enumerate is a
func (source *Source) Enumerate(searchURL string, domainRegexp *regexp.Regexp, tokens *Tokens, ses *session.Session, URLs chan sources.URLs) {
	token := tokens.Get()

	if token.RetryAfter > 0 {
		if len(tokens.pool) == 1 {
			time.Sleep(time.Duration(token.RetryAfter) * time.Second)
		} else {
			token = tokens.Get()
		}
	}

	headers := map[string]string{"Accept": "application/vnd.github.v3.text-match+json", "Authorization": "token " + token.Hash}

	// Initial request to GitHub search
	res, err := ses.Get(searchURL, headers)
	isForbidden := res != nil && res.StatusCode == http.StatusForbidden
	if err != nil && !isForbidden {
		ses.DiscardHTTPResponse(res)
		return
	}

	// Retry enumerarion after Retry-After seconds on rate limit abuse detected
	ratelimitRemaining, _ := strconv.ParseInt(string(res.Header.Get("X-Ratelimit-Remaining")), 10, 64)
	if isForbidden && ratelimitRemaining == 0 {
		retryAfterSeconds, _ := strconv.ParseInt(string(res.Header.Get("Retry-After")), 10, 64)
		tokens.setCurrentTokenExceeded(retryAfterSeconds)

		source.Enumerate(searchURL, domainRegexp, tokens, ses, URLs)
	}

	var results response

	body, err := ioutil.ReadAll(res.Body)

	// Marshall json response
	if err := json.Unmarshal(body, &results); err != nil {
		return
	}

	err = proccesItems(results.Items, domainRegexp, source.Name(), ses, URLs)
	if err != nil {
		return
	}

	// Links header, first, next, last...
	linksHeader := linkheader.Parse(string(res.Header.Get("Link")))
	// Process the next link recursively
	for _, link := range linksHeader {
		if link.Rel == "next" {
			nextURL, err := url.QueryUnescape(link.URL)
			if err != nil {
				return
			}
			source.Enumerate(nextURL, domainRegexp, tokens, ses, URLs)
		}
	}
}

// proccesItems procceses github response items
func proccesItems(items []item, domainRegexp *regexp.Regexp, name string, ses *session.Session, URLs chan sources.URLs) error {
	for _, item := range items {
		// find URLs in code
		res, err := ses.SimpleGet(rawContentURL(item.HTMLURL))
		if err != nil {
			if res != nil && res.StatusCode != http.StatusNotFound {
				ses.DiscardHTTPResponse(res)
			}
			return err
		}

		if res.StatusCode == fasthttp.StatusOK {
			scanner := bufio.NewScanner(res.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}

				for _, subdomain := range domainRegexp.FindAllString(normalizeContent(line), -1) {
					y := strings.ReplaceAll(subdomain, "&quot;", "\"")
					y = strings.ReplaceAll(y, "/>", "\"")
					y = strings.ReplaceAll(y, "'", "\"")
					y = strings.ReplaceAll(y, "`", "\"")
					y = strings.ReplaceAll(y, ",", "\"")
					y = strings.ReplaceAll(y, "*", "\"")
					y = strings.ReplaceAll(y, ")", "\"")
					y = strings.ReplaceAll(y, "<", "\"")
					y = strings.ReplaceAll(y, ">", "\"")
					y = strings.ReplaceAll(y, "]", "\"")
					x := strings.Split(y, "\"")
					URLs <- sources.URLs{Source: name, Value: x[0]}
				}
			}
		}

		// find subdomains in text matches
		for _, textMatch := range item.TextMatches {
			for _, subdomain := range domainRegexp.FindAllString(normalizeContent(textMatch.Fragment), -1) {
				y := strings.ReplaceAll(subdomain, "&quot;", "\"")
				y = strings.ReplaceAll(y, "/>", "\"")
				y = strings.ReplaceAll(y, "'", "\"")
				y = strings.ReplaceAll(y, "`", "\"")
				y = strings.ReplaceAll(y, ",", "\"")
				y = strings.ReplaceAll(y, "*", "\"")
				y = strings.ReplaceAll(y, ")", "\"")
				y = strings.ReplaceAll(y, "<", "\"")
				y = strings.ReplaceAll(y, ">", "\"")
				y = strings.ReplaceAll(y, "]", "\"")
				x := strings.Split(y, "\"")
				URLs <- sources.URLs{Source: name, Value: x[0]}
			}
		}
	}
	return nil
}

func normalizeContent(content string) string {
	content, _ = url.QueryUnescape(content)
	content = strings.ReplaceAll(content, "\\t", "")
	content = strings.ReplaceAll(content, "\\n", "")
	return content
}

func rawContentURL(URL string) string {
	URL = strings.ReplaceAll(URL, "https://github.com/", "https://raw.githubusercontent.com/")
	URL = strings.ReplaceAll(URL, "/blob/", "/")
	return URL
}

// DomainRegexp regular expression to match subdomains in github files code
func domainRegexp(host string, includeSubs bool) (URLRegex *regexp.Regexp) {
	escapedHost := strings.ReplaceAll(host, ".", "\\.")

	if includeSubs {
		URLRegex = regexp.MustCompile(fmt.Sprintf(`(https?)://[^\s?#\/]*%s/?[^\s]*`, escapedHost))
	} else {
		URLRegex = regexp.MustCompile(fmt.Sprintf(`(https?)://[^\s?#\/]%s/?[^\s]*`, escapedHost))
	}

	return URLRegex
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "github"
}
