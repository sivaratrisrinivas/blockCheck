package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/types"
)

type AddressCache struct {
	cache Cache
}

type AddressInfo struct {
	IsValid     bool   `json:"is_valid"`
	AddressType string `json:"address_type"`
	ENSName     string `json:"ens_name,omitempty"`
}

const defaultTTL = 1 * time.Hour

func NewAddressCache(cache Cache) *AddressCache {
	return &AddressCache{
		cache: cache,
	}
}

func (ac *AddressCache) GetAddressInfo(address string) (*AddressInfo, error) {
	data, err := ac.cache.Get(context.Background(), "addr:"+address)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	var info AddressInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (ac *AddressCache) SetAddressInfo(address string, info *AddressInfo) error {
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return ac.cache.Set(context.Background(), "addr:"+address, data, defaultTTL)
}

func (ac *AddressCache) GetENSAddress(name string) (string, error) {
	data, err := ac.cache.Get(context.Background(), "ens:"+name)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}
	return string(data), nil
}

func (ac *AddressCache) SetENSAddress(name, address string) error {
	return ac.cache.Set(context.Background(), "ens:"+name, []byte(address), defaultTTL)
}

func (ac *AddressCache) Clear() error {
	return ac.cache.Clear(context.Background())
}

func (ac *AddressCache) GetStats() types.Stats {
	return ac.cache.GetStats()
}
