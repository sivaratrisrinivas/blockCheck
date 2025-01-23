package chain

import (
	"fmt"
	"sync"
)

// Registry manages the available chain validators
type Registry struct {
	validators map[string]Validator
	mu         sync.RWMutex
}

// NewRegistry creates a new validator registry
func NewRegistry() *Registry {
	return &Registry{
		validators: make(map[string]Validator),
	}
}

// Register adds a new validator to the registry
func (r *Registry) Register(v Validator) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chainName := v.GetChainName()
	if _, exists := r.validators[chainName]; exists {
		return fmt.Errorf("validator for chain %s already registered", chainName)
	}

	r.validators[chainName] = v
	return nil
}

// Get returns a validator for the specified chain
func (r *Registry) Get(chainName string) (Validator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, exists := r.validators[chainName]
	if !exists {
		return nil, fmt.Errorf("no validator registered for chain %s", chainName)
	}

	return v, nil
}

// ListChains returns a list of all registered chain names
func (r *Registry) ListChains() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chains := make([]string, 0, len(r.validators))
	for chain := range r.validators {
		chains = append(chains, chain)
	}
	return chains
}
