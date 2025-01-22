package ens

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Resolver struct {
	client        bind.ContractBackend
	cache         map[string]cacheEntry
	cacheMutex    sync.RWMutex
	cacheDuration time.Duration
}

type cacheEntry struct {
	address   common.Address
	timestamp time.Time
}

type ResolveResult struct {
	Name    string         `json:"name"`
	Address common.Address `json:"address,omitempty"`
	Error   string         `json:"error,omitempty"`
}

func NewResolver(providerURL string, cacheDuration time.Duration) (*Resolver, error) {
	client, err := ethclient.Dial(providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &Resolver{
		client:        client,
		cache:         make(map[string]cacheEntry),
		cacheDuration: cacheDuration,
	}, nil
}

func (r *Resolver) Resolve(ctx context.Context, name string) (*ResolveResult, error) {
	// Normalize the name
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".eth") {
		name = name + ".eth"
	}

	// Check cache first
	if addr, ok := r.checkCache(name); ok {
		return &ResolveResult{
			Name:    name,
			Address: addr,
		}, nil
	}

	// Resolve using ENS
	address, err := r.resolveENS(ctx, name)
	if err != nil {
		return &ResolveResult{
			Name:  name,
			Error: err.Error(),
		}, nil
	}

	// Update cache
	r.updateCache(name, address)

	return &ResolveResult{
		Name:    name,
		Address: address,
	}, nil
}

func (r *Resolver) checkCache(name string) (common.Address, bool) {
	r.cacheMutex.RLock()
	defer r.cacheMutex.RUnlock()

	if entry, exists := r.cache[name]; exists {
		if time.Since(entry.timestamp) < r.cacheDuration {
			return entry.address, true
		}
	}
	return common.Address{}, false
}

func (r *Resolver) updateCache(name string, address common.Address) {
	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	r.cache[name] = cacheEntry{
		address:   address,
		timestamp: time.Now(),
	}
}

func (r *Resolver) resolveENS(ctx context.Context, name string) (common.Address, error) {
	// Calculate namehash
	node := NameHash(name)

	// ENS Registry address on mainnet
	registryAddress := common.HexToAddress("0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e")

	// Create ENS registry contract
	registry, err := NewENSRegistry(registryAddress, r.client)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create ENS registry contract: %w", err)
	}

	// Get resolver address
	resolverAddr, err := registry.Resolver(&bind.CallOpts{Context: ctx}, node)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get resolver address: %w", err)
	}

	if resolverAddr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("no resolver found for %s", name)
	}

	// Create resolver contract
	resolver, err := NewENSResolver(resolverAddr, r.client)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create resolver contract: %w", err)
	}

	// Get address
	addr, err := resolver.Addr(&bind.CallOpts{Context: ctx}, node)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to resolve address: %w", err)
	}

	if addr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("address not found for %s", name)
	}

	return addr, nil
}

// NameHash implements the ENS namehash algorithm
func NameHash(name string) [32]byte {
	if name == "" {
		return [32]byte{}
	}

	node := [32]byte{}

	labels := strings.Split(name, ".")
	for i := len(labels) - 1; i >= 0; i-- {
		labelHash := crypto.Keccak256([]byte(labels[i]))
		node = crypto.Keccak256Hash(append(node[:], labelHash...))
	}

	return node
}

func (r *Resolver) Close() {
	if client, ok := r.client.(*ethclient.Client); ok {
		client.Close()
	}
}
