package sources

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
	Run(string, bool) chan URLs
	// Name returns the name of the source
	Name() string
}
