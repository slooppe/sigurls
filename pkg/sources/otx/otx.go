package otx

import (
	"encoding/json"
	"fmt"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/valyala/fasthttp"
)

// Source is a
type Source struct{}

type response struct {
	HasNext    bool `json:"has_next"`
	ActualSize int  `json:"actual_size"`
	URLList    []struct {
		Domain   string `json:"domain"`
		URL      string `json:"url"`
		Hostname string `json:"hostname"`
		HTTPCode int    `json:"httpcode"`
		PageNum  int    `json:"page_num"`
		FullSize int    `json:"full_size"`
		Paged    bool   `json:"paged"`
	} `json:"url_list"`
}

// Run returns all URLS found from the source.
func (source *Source) Run(domain string, includeSubs bool) chan sources.Result {
	URLS := make(chan sources.Result)

	go func() {
		defer close(URLS)

		for page := 0; ; page++ {
			req := fasthttp.AcquireRequest()
			res := fasthttp.AcquireResponse()

			defer func() {
				fasthttp.ReleaseRequest(req)
				fasthttp.ReleaseResponse(res)
			}()

			req.SetRequestURI(fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list?limit=%d&page=%d", domain, 200, page))

			client := &fasthttp.Client{}
			if err := client.Do(req, res); err != nil {
				return
			}

			var results response

			if err := json.Unmarshal(res.Body(), &results); err != nil {
				return
			}

			for _, i := range results.URLList {
				parsedURL, err := gos.ParseURL(i.URL)
				if err != nil {
					continue
				}

				if parsedURL.ETLDPlus1 == domain {
					if includeSubs {
						URLS <- sources.Result{Source: source.Name(), URL: i.URL}
					} else {
						if parsedURL.SubDomainName == "" || parsedURL.SubDomainName == "www" {
							URLS <- sources.Result{Source: source.Name(), URL: i.URL}
						}
					}
				}
			}

			if !results.HasNext {
				break
			}
		}
	}()

	return URLS
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "otx"
}
