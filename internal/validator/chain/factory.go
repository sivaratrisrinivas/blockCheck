package chain

import (
	"fmt"
	"sync"
)

// Factory manages validator constructors and creates validator instances
type Factory struct {
	constructors map[string]ValidatorConstructor
	mu           sync.RWMutex
}

// NewFactory creates a new validator factory
func NewFactory() *Factory {
	return &Factory{
		constructors: make(map[string]ValidatorConstructor),
	}
}

// Register adds a new validator constructor to the factory
func (f *Factory) Register(chainName string, constructor ValidatorConstructor) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.constructors[chainName]; exists {
		return fmt.Errorf("constructor for chain %s already registered", chainName)
	}

	f.constructors[chainName] = constructor
	return nil
}

// Create instantiates a new validator for the specified chain
func (f *Factory) Create(chainName string, config map[string]interface{}) (Validator, error) {
	f.mu.RLock()
	constructor, exists := f.constructors[chainName]
	f.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no constructor registered for chain %s", chainName)
	}

	return constructor(config)
}

// ListSupportedChains returns a list of all chains that have registered constructors
func (f *Factory) ListSupportedChains() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	chains := make([]string, 0, len(f.constructors))
	for chain := range f.constructors {
		chains = append(chains, chain)
	}
	return chains
}
