package commoncrawl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
)

// Source is the passive scraping agent
type Source struct{}

// CommonPaginationResult is a
type CommonPaginationResult struct {
	Blocks   uint `json:"blocks"`
	PageSize uint `json:"pageSize"`
	Pages    uint `json:"pages"`
}

// CommonResult is a
type CommonResult struct {
	URL   string `json:"url"`
	Error string `json:"error"`
}

// CommonAPIResult is a
type CommonAPIResult []struct {
	API string `json:"cdx-api"`
}

var apiURL string

func formatURL(domain string, page uint, includeSubs bool) string {
	if includeSubs {
		domain = "*." + domain
	}

	return fmt.Sprintf("%s?url=%s/*&output=json&fl=url&page=%d", apiURL, domain, page)
}

// Fetch the number of pages.
func getPagination(domain string, ses *session.Session, includeSubs bool) (*CommonPaginationResult, error) {
	res, err := ses.SimpleGet(fmt.Sprintf("%s&showNumPages=true", formatURL(domain, 0, includeSubs)))
	if err != nil {
		ses.DiscardHTTPResponse(res)
		return nil, err
	}

	defer res.Body.Close()

	var paginationResult CommonPaginationResult

	body, err := ioutil.ReadAll(res.Body)

	// Marshall json response
	if err := json.Unmarshal(body, &paginationResult); err != nil {
		return nil, err
	}

	return &paginationResult, nil
}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, ses *session.Session, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		res, err := ses.SimpleGet("http://index.commoncrawl.org/collinfo.json")
		if err != nil {
			ses.DiscardHTTPResponse(res)
			return
		}

		defer res.Body.Close()

		var apiRresults CommonAPIResult

		body, err := ioutil.ReadAll(res.Body)

		// Marshall json response
		if err := json.Unmarshal(body, &apiRresults); err != nil {
			return
		}

		apiURL = apiRresults[0].API

		// pagination
		pagination, err := getPagination(domain, ses, includeSubs)
		if err != nil {
			fmt.Println(err)
		}

		// URLS
		for page := uint(0); page < pagination.Pages; page++ {
			res, err := ses.SimpleGet(formatURL(domain, page, includeSubs))
			if err != nil {
				ses.DiscardHTTPResponse(res)
				return
			}

			defer res.Body.Close()

			sc := bufio.NewScanner(res.Body)

			for sc.Scan() {
				var result CommonResult

				if err := json.Unmarshal(sc.Bytes(), &result); err != nil {
					return
				}

				if result.Error != "" {
					return
				}

				URLs <- sources.URLs{Source: source.Name(), Value: result.URL}
			}
		}
	}()

	return URLs
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "commoncrawl"
}
