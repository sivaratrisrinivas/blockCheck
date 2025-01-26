package ethereum

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/ens"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/logger"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/validator/chain"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
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

func (v *EthereumValidator) IsChecksumAddress(address string) bool {
	logger.Debug("Validating Ethereum address checksum",
		zap.String("address", address))

	// First check basic format
	if !addressRegex.MatchString(address) {
		logger.Debug("Basic format check failed")
		return false
	}

	// Generate checksum address
	checksummed, err := ToChecksumAddress(address)
	if err != nil {
		logger.Debug("Failed to generate checksum address",
			zap.Error(err))
		return false
	}

	logger.Debug("Comparing addresses",
		zap.String("input", address),
		zap.String("checksummed", checksummed))

	// For EIP-55, we need exact match
	return address == checksummed
}

// ToChecksumAddress converts an Ethereum address to mixed-case checksum format
func ToChecksumAddress(address string) (string, error) {
	if !addressRegex.MatchString(address) {
		return "", fmt.Errorf("invalid ethereum address format")
	}

	// Remove 0x prefix and convert to lowercase
	addr := strings.ToLower(address[2:])

	// Calculate hash of the lowercase address without 0x prefix
	hash := Keccak256([]byte(addr))

	result := "0x"
	for i := 0; i < len(addr); i++ {
		// Get the corresponding nibble from the hash
		// Each byte of the hash corresponds to two characters of the address
		hashByte := hash[i/2]
		// For even indices, use the high nibble
		// For odd indices, use the low nibble
		hashNibble := hashByte
		if i%2 == 0 {
			hashNibble = hashByte >> 4
		} else {
			hashNibble &= 0xf
		}

		// If the character is a letter (a-f) and the corresponding hash nibble is >= 8,
		// make it uppercase
		c := addr[i]
		if c >= '0' && c <= '9' {
			result += string(c)
		} else {
			if hashNibble >= 8 {
				result += strings.ToUpper(string(c))
			} else {
				result += string(c)
			}
		}
	}

	return result, nil
}

// Keccak256 calculates the Keccak-256 hash of a byte slice
func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
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
