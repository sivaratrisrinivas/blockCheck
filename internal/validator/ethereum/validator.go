package ethereum

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/ens"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
)

var (
	addressRegex = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	log          = logrus.New()
)

type EthereumValidator struct {
	client *ethclient.Client
	ens    *ens.Resolver
}

func NewValidator(config map[string]interface{}) (chain.Validator, error) {
	providerURL, ok := config["provider_url"].(string)
	if !ok {
		return nil, fmt.Errorf("provider_url not found in config")
	}

	log.Infof("Initializing Ethereum validator with provider: %s", providerURL)

	client, err := ethclient.Dial(providerURL)
	if err != nil {
		log.Errorf("Failed to connect to Ethereum node: %v", err)
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	// Test connection
	_, err = client.ChainID(context.Background())
	if err != nil {
		log.Errorf("Failed to get chain ID: %v", err)
		return nil, fmt.Errorf("failed to verify connection: %w", err)
	}

	cacheDuration, _ := config["cache_duration"].(int64)
	if cacheDuration == 0 {
		cacheDuration = 3600 // 1 hour default
	}
	log.Debugf("Using cache duration: %d seconds", cacheDuration)

	ensResolver, err := ens.NewResolver(providerURL, time.Duration(cacheDuration)*time.Second)
	if err != nil {
		log.Errorf("Failed to create ENS resolver: %v", err)
		return nil, fmt.Errorf("failed to create ENS resolver: %w", err)
	}

	log.Info("Successfully initialized Ethereum validator")
	return &EthereumValidator{
		client: client,
		ens:    ensResolver,
	}, nil
}

func (v *EthereumValidator) ValidateAddress(address string) bool {
	log.Debugf("Validating address: %s", address)
	return addressRegex.MatchString(address)
}

func (v *EthereumValidator) ResolveName(ctx context.Context, name string) (string, error) {
	log.Infof("Resolving ENS name: %s", name)
	result, err := v.ens.Resolve(ctx, name)
	if err != nil {
		log.Errorf("Failed to resolve ENS name: %v", err)
		return "", err
	}
	if result.Error != "" {
		log.Warnf("ENS resolution error: %s", result.Error)
		return "", fmt.Errorf("%s", result.Error)
	}
	log.Infof("Successfully resolved %s to %s", name, result.Address.Hex())
	return result.Address.Hex(), nil
}

func (v *EthereumValidator) IsContract(ctx context.Context, address string) (bool, error) {
	log.Infof("Checking if address is contract: %s", address)

	if !v.ValidateAddress(address) {
		log.Warn("Invalid address format")
		return false, fmt.Errorf("invalid address format")
	}

	code, err := v.client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		log.Errorf("Failed to get code at address: %v", err)
		return false, fmt.Errorf("failed to get code at address: %w", err)
	}

	isContract := len(code) > 0
	log.Infof("Address %s is contract: %v", address, isContract)
	return isContract, nil
}

func (v *EthereumValidator) GetChainName() string {
	return "ethereum"
}

func (v *EthereumValidator) Close() error {
	log.Info("Closing Ethereum validator")
	v.client.Close()
	return nil
}
