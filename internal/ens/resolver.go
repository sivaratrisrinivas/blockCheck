package ens

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// ENS Registry ABI
const ensRegistryABI = `[{"constant":true,"inputs":[{"name":"node","type":"bytes32"}],"name":"resolver","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"}]`

// ENS Resolver ABI
const ensResolverABI = `[{"constant":true,"inputs":[{"name":"node","type":"bytes32"}],"name":"addr","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"}]`

type Resolver struct {
	client        *ethclient.Client
	cache         map[string]cacheEntry
	cacheMutex    sync.RWMutex
	cacheDuration time.Duration
	registryABI   abi.ABI
	resolverABI   abi.ABI
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
	log.Infof("Connecting to Ethereum node at %s", providerURL)
	client, err := ethclient.Dial(providerURL)
	if err != nil {
		log.Errorf("Failed to connect to Ethereum node: %v", err)
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	log.Info("Successfully connected to Ethereum node")

	registryABI, err := abi.JSON(strings.NewReader(ensRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse registry ABI: %w", err)
	}

	resolverABI, err := abi.JSON(strings.NewReader(ensResolverABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse resolver ABI: %w", err)
	}

	return &Resolver{
		client:        client,
		cache:         make(map[string]cacheEntry),
		cacheDuration: cacheDuration,
		registryABI:   registryABI,
		resolverABI:   resolverABI,
	}, nil
}

func (r *Resolver) Resolve(ctx context.Context, name string) (*ResolveResult, error) {
	// Normalize the name
	name = strings.ToLower(strings.TrimSpace(name))
	if !strings.HasSuffix(name, ".eth") {
		name = name + ".eth"
	}

	log.Debugf("Resolving ENS name: %s", name)

	// Check cache first
	if addr, ok := r.checkCache(name); ok {
		log.Debugf("Cache hit for %s: %s", name, addr.Hex())
		return &ResolveResult{
			Name:    name,
			Address: addr,
		}, nil
	}

	// Resolve using ENS
	address, err := r.resolveENS(ctx, name)
	if err != nil {
		log.Errorf("Failed to resolve ENS name: %v", err)
		return &ResolveResult{
			Name:  name,
			Error: err.Error(),
		}, nil
	}

	// Update cache
	r.updateCache(name, address)

	log.Infof("Successfully resolved %s to %s", name, address.Hex())
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
	log.Debugf("Resolving ENS name using raw contract calls: %s", name)

	// Calculate namehash
	node := NameHash(name)
	log.Debugf("Calculated namehash for %s: %x", name, node)

	// ENS Registry address on mainnet
	registryAddress := common.HexToAddress("0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e")
	log.Debugf("Using ENS Registry at %s", registryAddress.Hex())

	// Call resolver() function
	data, err := r.registryABI.Pack("resolver", node)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to pack resolver call: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &registryAddress,
		Data: data,
	}

	result, err := r.client.CallContract(ctx, msg, nil)
	if err != nil {
		if strings.Contains(err.Error(), "Unauthorized") {
			return common.Address{}, fmt.Errorf("Infura authentication failed: %w", err)
		}
		return common.Address{}, fmt.Errorf("failed to call resolver: %w", err)
	}

	if len(result) == 0 {
		return common.Address{}, fmt.Errorf("no resolver found for %s", name)
	}

	var resolverAddr common.Address
	if err := r.registryABI.UnpackIntoInterface(&resolverAddr, "resolver", result); err != nil {
		return common.Address{}, fmt.Errorf("failed to unpack resolver address: %w", err)
	}

	if resolverAddr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("no resolver found for %s", name)
	}
	log.Debugf("Found resolver at %s", resolverAddr.Hex())

	// Call addr() function on resolver
	data, err = r.resolverABI.Pack("addr", node)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to pack addr call: %w", err)
	}

	msg = ethereum.CallMsg{
		To:   &resolverAddr,
		Data: data,
	}

	result, err = r.client.CallContract(ctx, msg, nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to call addr: %w", err)
	}

	if len(result) == 0 {
		return common.Address{}, fmt.Errorf("address not found for %s", name)
	}

	var address common.Address
	if err := r.resolverABI.UnpackIntoInterface(&address, "addr", result); err != nil {
		return common.Address{}, fmt.Errorf("failed to unpack address: %w", err)
	}

	if address == (common.Address{}) {
		return common.Address{}, fmt.Errorf("address not found for %s", name)
	}

	return address, nil
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
	if r.client != nil {
		r.client.Close()
	}
}
