package otx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
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
func (source *Source) Run(domain string, ses *session.Session, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		for page := 0; ; page++ {
			res, err := ses.SimpleGet(fmt.Sprintf("https://otx.alienvault.com/api/v1/indicators/domain/%s/url_list?limit=%d&page=%d", domain, 200, page))
			if err != nil {
				ses.DiscardHTTPResponse(res)
				return
			}

			defer res.Body.Close()

			var results response

			body, err := ioutil.ReadAll(res.Body)

			if err := json.Unmarshal(body, &results); err != nil {
				return
			}

			for _, i := range results.URLList {
				parsedURL, err := gos.ParseURL(i.URL)
				if err != nil {
					continue
				}

				if parsedURL.ETLDPlus1 == domain {
					if includeSubs {
						URLs <- sources.URLs{Source: source.Name(), Value: i.URL}
					} else {
						if parsedURL.SubDomainName == "" || parsedURL.SubDomainName == "www" {
							URLs <- sources.URLs{Source: source.Name(), Value: i.URL}
						}
					}
				}
			}

			if !results.HasNext {
				break
			}
		}
	}()

	return URLs
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "otx"
}
