package sources

import (
	"strings"

	"github.com/drsigned/sigurls/pkg/session"

	"github.com/drsigned/gos"
)

// NormalizeURL is a
func NormalizeURL(URL string, scope session.Scope) (string, bool) {
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

	if parsedURL.ETLDPlus1 == "" || parsedURL.ETLDPlus1 != scope.Domain {
		return URL, false
	}

	if !scope.IncludeSubs {
		if parsedURL.Host != scope.Domain && parsedURL.Host != "www."+scope.Domain {
			return URL, false
		}
	}

	return parsedURL.String(), true
}
