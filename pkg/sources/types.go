package sources

import (
	"github.com/drsigned/sigurls/pkg/session"
)

// URLs is a structure for results
type URLs struct {
	Source string
	Value  string
}

// Source is an interface inherited by each source
type Source interface {
	// Run takes a domain as argument and a session object
	// which contains the extractor for subdomains, http client
	// and other stuff.
	Run(string, *session.Session, bool) chan URLs
	// Name returns the name of the source
	Name() string
}

// Keys contains the current API Keys we have in store
type Keys struct {
	GitHub []string `json:"github"`
}
