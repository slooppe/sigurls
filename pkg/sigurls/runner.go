package sigurls

import (
	"strings"

	"github.com/drsigned/sigurls/pkg/sources"
)

// Runner is
type Runner struct {
	options *Options
	agent   *Agent
}

// NewRunner is
func NewRunner(options *Options) *Runner {
	var uses, exclusions []string

	if options.UseSources != "" {
		uses = append(uses, strings.Split(options.UseSources, ",")...)
	} else {
		uses = append(uses, sources.All...)
	}

	if options.ExcludeSources != "" {
		exclusions = append(exclusions, strings.Split(options.ExcludeSources, ",")...)
	}

	return &Runner{
		options: options,
		agent:   NewAgent(uses, exclusions),
	}
}

// Run is a
func (runner *Runner) Run() (chan sources.Result, error) {
	results := runner.agent.Fetch(runner.options.Domain, runner.options.IncludeSubs)

	URLs := make(chan sources.Result)

	uniqueMap := make(map[string]sources.Result)
	sourceMap := make(map[string]map[string]struct{})

	go func() {
		defer close(URLs)

		for result := range results {
			URL := result.URL

			if _, exists := uniqueMap[URL]; !exists {
				sourceMap[URL] = make(map[string]struct{})
			}

			sourceMap[URL][result.Source] = struct{}{}

			if _, exists := uniqueMap[URL]; exists {
				continue
			}

			hostEntry := sources.Result{URL: URL, Source: result.Source}

			uniqueMap[URL] = hostEntry

			URLs <- hostEntry
		}
	}()

	return URLs, nil
}
