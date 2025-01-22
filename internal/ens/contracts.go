// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ens

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ENSRegistryABI is the input ABI used to generate the binding from.
const ENSRegistryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"resolver\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"}]"

// ENSRegistry is an auto generated Go binding around an Ethereum contract.
type ENSRegistry struct {
	ENSRegistryCaller     // Read-only binding to the contract
	ENSRegistryTransactor // Write-only binding to the contract
	ENSRegistryFilterer   // Log filterer for contract events
}

// ENSRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type ENSRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ENSRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ENSRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewENSRegistry creates a new instance of ENSRegistry, bound to a specific deployed contract.
func NewENSRegistry(address common.Address, backend bind.ContractBackend) (*ENSRegistry, error) {
	contract, err := bindENSRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ENSRegistry{ENSRegistryCaller: ENSRegistryCaller{contract: contract}, ENSRegistryTransactor: ENSRegistryTransactor{contract: contract}, ENSRegistryFilterer: ENSRegistryFilterer{contract: contract}}, nil
}

// bindENSRegistry binds a generic wrapper to an already deployed contract.
func bindENSRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
func (c *ENSRegistryCaller) Resolver(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "resolver", node)
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
}

// ENSResolverABI is the input ABI used to generate the binding from.
const ENSResolverABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"addr\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"}]"

// ENSResolver is an auto generated Go binding around an Ethereum contract.
type ENSResolver struct {
	ENSResolverCaller     // Read-only binding to the contract
	ENSResolverTransactor // Write-only binding to the contract
	ENSResolverFilterer   // Log filterer for contract events
}

// ENSResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type ENSResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ENSResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ENSResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewENSResolver creates a new instance of ENSResolver, bound to a specific deployed contract.
func NewENSResolver(address common.Address, backend bind.ContractBackend) (*ENSResolver, error) {
	contract, err := bindENSResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ENSResolver{ENSResolverCaller: ENSResolverCaller{contract: contract}, ENSResolverTransactor: ENSResolverTransactor{contract: contract}, ENSResolverFilterer: ENSResolverFilterer{contract: contract}}, nil
}

// bindENSResolver binds a generic wrapper to an already deployed contract.
func bindENSResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
func (c *ENSResolverCaller) Addr(opts *bind.CallOpts, node [32]byte) (common.Address, error) {
	var out []interface{}
	err := c.contract.Call(opts, &out, "addr", node)
	if err != nil {
		return common.Address{}, err
	}
	return *abi.ConvertType(out[0], new(common.Address)).(*common.Address), nil
} 