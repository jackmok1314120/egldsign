package egld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"net/http"
	"sync"
	"time"
)

const (
	minimumCachingInterval = time.Second
)

type ArgsElrondBaseProxy struct {
	expirationTime    time.Duration
	httpClientWrapper httpClientWrapper
	endpointProvider  EndpointProvider
}

type elrondBaseProxy struct {
	httpClientWrapper
	mut                 sync.RWMutex
	fetchedConfigs      *NetworkConfig
	lastFetchedTime     time.Time
	cacheExpiryDuration time.Duration
	sinceTimeHandler    func(t time.Time) time.Duration
	endpointProvider    EndpointProvider
}

// newElrondBaseProxy will create a base elrond proxy with cache instance
func newElrondBaseProxy(args ArgsElrondBaseProxy) (*elrondBaseProxy, error) {
	err := checkArgsBaseProxy(args)
	if err != nil {
		return nil, err
	}

	return &elrondBaseProxy{
		httpClientWrapper:   args.httpClientWrapper,
		cacheExpiryDuration: args.expirationTime,
		endpointProvider:    args.endpointProvider,
		sinceTimeHandler:    since,
	}, nil
}

func checkArgsBaseProxy(args ArgsElrondBaseProxy) error {
	if args.expirationTime < minimumCachingInterval {
		return fmt.Errorf("%w, provided: %v, minimum: %v", ErrInvalidCacherDuration, args.expirationTime, minimumCachingInterval)
	}
	if IfNil(args.httpClientWrapper) {
		return ErrNilHTTPClientWrapper
	}
	if IfNil(args.endpointProvider) {
		return ErrNilEndpointProvider
	}

	return nil
}

func since(t time.Time) time.Duration {
	return time.Since(t)
}

// GetNetworkConfig will return the cached network configs fetching new values and saving them if necessary
func (proxy *elrondBaseProxy) GetNetworkConfig(ctx context.Context) (*NetworkConfig, error) {
	proxy.mut.RLock()
	cachedConfigs := proxy.getCachedConfigs()
	proxy.mut.RUnlock()

	if cachedConfigs != nil {
		return cachedConfigs, nil
	}

	return proxy.cacheConfigs(ctx)
}

func (proxy *elrondBaseProxy) getCachedConfigs() *NetworkConfig {
	if proxy.sinceTimeHandler(proxy.lastFetchedTime) > proxy.cacheExpiryDuration {
		return nil
	}

	return proxy.fetchedConfigs
}

func (proxy *elrondBaseProxy) cacheConfigs(ctx context.Context) (*NetworkConfig, error) {
	proxy.mut.Lock()
	defer proxy.mut.Unlock()

	// maybe another parallel running go routine already did the fetching
	cachedConfig := proxy.getCachedConfigs()
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	log.Debug("Network config not cached. caching...")
	configs, err := proxy.getNetworkConfigFromSource(ctx)
	if err != nil {
		return nil, err
	}

	proxy.lastFetchedTime = time.Now()
	proxy.fetchedConfigs = configs

	return configs, nil
}

// getNetworkConfigFromSource retrieves the network configuration from the proxy
func (proxy *elrondBaseProxy) getNetworkConfigFromSource(ctx context.Context) (*NetworkConfig, error) {
	buff, code, err := proxy.GetHTTP(ctx, proxy.endpointProvider.GetNetworkConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &NetworkConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetShardOfAddress returns the shard ID of a provided address by using a shardCoordinator object and querying the
// network config route

//func (proxy *elrondBaseProxy) GetShardOfAddress(ctx context.Context, bech32Address string) (uint32, error) {
//	addr, err := NewAddressFromBech32String(bech32Address)
//	if err != nil {
//		return 0, err
//	}
//
//	networkConfigs, err := proxy.GetNetworkConfig(ctx)
//	if err != nil {
//		return 0, err
//	}
//
//	shardCoordinatorInstance, err := NewShardCoordinator(networkConfigs.NumShardsWithoutMeta, 0)
//	if err != nil {
//		return 0, err
//	}
//
//	return shardCoordinatorInstance.ComputeShardId(addr)
//}

// GetNetworkStatus will return the network status of a provided shard
func (proxy *elrondBaseProxy) GetNetworkStatus(ctx context.Context, shardID uint32) (*NetworkStatus, error) {
	fmt.Println("shardid:", shardID)
	endpoint := proxy.endpointProvider.GetNodeStatus(shardID)
	fmt.Println("endpoint:", endpoint)
	buff, code, err := proxy.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	endpointProviderType := proxy.endpointProvider.GetRestAPIEntityType()
	switch endpointProviderType {
	case Proxy:
		return proxy.getNetworkStatus(buff, shardID)
	case ObserverNode:
		return proxy.getNodeStatus(buff, shardID)
	}

	return &NetworkStatus{}, ErrInvalidEndpointProvider
}

func (proxy *elrondBaseProxy) getNetworkStatus(buff []byte, shardID uint32) (*NetworkStatus, error) {
	response := &NetworkStatusResponse{}
	err := json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	err = proxy.checkReceivedNodeStatus(response.Data.Status, shardID)
	if err != nil {
		return nil, err
	}

	return response.Data.Status, nil
}

func (proxy *elrondBaseProxy) getNodeStatus(buff []byte, shardID uint32) (*NetworkStatus, error) {
	response := &NodeStatusResponse{}
	err := json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	err = proxy.checkReceivedNodeStatus(response.Data.Status, shardID)
	if err != nil {
		return nil, err
	}

	return response.Data.Status, nil
}

func (proxy *elrondBaseProxy) checkReceivedNodeStatus(networkStatus *NetworkStatus, shardID uint32) error {
	if networkStatus == nil {
		return fmt.Errorf("%w, requested from %d", ErrNilNetworkStatus, shardID)
	}
	if !proxy.endpointProvider.ShouldCheckShardIDForNodeStatus() {
		return nil
	}
	if networkStatus.ShardID == shardID {
		return nil
	}

	return fmt.Errorf("%w, requested from %d, got response from %d", ErrShardIDMismatch, shardID, networkStatus.ShardID)
}

// GetRestAPIEntityType returns the REST API entity type that this implementation works with
func (proxy *elrondBaseProxy) GetRestAPIEntityType() RestAPIEntityType {
	return proxy.endpointProvider.GetRestAPIEntityType()
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *elrondBaseProxy) IsInterfaceNil() bool {
	return proxy == nil
}
