package sigurls

import "errors"

// Options is a
type Options struct {
	Domain         string
	ExcludeSources string
	IncludeSubs    bool
	UseSources     string
}

// ParseOptions is a
func ParseOptions(options *Options) (*Options, error) {
	if options.Domain == "" {
		return options, errors.New("-d, not provided")
	}

	return options, nil
}
