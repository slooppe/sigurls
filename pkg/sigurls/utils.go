package sigurls

import (
	"net/url"
	"strings"
)

func decodeChars(URL string) string {
	source, err := url.QueryUnescape(URL)
	if err == nil {
		URL = source
	}

	// In case json encoded chars
	replacer := strings.NewReplacer(
		`\u002f`, "/",
		`\u0026`, "&",
	)

	// URL = replacer.Replace(strings.ToLower(URL))
	URL = replacer.Replace(URL)

	return URL
}
