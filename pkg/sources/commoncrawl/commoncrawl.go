package commoncrawl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/valyala/fasthttp"
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
func getPagination(domain string, includeSubs bool) (*CommonPaginationResult, error) {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	req.SetRequestURI(fmt.Sprintf("%s&showNumPages=true", formatURL(domain, 0, includeSubs)))

	client := &fasthttp.Client{}
	if err := client.Do(req, res); err != nil {
		return nil, err
	}

	var paginationResult CommonPaginationResult

	if err := json.NewDecoder(bytes.NewReader(res.Body())).Decode(&paginationResult); err != nil {
		return nil, err
	}

	return &paginationResult, nil
}

// Run function returns all subdomains found with the service
func (source *Source) Run(domain string, includeSubs bool) chan sources.Result {
	URLS := make(chan sources.Result)

	go func() {
		defer close(URLS)

		// collinfo
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()

		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(res)
		}()

		req.SetRequestURI("http://index.commoncrawl.org/collinfo.json")

		client := &fasthttp.Client{}
		if err := client.Do(req, res); err != nil {
			return
		}

		var apiRresults CommonAPIResult

		if err := json.Unmarshal(res.Body(), &apiRresults); err != nil {
			return
		}

		apiURL = apiRresults[0].API

		// pagination
		pagination, err := getPagination(domain, includeSubs)
		if err != nil {
			fmt.Println(err)
		}

		// URLS
		for page := uint(0); page < pagination.Pages; page++ {
			req := fasthttp.AcquireRequest()
			res := fasthttp.AcquireResponse()

			defer func() {
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(res)
			}()

			req.SetRequestURI(formatURL(domain, page, includeSubs))

			client := &fasthttp.Client{}
			if err := client.Do(req, res); err != nil {
				return
			}

			sc := bufio.NewScanner(bytes.NewReader(res.Body()))

			for sc.Scan() {
				var result CommonResult

				if err := json.Unmarshal(sc.Bytes(), &result); err != nil {
					return
				}

				if result.Error != "" {
					return
				}

				URLS <- sources.Result{Source: source.Name(), URL: result.URL}
			}
		}
	}()

	return URLS
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "commoncrawl"
}
