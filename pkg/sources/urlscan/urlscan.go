package urlscan

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/valyala/fasthttp"
)

type response struct {
	Results []struct {
		Page struct {
			URL string `json:"url"`
		} `json:"page"`
	} `json:"results"`
}

// Source is a
type Source struct{}

// Run returns all URLS found from the source.
func (source *Source) Run(domain string, includeSubs bool) chan sources.Result {
	URLS := make(chan sources.Result)

	go func() {
		defer close(URLS)

		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()

		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(res)
		}()

		req.SetRequestURI(fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain))

		client := &fasthttp.Client{}
		if err := client.Do(req, res); err != nil {
			return
		}

		var results response

		if err := json.Unmarshal(res.Body(), &results); err != nil {
			return
		}

		for _, i := range results.Results {
			parsedURL, err := gos.ParseURL(i.Page.URL)
			if err != nil {
				continue
			}

			if parsedURL.ETLDPlus1 == domain {
				if includeSubs {
					URLS <- sources.Result{Source: source.Name(), URL: i.Page.URL}
				} else {
					if parsedURL.SubDomainName == "" || parsedURL.SubDomainName == "www" {
						URLS <- sources.Result{Source: source.Name(), URL: i.Page.URL}
					}
				}
			}
		}
	}()

	return URLS
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "urlscan"
}
