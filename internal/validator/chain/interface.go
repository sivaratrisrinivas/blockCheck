package chain

import "context"

// Validator defines the interface that all chain validators must implement
type Validator interface {
	// ValidateAddress checks if the given address is valid for this chain
	ValidateAddress(address string) bool

	// ResolveName resolves a human-readable name to a chain address (e.g. ENS for Ethereum)
	ResolveName(ctx context.Context, name string) (string, error)

	// IsContract checks if the given address is a contract or a regular account
	IsContract(ctx context.Context, address string) (bool, error)

	// GetChainName returns the name of the chain this validator supports
	GetChainName() string
}

// ValidatorConstructor is a function type that creates new validators
type ValidatorConstructor func(config map[string]interface{}) (Validator, error)
