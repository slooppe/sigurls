package urlscan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/drsigned/gos"
	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
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
func (source *Source) Run(domain string, ses *session.Session, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		res, err := ses.SimpleGet(fmt.Sprintf("https://urlscan.io/api/v1/search/?q=domain:%s", domain))
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

		for _, i := range results.Results {
			parsedURL, err := gos.ParseURL(i.Page.URL)
			if err != nil {
				continue
			}

			if parsedURL.ETLDPlus1 == domain {
				if includeSubs {
					URLs <- sources.URLs{Source: source.Name(), Value: i.Page.URL}
				} else {
					if parsedURL.SubDomainName == "" || parsedURL.SubDomainName == "www" {
						URLs <- sources.URLs{Source: source.Name(), Value: i.Page.URL}
					}
				}
			}
		}
	}()

	return URLs
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "urlscan"
}
