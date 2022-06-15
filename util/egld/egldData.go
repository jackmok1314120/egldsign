package egld

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"
)

// NetworkConfigResponse holds the network config endpoint response
type NetworkConfigResponse struct {
	Data struct {
		Config *NetworkConfig `json:"config"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// NetworkConfig holds the network configuration parameters
type NetworkConfig struct {
	ChainID                  string  `json:"erd_chain_id"`
	Denomination             int     `json:"erd_denomination"`
	GasPerDataByte           uint64  `json:"erd_gas_per_data_byte"`
	LatestTagSoftwareVersion string  `json:"erd_latest_tag_software_version"`
	MetaConsensusGroup       uint32  `json:"erd_meta_consensus_group_size"`
	MinGasLimit              uint64  `json:"erd_min_gas_limit"`
	MinGasPrice              uint64  `json:"erd_min_gas_price"`
	MinTransactionVersion    uint32  `json:"erd_min_transaction_version"`
	NumMetachainNodes        uint32  `json:"erd_num_metachain_nodes"`
	NumNodesInShard          uint32  `json:"erd_num_nodes_in_shard"`
	NumShardsWithoutMeta     uint32  `json:"erd_num_shards_without_meta"`
	RoundDuration            int64   `json:"erd_round_duration"`
	ShardConsensusGroupSize  uint64  `json:"erd_shard_consensus_group_size"`
	StartTime                int64   `json:"erd_start_time"`
	Adaptivity               bool    `json:"erd_adaptivity,string"`
	Hysteresys               float32 `json:"erd_hysteresis,string"`
}

var errInvalidBalance = errors.New("invalid balance")

// AccountResponse holds the account endpoint response
type AccountResponse struct {
	Data struct {
		Account *Account `json:"account"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// Account holds an Account's information
type Account struct {
	Address         string `json:"address"`
	Nonce           uint64 `json:"nonce"`
	Balance         string `json:"balance"`
	Code            string `json:"code"`
	CodeHash        []byte `json:"codeHash"`
	RootHash        []byte `json:"rootHash"`
	CodeMetadata    []byte `json:"codeMetadata"`
	Username        string `json:"username"`
	DeveloperReward string `json:"developerReward"`
	OwnerAddress    string `json:"ownerAddress"`
}

// GetBalance computes the float representation of the balance,
// based on the provided number of decimals
func (a *Account) GetBalance(decimals int) (float64, error) {
	balance, ok := big.NewFloat(0).SetString(a.Balance)
	if !ok {
		return 0, errInvalidBalance
	}
	// Compute denominated balance to 18 decimals
	denomination := big.NewInt(int64(decimals))
	denominationMultiplier := big.NewInt(10)
	denominationMultiplier.Exp(denominationMultiplier, denomination, nil)
	floatDenomination, _ := big.NewFloat(0).SetString(denominationMultiplier.String())
	balance.Quo(balance, floatDenomination)
	floatBalance, _ := balance.Float64()

	return floatBalance, nil
}

// SendTransactionResponse holds the response received from the network when broadcasting a transaction
type SendTransactionResponse struct {
	Data struct {
		TxHash string `json:"txHash"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// SendTransactionsResponse holds the response received from the network when broadcasting multiple transactions
type SendTransactionsResponse struct {
	Data struct {
		NumOfSentTxs int            `json:"numOfSentTxs"`
		TxsHashes    map[int]string `json:"txsHashes"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// Transaction holds the fields of a transaction to be broadcasted to the network
type Transaction struct {
	Nonce     uint64 `json:"nonce"`
	Value     string `json:"value"`
	RcvAddr   string `json:"receiver"`
	SndAddr   string `json:"sender"`
	GasPrice  uint64 `json:"gasPrice,omitempty"`
	GasLimit  uint64 `json:"gasLimit,omitempty"`
	Data      []byte `json:"data,omitempty"`
	Signature string `json:"signature,omitempty"`
	ChainID   string `json:"chainID"`
	Version   uint32 `json:"version"`
	Options   uint32 `json:"options,omitempty"`
}

// TransactionStatus holds a transaction's status response from the network
type TransactionStatus struct {
	Data struct {
		Status string `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// TransactionInfo holds a transaction info response from the network
type TransactionInfo struct {
	Data struct {
		Transaction TransactionOnNetwork `json:"transaction"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// TransactionOnNetwork holds a transaction's info entry in a hyper block
type TransactionOnNetwork struct {
	Type             string                    `json:"type"`
	Hash             string                    `json:"hash"`
	Nonce            int64                     `json:"nonce"`
	Value            *big.Int                  `json:"value"`
	Receiver         string                    `json:"receiver"`
	Sender           string                    `json:"sender"`
	GasPrice         *big.Int                  `json:"gasPrice"`
	GasLimit         int64                     `json:"gasLimit"`
	Data             []byte                    `json:"data"`
	Signature        string                    `json:"signature"`
	SourceShard      uint32                    `json:"sourceShard"`
	DestinationShard uint32                    `json:"destinationShard"`
	MiniblockType    string                    `json:"miniblockType"`
	MiniblockHash    string                    `json:"miniblockHash"`
	Status           string                    `json:"status"`
	HyperBlockNonce  int64                     `json:"hyperblockNonce"`
	HyperBlockHash   string                    `json:"hyperblockHash"`
	ScResults        []*ApiSmartContractResult `json:"smartContractResults,omitempty"`
}

// TxCostResponseData follows the format of the data field of a transaction cost request
type TxCostResponseData struct {
	TxCost     uint64 `json:"txGasUnits"`
	RetMessage string `json:"returnMessage"`
}

// ResponseTxCost defines a response from the node holding the transaction cost
type ResponseTxCost struct {
	Data  TxCostResponseData `json:"data"`
	Error string             `json:"error"`
	Code  string             `json:"code"`
}

// ArgCreateTransaction will hold the transaction fields
type ArgCreateTransaction struct {
	Nonce            uint64
	Value            string
	RcvAddr          string
	SndAddr          string
	GasPrice         uint64
	GasLimit         uint64
	Data             []byte
	Signature        string
	ChainID          string
	Version          uint32
	Options          uint32
	AvailableBalance string
}

// ApiTransactionResult is the data transfer object which will be returned on the get transaction by hash endpoint
type ApiTransactionResult struct {
	Tx                                TransactionHandler        `json:"-"`
	Type                              string                    `json:"type"`
	Hash                              string                    `json:"hash,omitempty"`
	Nonce                             uint64                    `json:"nonce,omitempty"`
	Round                             uint64                    `json:"round,omitempty"`
	Epoch                             uint32                    `json:"epoch,omitempty"`
	Value                             string                    `json:"value,omitempty"`
	Receiver                          string                    `json:"receiver,omitempty"`
	Sender                            string                    `json:"sender,omitempty"`
	SenderUsername                    []byte                    `json:"senderUsername,omitempty"`
	ReceiverUsername                  []byte                    `json:"receiverUsername,omitempty"`
	GasPrice                          uint64                    `json:"gasPrice,omitempty"`
	GasLimit                          uint64                    `json:"gasLimit,omitempty"`
	Data                              []byte                    `json:"data,omitempty"`
	CodeMetadata                      []byte                    `json:"codeMetadata,omitempty"`
	Code                              string                    `json:"code,omitempty"`
	PreviousTransactionHash           string                    `json:"previousTransactionHash,omitempty"`
	OriginalTransactionHash           string                    `json:"originalTransactionHash,omitempty"`
	ReturnMessage                     string                    `json:"returnMessage,omitempty"`
	OriginalSender                    string                    `json:"originalSender,omitempty"`
	Signature                         string                    `json:"signature,omitempty"`
	SourceShard                       uint32                    `json:"sourceShard"`
	DestinationShard                  uint32                    `json:"destinationShard"`
	BlockNonce                        uint64                    `json:"blockNonce,omitempty"`
	BlockHash                         string                    `json:"blockHash,omitempty"`
	NotarizedAtSourceInMetaNonce      uint64                    `json:"notarizedAtSourceInMetaNonce,omitempty"`
	NotarizedAtSourceInMetaHash       string                    `json:"NotarizedAtSourceInMetaHash,omitempty"`
	NotarizedAtDestinationInMetaNonce uint64                    `json:"notarizedAtDestinationInMetaNonce,omitempty"`
	NotarizedAtDestinationInMetaHash  string                    `json:"notarizedAtDestinationInMetaHash,omitempty"`
	MiniBlockType                     string                    `json:"miniblockType,omitempty"`
	MiniBlockHash                     string                    `json:"miniblockHash,omitempty"`
	Timestamp                         int64                     `json:"timestamp,omitempty"`
	Receipt                           *ApiReceipt               `json:"receipt,omitempty"`
	SmartContractResults              []*ApiSmartContractResult `json:"smartContractResults,omitempty"`
	Logs                              *ApiLogs                  `json:"logs,omitempty"`
	Status                            TxStatus                  `json:"status,omitempty"`
	Tokens                            []string                  `json:"tokens,omitempty"`
	ESDTValues                        []string                  `json:"esdtValues,omitempty"`
	Receivers                         []string                  `json:"receivers,omitempty"`
	ReceiversShardIDs                 []uint32                  `json:"receiversShardIDs,omitempty"`
	Operation                         string                    `json:"operation,omitempty"`
	Function                          string                    `json:"function,omitempty"`
	IsRelayed                         bool                      `json:"isRelayed,omitempty"`
}

// ApiSmartContractResult represents a smart contract result with changed fields' types in order to make it friendly for API's json
type ApiSmartContractResult struct {
	Hash              string   `json:"hash,omitempty"`
	Nonce             uint64   `json:"nonce"`
	Value             *big.Int `json:"value"`
	RcvAddr           string   `json:"receiver"`
	SndAddr           string   `json:"sender"`
	RelayerAddr       string   `json:"relayerAddress,omitempty"`
	RelayedValue      *big.Int `json:"relayedValue,omitempty"`
	Code              string   `json:"code,omitempty"`
	Data              string   `json:"data,omitempty"`
	PrevTxHash        string   `json:"prevTxHash"`
	OriginalTxHash    string   `json:"originalTxHash"`
	GasLimit          uint64   `json:"gasLimit"`
	GasPrice          uint64   `json:"gasPrice"`
	CallType          CallType `json:"callType"`
	CodeMetadata      string   `json:"codeMetadata,omitempty"`
	ReturnMessage     string   `json:"returnMessage,omitempty"`
	OriginalSender    string   `json:"originalSender,omitempty"`
	Logs              *ApiLogs `json:"logs,omitempty"`
	Tokens            []string `json:"tokens,omitempty"`
	ESDTValues        []string `json:"esdtValues,omitempty"`
	Receivers         []string `json:"receivers,omitempty"`
	ReceiversShardIDs []uint32 `json:"receiversShardIDs,omitempty"`
	Operation         string   `json:"operation,omitempty"`
	Function          string   `json:"function,omitempty"`
	IsRelayed         bool     `json:"isRelayed,omitempty"`
}

// ApiReceipt represents a receipt with changed fields' types in order to make it friendly for API's json
type ApiReceipt struct {
	Value   *big.Int `json:"value"`
	SndAddr string   `json:"sender"`
	Data    string   `json:"data,omitempty"`
	TxHash  string   `json:"txHash"`
}

// ApiLogs represents logs with changed fields' types in order to make it friendly for API's json
type ApiLogs struct {
	Address string    `json:"address"`
	Events  []*Events `json:"events"`
}

// Events represents the events generated by a transaction with changed fields' types in order to make it friendly for API's json
type Events struct {
	Address    string   `json:"address"`
	Identifier string   `json:"identifier"`
	Topics     [][]byte `json:"topics"`
	Data       []byte   `json:"data"`
}

// CostResponse is structure used to return the transaction cost in gas units
type CostResponse struct {
	GasUnits             uint64                             `json:"txGasUnits"`
	ReturnMessage        string                             `json:"returnMessage"`
	SmartContractResults map[string]*ApiSmartContractResult `json:"smartContractResults"`
}

// NetworkStatusResponse holds the network status response (for a specified shard)
type NetworkStatusResponse struct {
	Data struct {
		Status *NetworkStatus `json:"status"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// NetworkStatus holds the network status details of a specified shard
type NetworkStatus struct {
	CurrentRound               uint64 `json:"erd_current_round"`
	EpochNumber                uint64 `json:"erd_epoch_number"`
	Nonce                      int64  `json:"erd_nonce"`
	NonceAtEpochStart          uint64 `json:"erd_nonce_at_epoch_start"`
	NoncesPassedInCurrentEpoch uint64 `json:"erd_nonces_passed_in_current_epoch"`
	RoundAtEpochStart          uint64 `json:"erd_round_at_epoch_start"`
	RoundsPassedInCurrentEpoch uint64 `json:"erd_rounds_passed_in_current_epoch"`
	RoundsPerEpoch             uint64 `json:"erd_rounds_per_epoch"`
	CrossCheckBlockHeight      string `json:"erd_cross_check_block_height"`
	HighestNonce               int64  `json:"erd_highest_final_nonce"`
	ProbableHighestNonce       int64  `json:"erd_probable_highest_nonce"`
	ShardID                    uint32 `json:"erd_shard_id"`
}

// NodeStatusResponse holds the node status response
type NodeStatusResponse struct {
	Data struct {
		Status *NetworkStatus `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// HyperBlock holds a hyper block's details
type HyperBlock struct {
	Nonce         int64  `json:"nonce"`
	Round         uint64 `json:"round"`
	Hash          string `json:"hash"`
	PrevBlockHash string `json:"prevBlockHash"`
	Epoch         uint64 `json:"epoch"`
	NumTxs        uint64 `json:"numTxs"`
	ShardBlocks   []struct {
		Hash  string `json:"hash"`
		Nonce uint64 `json:"nonce"`
		Shard uint32 `json:"shard"`
	} `json:"shardBlocks"`
	Transactions []*TransactionOnNetwork
	Timestamp    int64 `json:"timestamp"`
}

// HyperBlockResponse holds a hyper block info response from the network
type HyperBlockResponse struct {
	Data struct {
		HyperBlock HyperBlock `json:"hyperblock"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

// RatingsConfigResponse holds the ratings config endpoint response
type RatingsConfigResponse struct {
	Data struct {
		Config *RatingsConfig `json:"config"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// SelectionChances holds the selection chances parameters
type SelectionChances struct {
	ChancePercent uint32 `json:"erd_chance_percent"`
	MaxThreshold  uint32 `json:"erd_max_threshold"`
}

// RatingsConfig holds the ratings configuration parameters
type RatingsConfig struct {
	GeneralMaxRating                          uint32              `json:"erd_ratings_general_max_rating"`
	GeneralMinRating                          uint32              `json:"erd_ratings_general_min_rating"`
	GeneralSignedBlocksThreshold              float32             `json:"erd_ratings_general_signed_blocks_threshold,string"`
	GeneralStartRating                        uint32              `json:"erd_ratings_general_start_rating"`
	GeneralSelectionChances                   []*SelectionChances `json:"erd_ratings_general_selection_chances"`
	MetachainConsecutiveMissedBlocksPenalty   float32             `json:"erd_ratings_metachain_consecutive_missed_blocks_penalty,string"`
	MetachainHoursToMaxRatingFromStartRating  uint32              `json:"erd_ratings_metachain_hours_to_max_rating_from_start_rating"`
	MetachainProposerDecreaseFactor           float32             `json:"erd_ratings_metachain_proposer_decrease_factor,string"`
	MetachainProposerValidatorImportance      float32             `json:"erd_ratings_metachain_proposer_validator_importance,string"`
	MetachainValidatorDecreaseFactor          float32             `json:"erd_ratings_metachain_validator_decrease_factor,string"`
	PeerhonestyBadPeerThreshold               float64             `json:"erd_ratings_peerhonesty_bad_peer_threshold,string"`
	PeerhonestyDecayCoefficient               float64             `json:"erd_ratings_peerhonesty_decay_coefficient,string"`
	PeerhonestyDecayUpdateIntervalInseconds   uint32              `json:"erd_ratings_peerhonesty_decay_update_interval_inseconds"`
	PeerhonestyMaxScore                       float64             `json:"erd_ratings_peerhonesty_max_score,string"`
	PeerhonestyMinScore                       float64             `json:"erd_ratings_peerhonesty_min_score,string"`
	PeerhonestyUnitValue                      float64             `json:"erd_ratings_peerhonesty_unit_value,string"`
	ShardchainConsecutiveMissedBlocksPenalty  float32             `json:"erd_ratings_shardchain_consecutive_missed_blocks_penalty,string"`
	ShardchainHoursToMaxRatingFromStartRating uint32              `json:"erd_ratings_shardchain_hours_to_max_rating_from_start_rating"`
	ShardchainProposerDecreaseFactor          float32             `json:"erd_ratings_shardchain_proposer_decrease_factor,string"`
	ShardchainProposerValidatorImportance     float32             `json:"erd_ratings_shardchain_proposer_validator_importance,string"`
	ShardchainValidatorDecreaseFactor         float32             `json:"erd_ratings_shardchain_validator_decrease_factor,string"`
}

// RawBlockRespone holds the raw blocks endpoint response
type RawBlockRespone struct {
	Data struct {
		Block []byte `json:"block"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

// RawMiniBlockRespone holds the raw miniblock endpoint respone
type RawMiniBlockRespone struct {
	Data struct {
		MiniBlock []byte `json:"miniblock"`
	}
	Error string `json:"error"`
	Code  string `json:"code"`
}

// EnableEpochsConfigResponse holds the enable epochs config endpoint response
type EnableEpochsConfigResponse struct {
	Data struct {
		Config *EnableEpochsConfig `json:"enableEpochs"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// MaxNodesChangeConfig holds the max nodes change config
type MaxNodesChangeConfig struct {
	EpochEnable            uint32 `json:"erd_epoch_enable"`
	MaxNumNodes            uint32 `json:"erd_max_num_nodes"`
	NodesToShufflePerShard uint32 `json:"erd_nodes_to_shuffle_per_shard"`
}

// EnableEpochsConfig holds the enable epochs configuration parameters
type EnableEpochsConfig struct {
	BalanceWaitingListsEnableEpoch uint32                 `json:"erd_balance_waiting_lists_enable_epoch"`
	WaitingListFixEnableEpoch      uint32                 `json:"erd_waiting_list_fix_enable_epoch"`
	MaxNodesChangeEnableEpoch      []MaxNodesChangeConfig `json:"erd_max_nodes_change_enable_epoch"`
}

// GenesisNodesResponse holds the network genesis nodes endpoint reponse
type GenesisNodesResponse struct {
	Data struct {
		Nodes *GenesisNodes `json:"nodes"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// GenesisNodes holds the genesis nodes public keys per shard
type GenesisNodes struct {
	Eligible map[uint32][]string `json:"eligible"`
	Waiting  map[uint32][]string `json:"waiting"`
}

// NetworkEconomicsResponse holds the network economics endpoint response
type NetworkEconomicsResponse struct {
	Data struct {
		Economics *NetworkEconomics `json:"metrics"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

// NetworkEconomics holds the network economics details
type NetworkEconomics struct {
	DevRewards            string `json:"erd_dev_rewards"`
	EpochForEconomicsData uint32 `json:"erd_epoch_for_economics_data"`
	Inflation             string `json:"erd_inflation"`
	TotalFees             string `json:"erd_total_fees"`
	TotalStakedValue      string `json:"erd_total_staked_value"`
	TotalSupply           string `json:"erd_total_supply"`
	TotalTopUpValue       string `json:"erd_total_top_up_value"`
}

// VmValuesResponseData follows the format of the data field in an API response for a VM values query
type VmValuesResponseData struct {
	Data *VMOutputApi `json:"data"`
}

// ResponseVmValue defines a wrapper over string containing returned data in hex format
type ResponseVmValue struct {
	Data  VmValuesResponseData `json:"data"`
	Error string               `json:"error"`
	Code  string               `json:"code"`
}

// VmValueRequest defines the request struct for values available in a VM
type VmValueRequest struct {
	Address    string   `json:"scAddress"`
	FuncName   string   `json:"funcName"`
	CallerAddr string   `json:"caller"`
	CallValue  string   `json:"value"`
	Args       []string `json:"args"`
}

// VmValueRequest defines the request struct for values available in a VM
type VmValueRequestWithOptionalParameters struct {
	*VmValueRequest
	SameScState    bool `json:"sameScState"`
	ShouldBeSynced bool `json:"shouldBeSynced"`
}

// TransactionHandler defines the type of executable transaction
type TransactionHandler interface {
	IsInterfaceNil() bool

	GetValue() *big.Int
	GetNonce() uint64
	GetData() []byte
	GetRcvAddr() []byte
	GetRcvUserName() []byte
	GetSndAddr() []byte
	GetGasLimit() uint64
	GetGasPrice() uint64

	SetValue(*big.Int)
	SetData([]byte)
	SetRcvAddr([]byte)
	SetSndAddr([]byte)
	Size() int

	CheckIntegrity() error
}

// TxStatus is the status of a transaction
type TxStatus string

const (
	// TxStatusPending = received and maybe executed on source shard, but not on destination shard
	TxStatusPending TxStatus = "pending"
	// TxStatusSuccess = received and executed
	TxStatusSuccess TxStatus = "success"
	// TxStatusFail = received and executed with error
	TxStatusFail TxStatus = "fail"
	// TxStatusInvalid = considered invalid
	TxStatusInvalid TxStatus = "invalid"
	// TxStatusRewardReverted represents the identifier for a reverted reward transaction
	TxStatusRewardReverted TxStatus = "reward-reverted"
)

// String returns the string representation of the status
func (tx TxStatus) String() string {
	return string(tx)
}

// CallType specifies the type of SC invocation (in terms of asynchronicity)
type CallType int

const (
	// DirectCall means that the call is an explicit SC invocation originating from a user Transaction
	DirectCall CallType = iota

	// AsynchronousCall means that the invocation was performed from within
	// another SmartContract from another Shard, using asyncCall
	AsynchronousCall

	// AsynchronousCallBack means that an AsynchronousCall was performed
	// previously, and now the control returns to the caller SmartContract's callBack method
	AsynchronousCallBack

	// ESDTTransferAndExecute means that there is a smart contract execution after the ESDT transfer
	// this is needed in order to skip the check whether a contract is payable or not
	ESDTTransferAndExecute

	// ExecOnDestByCaller means that the call is an invocation of a built in function / smart contract from
	// another smart contract but the caller is from the previous caller
	ExecOnDestByCaller
)

// AddressBytesLen represents the number of bytes of an address
const AddressBytesLen = 32

// MinAllowedDeltaToFinal is the minimum value between nonces allowed when checking finality on a shard
const MinAllowedDeltaToFinal = 1

//var log = logger.GetOrCreate("elrond-sdk-erdgo/core")
//
//// AddressPublicKeyConverter represents the default address public key converter
//var AddressPublicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(AddressBytesLen, log)

type address struct {
	bytes []byte
}

// NewAddressFromBytes returns a new address from provided bytes
func NewAddressFromBytes(bytes []byte) *address {
	addr := &address{
		bytes: make([]byte, len(bytes)),
	}
	copy(addr.bytes, bytes)

	return addr
}

//// NewAddressFromBech32String returns a new address from provided bech32 string
//func NewAddressFromBech32String(bech32 string) (*address, error) {
//	buff, err := AddressPublicKeyConverter.Decode(bech32)
//	if err != nil {
//		return nil, err
//	}
//
//	return &address{
//		bytes: buff,
//	}, err
//}

//// AddressAsBech32String returns the address as a bech32 string
//func (a *address) AddressAsBech32String() string {
//	return AddressPublicKeyConverter.Encode(a.bytes)
//}

// AddressBytes returns the raw address' bytes
func (a *address) AddressBytes() []byte {
	return a.bytes
}

// AddressSlice will convert the provided buffer to its [32]byte representation
func (a *address) AddressSlice() [32]byte {
	var result [32]byte
	copy(result[:], a.bytes)

	return result
}

// IsValid returns true if the contained address is valid
func (a *address) IsValid() bool {
	return len(a.bytes) == AddressBytesLen
}

// IsInterfaceNil returns true if there is no value under the interface
func (a *address) IsInterfaceNil() bool {
	return a == nil
}

// ReturnDataKind specifies how to interpret VMOutputs's return data.
// More specifically, how to interpret returned data's first item.
type ReturnDataKind int

const (
	// AsBigInt to interpret as big int
	AsBigInt ReturnDataKind = 1 << iota
	// AsBigIntString to interpret as big int string
	AsBigIntString
	// AsString to interpret as string
	AsString
	// AsHex to interpret as hex
	AsHex
)

// VMOutputApi is a wrapper over the vmcommon's VMOutput
type VMOutputApi struct {
	ReturnData      [][]byte                     `json:"returnData"`
	ReturnCode      string                       `json:"returnCode"`
	ReturnMessage   string                       `json:"returnMessage"`
	GasRemaining    uint64                       `json:"gasRemaining"`
	GasRefund       *big.Int                     `json:"gasRefund"`
	OutputAccounts  map[string]*OutputAccountApi `json:"outputAccounts"`
	DeletedAccounts [][]byte                     `json:"deletedAccounts"`
	TouchedAccounts [][]byte                     `json:"touchedAccounts"`
	Logs            []*LogEntryApi               `json:"logs"`
}

// StorageUpdateApi is a wrapper over vmcommon's StorageUpdate
type StorageUpdateApi struct {
	Offset []byte `json:"offset"`
	Data   []byte `json:"data"`
}

// OutputAccountApi is a wrapper over vmcommon's OutputAccount
type OutputAccountApi struct {
	Address         string                       `json:"address"`
	Nonce           uint64                       `json:"nonce"`
	Balance         *big.Int                     `json:"balance"`
	BalanceDelta    *big.Int                     `json:"balanceDelta"`
	StorageUpdates  map[string]*StorageUpdateApi `json:"storageUpdates"`
	Code            []byte                       `json:"code"`
	CodeMetadata    []byte                       `json:"codeMetaData"`
	OutputTransfers []OutputTransferApi          `json:"outputTransfers"`
	CallType        CallType                     `json:"callType"`
}

// OutputTransferApi is a wrapper over vmcommon's OutputTransfer
type OutputTransferApi struct {
	Value         *big.Int `json:"value"`
	GasLimit      uint64   `json:"gasLimit"`
	Data          []byte   `json:"data"`
	CallType      CallType `json:"callType"`
	SenderAddress string   `json:"senderAddress"`
}

// LogEntryApi is a wrapper over vmcommon's LogEntry
type LogEntryApi struct {
	Identifier []byte   `json:"identifier"`
	Address    string   `json:"address"`
	Topics     [][]byte `json:"topics"`
	Data       []byte   `json:"data"`
}

// GetFirstReturnData is a helper function that returns the first ReturnData of VMOutput, interpreted as specified.
func (vmOutput *VMOutputApi) GetFirstReturnData(asType ReturnDataKind) (interface{}, error) {
	if len(vmOutput.ReturnData) == 0 {
		return nil, fmt.Errorf("no return data")
	}

	returnData := vmOutput.ReturnData[0]

	switch asType {
	case AsBigInt:
		return big.NewInt(0).SetBytes(returnData), nil
	case AsBigIntString:
		return big.NewInt(0).SetBytes(returnData).String(), nil
	case AsString:
		return string(returnData), nil
	case AsHex:
		return hex.EncodeToString(returnData), nil
	}

	return nil, fmt.Errorf("can't interpret return data")
}

const (
	nodeGetNodeStatusEndpoint      = "node/status"
	nodeRawBlockByHashEndpoint     = "internal/raw/block/by-hash/%s"
	nodeRawBlockByNonceEndpoint    = "internal/raw/block/by-nonce/%d"
	nodeRawMiniBlockByHashEndpoint = "internal/raw/miniblock/by-hash/%s/epoch/%d"
)

// NewNodeEndpointProvider returns a new instance of a nodeEndpointProvider
func NewNodeEndpointProvider() *nodeEndpointProvider {
	return &nodeEndpointProvider{}
}

// GetNodeStatus returns the node status endpoint
func (node *nodeEndpointProvider) GetNodeStatus(_ uint32) string {
	return nodeGetNodeStatusEndpoint
}

// ShouldCheckShardIDForNodeStatus returns true as some extra check will need to be done when requesting from an observer
func (node *nodeEndpointProvider) ShouldCheckShardIDForNodeStatus() bool {
	return true
}

// GetRawBlockByHash returns the raw block by hash endpoint
func (node *nodeEndpointProvider) GetRawBlockByHash(_ uint32, hexHash string) string {
	return fmt.Sprintf(nodeRawBlockByHashEndpoint, hexHash)
}

// GetRawBlockByNonce returns the raw block by nonce endpoint
func (node *nodeEndpointProvider) GetRawBlockByNonce(_ uint32, nonce int64) string {
	return fmt.Sprintf(nodeRawBlockByNonceEndpoint, nonce)
}

// GetRawMiniBlockByHash returns the raw miniblock by hash endpoint
func (node *nodeEndpointProvider) GetRawMiniBlockByHash(_ uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(nodeRawMiniBlockByHashEndpoint, hexHash, epoch)
}

// GetRestAPIEntityType returns the observer node constant
func (node *nodeEndpointProvider) GetRestAPIEntityType() RestAPIEntityType {
	return ObserverNode
}

// IsInterfaceNil returns true if there is no value under the interface
func (node *nodeEndpointProvider) IsInterfaceNil() bool {
	return node == nil
}

// nodeEndpointProvider is suitable to work with an Elrond node (observer)
type nodeEndpointProvider struct {
	*baseEndpointProvider
}

const (
	networkConfig            = "network/config"
	networkEconomics         = "network/economics"
	ratingsConfig            = "network/ratings"
	enableEpochsConfig       = "network/enable-epochs"
	account                  = "address/%s"
	costTransaction          = "transaction/cost"
	sendTransaction          = "transaction/send"
	sendMultipleTransactions = "transaction/send-multiple"
	transactionStatus        = "transaction/%s/status"
	transactionInfo          = "transaction/%s"
	hyperBlockByNonce        = "hyperblock/by-nonce/%d"
	hyperBlockByHash         = "hyperblock/by-hash/%s"
	vmValues                 = "vm-values/query"
	genesisNodesConfig       = "network/genesis-nodes"
	rawStartOfEpochMetaBlock = "internal/raw/startofepoch/metablock/by-epoch/%d"
)

type baseEndpointProvider struct{}

// GetNetworkConfig returns the network config endpoint
func (base *baseEndpointProvider) GetNetworkConfig() string {
	return networkConfig
}

// GetNetworkEconomics returns the network economics endpoint
func (base *baseEndpointProvider) GetNetworkEconomics() string {
	return networkEconomics
}

// GetRatingsConfig returns the ratings config endpoint
func (base *baseEndpointProvider) GetRatingsConfig() string {
	return ratingsConfig
}

// GetEnableEpochsConfig returns the enable epochs config endpoint
func (base *baseEndpointProvider) GetEnableEpochsConfig() string {
	return enableEpochsConfig
}

// GetAccount returns the account endpoint
func (base *baseEndpointProvider) GetAccount(addressAsBech32 string) string {
	return fmt.Sprintf(account, addressAsBech32)
}

// GetCostTransaction returns the transaction cost endpoint
func (base *baseEndpointProvider) GetCostTransaction() string {
	return costTransaction
}

// GetSendTransaction returns the send transaction endpoint
func (base *baseEndpointProvider) GetSendTransaction() string {
	return sendTransaction
}

// GetSendMultipleTransactions returns the send multiple transactions endpoint
func (base *baseEndpointProvider) GetSendMultipleTransactions() string {
	return sendMultipleTransactions
}

// GetTransactionStatus returns the transaction status endpoint
func (base *baseEndpointProvider) GetTransactionStatus(hexHash string) string {
	return fmt.Sprintf(transactionStatus, hexHash)
}

// GetTransactionInfo returns the transaction info endpoint
func (base *baseEndpointProvider) GetTransactionInfo(hexHash string) string {
	return fmt.Sprintf(transactionInfo, hexHash)
}

// GetHyperBlockByNonce returns the hyper block by nonce endpoint
func (base *baseEndpointProvider) GetHyperBlockByNonce(nonce int64) string {
	return fmt.Sprintf(hyperBlockByNonce, nonce)
}

// GetHyperBlockByHash returns the hyper block by hash endpoint
func (base *baseEndpointProvider) GetHyperBlockByHash(hexHash string) string {
	return fmt.Sprintf(hyperBlockByHash, hexHash)
}

// GetVmValues returns the VM values endpoint
func (base *baseEndpointProvider) GetVmValues() string {
	return vmValues
}

// GetGenesisNodesConfig returns the genesis nodes config endpoint
func (base *baseEndpointProvider) GetGenesisNodesConfig() string {
	return genesisNodesConfig
}

// GetRawStartOfEpochMetaBlock returns the raw start of epoch metablock endpoint
func (base *baseEndpointProvider) GetRawStartOfEpochMetaBlock(epoch uint32) string {
	return fmt.Sprintf(rawStartOfEpochMetaBlock, epoch)
}

const (
	proxyGetNodeStatus      = "network/status/%d"
	proxyRawBlockByHash     = "internal/%d/raw/block/by-hash/%s"
	proxyRawBlockByNonce    = "internal/%d/raw/block/by-nonce/%d"
	proxyRawMiniBlockByHash = "internal/%d/raw/miniblock/by-hash/%s/epoch/%d"
)

// proxyEndpointProvider is suitable to work with an Elrond Proxy
type proxyEndpointProvider struct {
	*baseEndpointProvider
}

// NewProxyEndpointProvider returns a new instance of a proxyEndpointProvider
func NewProxyEndpointProvider() *proxyEndpointProvider {
	return &proxyEndpointProvider{}
}

// GetNodeStatus returns the node status endpoint
func (proxy *proxyEndpointProvider) GetNodeStatus(shardID uint32) string {
	return fmt.Sprintf(proxyGetNodeStatus, shardID)
}

// ShouldCheckShardIDForNodeStatus returns false as the proxy will ensure the correct shard dispatching of the request
func (proxy *proxyEndpointProvider) ShouldCheckShardIDForNodeStatus() bool {
	return false
}

// GetRawBlockByHash returns the raw block by hash endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByHash(shardID uint32, hexHash string) string {
	return fmt.Sprintf(proxyRawBlockByHash, shardID, hexHash)
}

// GetRawBlockByNonce returns the raw block by nonce endpoint
func (proxy *proxyEndpointProvider) GetRawBlockByNonce(shardID uint32, nonce int64) string {
	return fmt.Sprintf(proxyRawBlockByNonce, shardID, nonce)
}

// GetRawMiniBlockByHash returns the raw miniblock by hash endpoint
func (proxy *proxyEndpointProvider) GetRawMiniBlockByHash(shardID uint32, hexHash string, epoch uint32) string {
	return fmt.Sprintf(proxyRawMiniBlockByHash, shardID, hexHash, epoch)
}

// GetRestAPIEntityType returns the proxy constant
func (proxy *proxyEndpointProvider) GetRestAPIEntityType() RestAPIEntityType {
	return Proxy
}

// IsInterfaceNil returns true if there is no value under the interface
func (proxy *proxyEndpointProvider) IsInterfaceNil() bool {
	return proxy == nil
}

// CreateEndpointProvider creates a new instance of EndpointProvider
func CreateEndpointProvider(entityType RestAPIEntityType) (EndpointProvider, error) {
	switch entityType {
	case ObserverNode:
		return NewNodeEndpointProvider(), nil
	case Proxy:
		return NewProxyEndpointProvider(), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownRestAPIEntityType, entityType)
	}
}

type proxy interface {
	GetNetworkStatus(ctx context.Context, shardID uint32) (*NetworkStatus, error)
	GetRestAPIEntityType() RestAPIEntityType
	IsInterfaceNil() bool
}

type nonces struct {
	current  int64
	highest  int64
	probable int64
}

type nodeFinalityProvider struct {
	proxy proxy
}

// ErrInvalidNonceCrossCheckValueFormat signals that an invalid nonce cross-check value has been provided
var ErrInvalidNonceCrossCheckValueFormat = errors.New("invalid nonce cross check value format")

// ErrInvalidAllowedDeltaToFinal signals that an invalid allowed delta to final value has been provided
//var ErrInvalidAllowedDeltaToFinal = errors.New("invalid allowed delta to final value")

// ErrNilProxy signals that a nil proxy has been provided
var ErrNilProxy = errors.New("nil proxy")

// ErrNodeNotStarted signals that the node is not started
var ErrNodeNotStarted = errors.New("node not started")

// ErrUnknownRestAPIEntityType signals that an unknown REST API entity type has been provided
var ErrUnknownRestAPIEntityType = errors.New("unknown REST API entity type")

// NewNodeFinalityProvider creates a new instance of type nodeFinalityProvider
func NewNodeFinalityProvider(proxy proxy) (*nodeFinalityProvider, error) {
	if IfNil(proxy) {
		return nil, ErrNilProxy
	}

	return &nodeFinalityProvider{
		proxy: proxy,
	}, nil
}

// CheckShardFinalization will query the proxy and check if the target shard ID has a current nonce close to the highest nonce
// nonce <= highest_nonce + maxNoncesDelta
// it also checks the probable nonce to determine (with high degree of precision) if the node is syncing:
// nonce + maxNoncesDelta < probable_nonce
func (provider *nodeFinalityProvider) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta int64) error {
	if maxNoncesDelta < MinAllowedDeltaToFinal {
		return fmt.Errorf("%w, provided: %d, minimum: %d", ErrInvalidAllowedDeltaToFinal, maxNoncesDelta, MinAllowedDeltaToFinal)
	}

	result, err := provider.getNonces(ctx, targetShardID)
	if err != nil {
		return err
	}

	if result.current+maxNoncesDelta < result.probable {
		return fmt.Errorf("shardID %d is syncing, probable nonce is %d, current nonce is %d, max delta: %d",
			targetShardID, result.probable, result.current, maxNoncesDelta)
	}
	if result.current <= result.highest+maxNoncesDelta {
		log.Trace("nodeFinalityProvider.CheckShardFinalization - shard is in sync",
			"shardID", targetShardID, "highest nonce", result.highest, "probable nonce", result.probable,
			"current nonce", result.current, "max delta", maxNoncesDelta)
		return nil
	}

	return fmt.Errorf("shardID %d is stuck, highest nonce is %d, current nonce is %d, max delta: %d",
		targetShardID, result.highest, result.current, maxNoncesDelta)
}

func (provider *nodeFinalityProvider) getNonces(ctx context.Context, targetShardID uint32) (nonces, error) {
	networkStatusShard, err := provider.proxy.GetNetworkStatus(ctx, targetShardID)
	if err != nil {
		return nonces{}, err
	}

	result := nonces{
		current:  networkStatusShard.Nonce,
		highest:  networkStatusShard.HighestNonce,
		probable: networkStatusShard.ProbableHighestNonce,
	}

	isEmpty := result.current == 0 && result.highest == 0 && result.probable == 0
	if isEmpty {
		return nonces{}, ErrNodeNotStarted
	}

	return result, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *nodeFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}

type proxyFinalityProvider struct {
	proxy proxy
}

// NewProxyFinalityProvider creates a new instance of type proxyFinalityProvider
func NewProxyFinalityProvider(proxy proxy) (*proxyFinalityProvider, error) {
	if IfNil(proxy) {
		return nil, ErrNilProxy
	}

	return &proxyFinalityProvider{
		proxy: proxy,
	}, nil
}

// CheckShardFinalization will query the proxy and check if the target shard ID has a current nonce close to the cross
// check nonce from the metachain
// nonce(target shard ID) <= nonce(target shard ID notarized by meta) + maxNoncesDelta
func (provider *proxyFinalityProvider) CheckShardFinalization(ctx context.Context, targetShardID uint32, maxNoncesDelta int64) error {
	if maxNoncesDelta < MinAllowedDeltaToFinal {
		return fmt.Errorf("%w, provided: %d, minimum: %d", ErrInvalidAllowedDeltaToFinal, maxNoncesDelta, MinAllowedDeltaToFinal)
	}
	if targetShardID == MetachainShardId {
		// we consider this final since the minAllowedDeltaToFinal is 1
		return nil
	}

	nonceFromMeta, nonceFromShard, err := provider.getNoncesFromMetaAndShard(ctx, targetShardID)
	if err != nil {
		return err
	}

	if nonceFromShard < nonceFromMeta {
		return fmt.Errorf("shardID %d is syncing, meta cross check nonce is %d, current nonce is %d, max delta: %d",
			targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
	}
	if nonceFromShard <= nonceFromMeta+maxNoncesDelta {
		log.Trace("proxyFinalityProvider.CheckShardFinalization - shard is in sync",
			"shardID", targetShardID, "meta cross check nonce", nonceFromMeta,
			"current nonce", nonceFromShard, "max delta", maxNoncesDelta)
		return nil
	}

	return fmt.Errorf("shardID %d is stuck, meta cross check nonce is %d, current nonce is %d, max delta: %d",
		targetShardID, nonceFromMeta, nonceFromShard, maxNoncesDelta)
}

func (provider *proxyFinalityProvider) getNoncesFromMetaAndShard(ctx context.Context, targetShardID uint32) (int64, int64, error) {
	networkStatusMeta, err := provider.proxy.GetNetworkStatus(ctx, MetachainShardId)
	if err != nil {
		return 0, 0, err
	}

	crossCheckValue := networkStatusMeta.CrossCheckBlockHeight
	nonceFromMeta, err := extractNonceOfShardID(crossCheckValue, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	networkStatusShard, err := provider.proxy.GetNetworkStatus(ctx, targetShardID)
	if err != nil {
		return 0, 0, err
	}

	nonceFromShard := networkStatusShard.Nonce

	return nonceFromMeta, nonceFromShard, nil
}

func extractNonceOfShardID(crossCheckValue string, shardID uint32) (int64, error) {
	// the value will come in this format: "0: 9169897, 1: 9166353, 2: 9170524, "
	if len(crossCheckValue) == 0 {
		return 0, fmt.Errorf("%w: empty value, maybe bad observer version", ErrInvalidNonceCrossCheckValueFormat)
	}
	shardsData := strings.Split(crossCheckValue, ",")
	shardIdAsString := fmt.Sprintf("%d", shardID)

	for _, shardData := range shardsData {
		shardNonce := strings.Split(shardData, ":")
		if len(shardNonce) != 2 {
			continue
		}

		shardNonce[0] = strings.TrimSpace(shardNonce[0])
		shardNonce[1] = strings.TrimSpace(shardNonce[1])
		if shardNonce[0] != shardIdAsString {
			continue
		}

		val, ok := big.NewInt(0).SetString(shardNonce[1], 10)
		if !ok {
			return 0, fmt.Errorf("%w: %s is not a valid number as found in this response: %s",
				ErrInvalidNonceCrossCheckValueFormat, shardNonce[1], crossCheckValue)
		}

		return val.Int64(), nil
	}

	return 0, fmt.Errorf("%w: value not found for shard %d from this response: %s",
		ErrInvalidNonceCrossCheckValueFormat, shardID, crossCheckValue)
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *proxyFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}

// CreateFinalityProvider creates a new instance of FinalityProvider
func CreateFinalityProvider(proxy proxy, checkFinality bool) (FinalityProvider, error) {
	if !checkFinality {
		return NewDisabledFinalityProvider(), nil
	}

	switch proxy.GetRestAPIEntityType() {
	case ObserverNode:
		return NewNodeFinalityProvider(proxy)
	case Proxy:
		return NewProxyFinalityProvider(proxy)
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownRestAPIEntityType, proxy.GetRestAPIEntityType())
	}
}

type disabledFinalityProvider struct {
}

// NewDisabledFinalityProvider returns a new instance of type disabledFinalityProvider
func NewDisabledFinalityProvider() *disabledFinalityProvider {
	return &disabledFinalityProvider{}
}

// CheckShardFinalization will always return nil
func (provider *disabledFinalityProvider) CheckShardFinalization(_ context.Context, _ uint32, _ int64) error {
	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (provider *disabledFinalityProvider) IsInterfaceNil() bool {
	return provider == nil
}
