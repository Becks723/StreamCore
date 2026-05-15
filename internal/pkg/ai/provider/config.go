package provider

import (
	"log"

	"StreamCore/config"
)

// Registry is an in-memory lookup of providers by name.
type Registry map[string]config.ProviderConfig

// BuildRegistry creates a provider name → config map.
func BuildRegistry(providers []config.ProviderConfig) Registry {
	r := make(Registry, len(providers))
	for _, p := range providers {
		r[p.Name] = p
	}
	return r
}

// Get returns the provider config for the given name, or nil if not found.
func (r Registry) Get(name string) (*config.ProviderConfig, bool) {
	p, ok := r[name]
	if !ok {
		return nil, false
	}
	return &p, true
}

// ResolveModel returns bot_config.model_name if set, otherwise provider's default_model.
func ResolveModel(p *config.ProviderConfig, modelOverride string) string {
	if modelOverride != "" {
		return modelOverride
	}
	if p.DefaultModel != "" {
		return p.DefaultModel
	}
	return "gpt-4o"
}

var globalRegistry Registry

// Init builds the global provider registry from the loaded config.
func Init() {
	cfg := config.Instance()
	globalRegistry = BuildRegistry(cfg.Providers)
	if len(globalRegistry) == 0 {
		log.Fatalf("[ai] provider.Init: no providers configured")
	}
	log.Printf("[ai] loaded %d providers", len(globalRegistry))
}

// GlobalRegistry returns the global provider registry.
func GlobalRegistry() Registry {
	return globalRegistry
}
