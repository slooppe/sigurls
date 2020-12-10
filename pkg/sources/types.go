package sources

// Source is an interface inherited by each source
type Source interface {
	Run(string, bool) chan Result
	Name() string
}

// Result is a
type Result struct {
	Source string
	URL    string
}
