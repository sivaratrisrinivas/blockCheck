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
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/logger"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
	"go.uber.org/zap"
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

func (v *EthereumValidator) IsValidAddress(address string) bool {
	logger.Debug("Validating Ethereum address",
		zap.String("address", address))
	return addressRegex.MatchString(address)
}

func (v *EthereumValidator) ResolveENS(name string) (string, error) {
	logger.Debug("Resolving ENS name",
		zap.String("name", name))

	result, err := v.ens.Resolve(context.Background(), name)
	if err != nil {
		logger.Error("Failed to resolve ENS name",
			zap.String("name", name),
			zap.Error(err))
		return "", err
	}
	if result.Error != "" {
		logger.Warn("ENS resolution error",
			zap.String("name", name),
			zap.String("error", result.Error))
		return "", fmt.Errorf("%s", result.Error)
	}
	logger.Info("Successfully resolved ENS name",
		zap.String("name", name),
		zap.String("address", result.Address.Hex()))
	return result.Address.Hex(), nil
}

func (v *EthereumValidator) IsContract(ctx context.Context, address string) (bool, error) {
	logger.Debug("Checking if address is contract",
		zap.String("address", address))

	if !v.IsValidAddress(address) {
		logger.Warn("Invalid address format",
			zap.String("address", address))
		return false, fmt.Errorf("invalid address format")
	}

	code, err := v.client.CodeAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		logger.Error("Failed to get code at address",
			zap.String("address", address),
			zap.Error(err))
		return false, err
	}

	isContract := len(code) > 0
	logger.Info("Contract check completed",
		zap.String("address", address),
		zap.Bool("isContract", isContract))
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
