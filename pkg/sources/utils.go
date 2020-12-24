package sources

import (
	"strings"

	"github.com/drsigned/gos"
)

// NormalizeURL is a
func NormalizeURL(URL, domain string, includeSubs bool) (string, bool) {
	// Remove the single and double quotes from the parsed link on the ends
	URL = strings.Trim(URL, "\"")
	URL = strings.Trim(URL, "'")
	// Trim the trailing slash
	URL = strings.TrimRight(URL, "/")
	URL = strings.Trim(URL, " ")

	parsedURL, err := gos.ParseURL(URL)
	if err != nil {
		return URL, false
	}

	if parsedURL.Host == "" || parsedURL.ETLDPlus1 != domain {
		return URL, false
	}

	return parsedURL.String(), true
}
