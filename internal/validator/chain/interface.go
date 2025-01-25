package chain

import "context"

// Validator defines the interface that all chain validators must implement
type Validator interface {
	// IsValidAddress checks if the given string is a valid address format
	IsValidAddress(address string) bool

	// ResolveENS resolves an ENS name to its Ethereum address
	ResolveENS(name string) (string, error)

	// IsContract checks if the given address is a contract
	IsContract(ctx context.Context, address string) (bool, error)

	// GetChainName returns the name of the chain this validator supports
	GetChainName() string
}

// ValidatorConstructor is a function type that creates new validators
type ValidatorConstructor func(config map[string]interface{}) (Validator, error)
