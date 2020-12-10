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
func (source *Source) Run(domain string, includeSubs bool) chan sources.Result {
	URLS := make(chan sources.Result)

	go func() {
		defer close(URLS)

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
			line := scanner.Text()

			if line == "" {
				continue
			}

			URLS <- sources.Result{Source: source.Name(), URL: line}
		}
	}()

	return URLS
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "wayback"
}
