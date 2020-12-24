package runner

import (
	"strings"

	"github.com/drsigned/sigurls/pkg/agent"
	"github.com/drsigned/sigurls/pkg/sources"
)

// Runner is
type Runner struct {
	options *Options
	agent   *agent.Agent
}

// New is
func New(options *Options) *Runner {
	var uses, exclusions []string

	if options.Use != "" {
		uses = append(uses, strings.Split(options.Use, ",")...)
	} else {
		uses = append(uses, sources.All...)
	}

	if options.Exclude != "" {
		exclusions = append(exclusions, strings.Split(options.Exclude, ",")...)
	}

	return &Runner{
		options: options,
		agent:   agent.New(uses, exclusions),
	}
}

// Run is a
func (runner *Runner) Run() (chan sources.URLs, error) {
	URLs := make(chan sources.URLs)

	uniqueMap := make(map[string]sources.URLs)
	sourceMap := make(map[string]map[string]struct{})

	keys := runner.options.YAMLConfig.GetKeys()
	results := runner.agent.Run(runner.options.Domain, keys, runner.options.IncludeSubs)

	go func() {
		defer close(URLs)

		for result := range results {
			URL := result.Value

			if _, exists := uniqueMap[URL]; !exists {
				sourceMap[URL] = make(map[string]struct{})
			}

			sourceMap[URL][result.Source] = struct{}{}

			if _, exists := uniqueMap[URL]; exists {
				continue
			}

			hostEntry := sources.URLs{Source: result.Source, Value: URL}

			uniqueMap[URL] = hostEntry

			URLs <- hostEntry
		}
	}()

	return URLs, nil
}
