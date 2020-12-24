package wayback

import (
	"bufio"
	"fmt"
	"net/url"
	"strings"

	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
)

// Source is a
type Source struct{}

// Run returns all URLS found from the source.
func (source *Source) Run(domain string, ses *session.Session, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		if includeSubs {
			domain = "*." + domain
		}

		res, err := ses.SimpleGet(fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s/*&output=txt&fl=original&collapse=urlkey", domain))
		if err != nil {
			ses.DiscardHTTPResponse(res)
			return
		}

		defer res.Body.Close()

		scanner := bufio.NewScanner(res.Body)

		for scanner.Scan() {
			URL := scanner.Text()
			if URL == "" {
				continue
			}

			URL, err = url.QueryUnescape(URL)
			if err != nil {
				return
			}

			if URL != "" {
				// fix for triple encoded URL
				URL = strings.ToLower(URL)
				URL = strings.TrimPrefix(URL, "25")
				URL = strings.TrimPrefix(URL, "2f")

				URLs <- sources.URLs{Source: source.Name(), Value: URL}
			}
		}
	}()

	return URLs
}

// Name returns the name of the source
func (source *Source) Name() string {
	return "wayback"
}
