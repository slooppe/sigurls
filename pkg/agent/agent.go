package agent

import (
	"sync"

	"github.com/drsigned/sigurls/pkg/session"
	"github.com/drsigned/sigurls/pkg/sources"
	"github.com/drsigned/sigurls/pkg/sources/commoncrawl"
	"github.com/drsigned/sigurls/pkg/sources/github"
	"github.com/drsigned/sigurls/pkg/sources/otx"
	"github.com/drsigned/sigurls/pkg/sources/urlscan"
	"github.com/drsigned/sigurls/pkg/sources/wayback"
)

// Agent is a
type Agent struct {
	Sources map[string]sources.Source
}

// New is a
func New(Sources, exclusions []string) (agent *Agent) {
	agent = &Agent{
		Sources: make(map[string]sources.Source),
	}

	// Add Sources
	for _, source := range Sources {
		switch source {
		case "commoncrawl":
			agent.Sources[source] = &commoncrawl.Source{}
		case "github":
			agent.Sources[source] = &github.Source{}
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
func (agent *Agent) Fetch(domain string, keys session.Keys, includeSubs bool) chan sources.URLs {
	URLs := make(chan sources.URLs)

	go func() {
		defer close(URLs)

		ses, err := session.New(domain, 10, keys)
		if err != nil {
			return
		}

		wg := new(sync.WaitGroup)

		for source, runner := range agent.Sources {
			wg.Add(1)

			go func(source string, runner sources.Source) {
				for res := range runner.Run(domain, ses, includeSubs) {
					URLs <- res
				}

				wg.Done()
			}(source, runner)
		}

		wg.Wait()
	}()

	return URLs
}
