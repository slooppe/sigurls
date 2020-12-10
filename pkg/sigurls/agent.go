package sigurls

import (
	"sync"

	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/drsigned/sigurls/pkg/sources/commoncrawl"
	"github.com/drsigned/sigurls/pkg/sources/otx"
	"github.com/drsigned/sigurls/pkg/sources/urlscan"
	"github.com/drsigned/sigurls/pkg/sources/wayback"
)

// Agent is a
type Agent struct {
	Sources map[string]sources.Source
}

// NewAgent is a
func NewAgent(Sources, exclusions []string) *Agent {
	agent := &Agent{
		Sources: make(map[string]sources.Source),
	}

	// Add Sources
	for _, source := range Sources {
		switch source {
		case "commoncrawl":
			agent.Sources[source] = &commoncrawl.Source{}
		case "otx":
			agent.Sources[source] = &otx.Source{}
		case "urlscan":
			agent.Sources[source] = &urlscan.Source{}
		case "wayback":
			agent.Sources[source] = &wayback.Source{}
		}
	}

	// Exclude Sources
	for _, source := range exclusions {
		delete(agent.Sources, source)
	}

	return agent
}

// Fetch is a
func (agent *Agent) Fetch(domain string, includeSubs bool) chan sources.Result {
	results := make(chan sources.Result)

	go func() {
		defer close(results)

		wg := new(sync.WaitGroup)

		for source, runner := range agent.Sources {
			wg.Add(1)

			go func(source string, runner sources.Source) {
				for res := range runner.Run(domain, includeSubs) {
					results <- res
				}

				wg.Done()
			}(source, runner)
		}

		wg.Wait()
	}()

	return results
}
