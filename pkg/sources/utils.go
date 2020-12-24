package sources

import "strings"

// NormalizeURL is a
func NormalizeURL(URL string) string {
	// Remove the single and double quotes from the parsed link on the ends
	URL = strings.Trim(URL, "\"")
	URL = strings.Trim(URL, "'")
	// Trim the trailing slash
	URL = strings.TrimRight(URL, "/")
	URL = strings.Trim(URL, " ")

	return URL
}
