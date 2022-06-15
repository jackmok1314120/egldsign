package egld

import (
	"context"
)

// EgldProxy holds the primitive functions that the elrond proxy engine supports & implements
type EgldProxy interface {
	GetNetworkConfig(ctx context.Context) (*NetworkConfig, error)
	//GetAccount(ctx context.Context, address AddressHandler) (*Account, error)
	SendTransaction(ctx context.Context, tx *Transaction) (string, error)
	SendTransactions(ctx context.Context, txs []*Transaction) ([]string, error)
	IsInterfaceNil() bool
}

type httpClientWrapper interface {
	GetHTTP(ctx context.Context, endpoint string) ([]byte, int, error)
	PostHTTP(ctx context.Context, endpoint string, data []byte) ([]byte, int, error)
	IsInterfaceNil() bool
}

// EndpointProvider is able to return endpoint routes strings
type EndpointProvider interface {
	GetNetworkConfig() string
	GetNetworkEconomics() string
	GetRatingsConfig() string
	GetEnableEpochsConfig() string
	GetAccount(addressAsBech32 string) string
	GetCostTransaction() string
	GetSendTransaction() string
	GetSendMultipleTransactions() string
	GetTransactionStatus(hexHash string) string
	GetTransactionInfo(hexHash string) string
	GetHyperBlockByNonce(nonce int64) string
	GetHyperBlockByHash(hexHash string) string
	GetVmValues() string
	GetGenesisNodesConfig() string
	GetRawStartOfEpochMetaBlock(epoch uint32) string
	GetNodeStatus(shardID uint32) string
	ShouldCheckShardIDForNodeStatus() bool
	GetRawBlockByHash(shardID uint32, hexHash string) string
	GetRawBlockByNonce(shardID uint32, nonce int64) string
	GetRawMiniBlockByHash(shardID uint32, hexHash string, epoch uint32) string
	GetRestAPIEntityType() RestAPIEntityType
	IsInterfaceNil() bool
}

// FinalityProvider is able to check the shard finalization status
type FinalityProvider interface {
	CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta int64) error
	IsInterfaceNil() bool
}
type AddressHandler interface {
	AddressAsBech32String() string
	AddressBytes() []byte
	AddressSlice() [32]byte
	IsValid() bool
	IsInterfaceNil() bool
}

// ProxyHandler defines the behavior of a proxy handler that can process requests
type ProxyHandler interface {
	GetLatestHyperBlockNonce(ctx context.Context) (int64, error)
	GetHyperBlockByNonce(ctx context.Context, nonce int64) (*HyperBlock, error)
	GetDefaultTransactionArguments(ctx context.Context, address AddressHandler, networkConfigs *NetworkConfig) (ArgCreateTransaction, error)
	GetNetworkConfig(ctx context.Context) (*NetworkConfig, error)
	IsInterfaceNil() bool
}

// Coordinator defines what a shard state coordinator should hold
type Coordinator interface {
	NumberOfShards() uint32
	ComputeId(address []byte) uint32
	SelfId() uint32
	SameShard(firstAddress, secondAddress []byte) bool
	CommunicationIdentifier(destShardID uint32) string
	IsInterfaceNil() bool
}

// EpochHandler defines what a component which handles current epoch should be able to do
type EpochHandler interface {
	MetaEpoch() uint32
	IsInterfaceNil() bool
}

//PeerAccountListAndRatingHandler provides Rating Computation Capabilites for the Nodes Coordinator and ValidatorStatistics
type PeerAccountListAndRatingHandler interface {
	//GetChance returns the chances for the the rating
	GetChance(uint32) uint32
	//GetStartRating gets the start rating values
	GetStartRating() uint32
	//GetSignedBlocksThreshold gets the threshold for the minimum signed blocks
	GetSignedBlocksThreshold() float32
	//ComputeIncreaseProposer computes the new rating for the increaseLeader
	ComputeIncreaseProposer(shardId uint32, currentRating uint32) uint32
	//ComputeDecreaseProposer computes the new rating for the decreaseLeader
	ComputeDecreaseProposer(shardId uint32, currentRating uint32, consecutiveMisses uint32) uint32
	//RevertIncreaseValidator computes the new rating if a revert for increaseProposer should be done
	RevertIncreaseValidator(shardId uint32, currentRating uint32, nrReverts uint32) uint32
	//ComputeIncreaseValidator computes the new rating for the increaseValidator
	ComputeIncreaseValidator(shardId uint32, currentRating uint32) uint32
	//ComputeDecreaseValidator computes the new rating for the decreaseValidator
	ComputeDecreaseValidator(shardId uint32, currentRating uint32) uint32
	//IsInterfaceNil verifies if the interface is nil
	IsInterfaceNil() bool
}

// GenesisNodesSetupHandler returns the genesis nodes info
type GenesisNodesSetupHandler interface {
	//AllInitialNodes() []GenesisNodeInfoHandler
	InitialNodesPubKeys() map[uint32][]string
	GetShardIDForPubKey(pubkey []byte) (uint32, error)
	InitialEligibleNodesPubKeysForShard(shardId uint32) ([]string, error)
	//InitialNodesInfoForShard(shardId uint32) ([]GenesisNodeInfoHandler, []GenesisNodeInfoHandler, error)
	//InitialNodesInfo() (map[uint32][]GenesisNodeInfoHandler, map[uint32][]GenesisNodeInfoHandler)
	GetStartTime() int64
	GetRoundDuration() uint64
	GetShardConsensusGroupSize() uint32
	GetMetaConsensusGroupSize() uint32
	NumberOfShards() uint32
	MinNumberOfNodes() uint32
	MinNumberOfShardNodes() uint32
	MinNumberOfMetaNodes() uint32
	GetHysteresis() float32
	GetAdaptivity() bool
	MinNumberOfNodesWithHysteresis() uint32
	IsInterfaceNil() bool
}
