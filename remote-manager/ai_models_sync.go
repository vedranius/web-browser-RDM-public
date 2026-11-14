package main

import (
	"context"
	"strings"
)

// aiModelsMaxChars bounds the comma-separated model list of one provider.
const aiModelsMaxChars = 20000

// aiSyncProviderModels connects to the provider, lists every model it offers (for
// Ollama and other local servers: every installed model) and adds the ones that are
// not on the provider's list yet. The models already listed stay first, so the
// default model and the order the administrator chose do not change. It returns the
// added models; the provider is saved only when something was added.
func aiSyncProviderModels(ctx context.Context, p *aiProvider) ([]string, error) {
	avail, err := aiListModels(ctx, *p)
	if err != nil {
		return nil, err
	}
	have := p.modelList()
	size := len(strings.Join(have, ","))
	added := []string{}
	for _, m := range avail {
		if oneOf(m, have...) {
			continue
		}
		if size+len(m)+1 > aiModelsMaxChars {
			break
		}
		have = append(have, m)
		added = append(added, m)
		size += len(m) + 1
	}
	if len(added) == 0 {
		return added, nil
	}
	p.Models = strings.Join(have, ",")
	if p.DefaultModel == "" {
		p.DefaultModel = have[0]
	}
	return added, saveAIProvider(p)
}
