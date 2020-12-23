package wayback

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/valyala/fasthttp"
)

// Source is a
type Source struct{}

// Run returns all URLS found from the source.
func (source *Source) Run(domain string, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		if includeSubs {
			domain = "*." + domain
		}

		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()

		defer func() {
			fasthttp.ReleaseRequest(req)
			fasthttp.ReleaseResponse(res)
		}()

		req.SetRequestURI(fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s/*&output=txt&fl=original&collapse=urlkey", domain))

		client := &fasthttp.Client{}
		if err := client.Do(req, res); err != nil {
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(res.Body()))

		for scanner.Scan() {
			URL := scanner.Text()

			if URL == "" {
				continue
			}

			URLs <- sources.URLs{Source: source.Name(), Value: URL}
		}
	}()

	return URLs
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "wayback"
}
