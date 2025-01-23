package ethereum

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/ens"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
)

var addressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

type EthereumValidator struct {
	client *ethclient.Client
	ens    *ens.Resolver
}

func NewValidator(config map[string]interface{}) (chain.Validator, error) {
	providerURL, ok := config["provider_url"].(string)
	if !ok {
		return nil, fmt.Errorf("provider_url not found in config")
	}

	client, err := ethclient.Dial(providerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	cacheDuration, _ := config["cache_duration"].(int64)
	if cacheDuration == 0 {
		cacheDuration = 3600 // 1 hour default
	}

	ensResolver, err := ens.NewResolver(providerURL, time.Duration(cacheDuration)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to create ENS resolver: %w", err)
	}

	return &EthereumValidator{
		client: client,
		ens:    ensResolver,
	}, nil
}

func (v *EthereumValidator) ValidateAddress(address string) bool {
	return addressRegex.MatchString(address)
}

func (v *EthereumValidator) ResolveName(ctx context.Context, name string) (string, error) {
	result, err := v.ens.Resolve(ctx, name)
	if err != nil {
		return "", err
	}
	return result.Address.Hex(), nil
}

func (v *EthereumValidator) IsContract(ctx context.Context, address string) (bool, error) {
	if !v.ValidateAddress(address) {
		return false, fmt.Errorf("invalid Ethereum address format")
	}

	code, err := v.client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return false, fmt.Errorf("failed to get code at address: %w", err)
	}

	return len(code) > 0, nil
}

func (v *EthereumValidator) GetChainName() string {
	return "ethereum"
}

func (v *EthereumValidator) Close() {
	if v.client != nil {
		v.client.Close()
	}
	if v.ens != nil {
		v.ens.Close()
	}
}
