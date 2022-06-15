package egld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"
)

const (
	withResultsQueryParam = "?withResults=true"
)

// ArgsElrondProxy is the DTO used in the elrond proxy constructor
type ArgsElrondProxy struct {
	ProxyURL            string
	Client              Client
	SameScState         bool
	ShouldBeSynced      bool
	FinalityCheck       bool
	AllowedDeltaToFinal int
	CacheExpirationTime time.Duration
	EntityType          RestAPIEntityType
}

// ElrondProxy implements basic functions for interacting with an Elrond Proxy
type ElrondProxy struct {
	*elrondBaseProxy
	sameScState         bool
	shouldBeSynced      bool
	finalityCheck       bool
	allowedDeltaToFinal int
	finalityProvider    FinalityProvider
}

// NewElrondProxy initializes and returns an ElrondProxy object
func NewElrondProxy(args ArgsElrondProxy) (*ElrondProxy, error) {
	err := checkArgsProxy(args)
	if err != nil {
		return nil, err
	}

	endpointProvider, err := CreateEndpointProvider(args.EntityType)
	if err != nil {
		return nil, err
	}

	clientWrapper := NewHttpClientWrapper(args.Client, args.ProxyURL)
	baseArgs := ArgsElrondBaseProxy{
		httpClientWrapper: clientWrapper,
		expirationTime:    args.CacheExpirationTime,
		endpointProvider:  endpointProvider,
	}
	baseProxy, err := newElrondBaseProxy(baseArgs)
	if err != nil {
		return nil, err
	}

	finalityProvider, err := CreateFinalityProvider(baseProxy, args.FinalityCheck)
	if err != nil {
		return nil, err
	}

	ep := &ElrondProxy{
		elrondBaseProxy:     baseProxy,
		sameScState:         args.SameScState,
		shouldBeSynced:      args.ShouldBeSynced,
		finalityCheck:       args.FinalityCheck,
		allowedDeltaToFinal: args.AllowedDeltaToFinal,
		finalityProvider:    finalityProvider,
	}

	return ep, nil
}

func checkArgsProxy(args ArgsElrondProxy) error {
	if args.FinalityCheck {
		if args.AllowedDeltaToFinal < MinAllowedDeltaToFinal {
			return fmt.Errorf("%w, provided: %d, minimum: %d",
				ErrInvalidAllowedDeltaToFinal, args.AllowedDeltaToFinal, MinAllowedDeltaToFinal)
		}
	}

	return nil
}

// ExecuteVMQuery retrieves data from existing SC trie through the use of a VM
//func (ep *ElrondProxy) ExecuteVMQuery(ctx context.Context, vmRequest *VmValueRequest) (*VmValuesResponseData, error) {
//	err := ep.checkFinalState(ctx, vmRequest.Address)
//	if err != nil {
//		return nil, err
//	}
//
//	jsonVMRequestWithOptionalParams := VmValueRequestWithOptionalParameters{
//		VmValueRequest: vmRequest,
//		SameScState:    ep.sameScState,
//		ShouldBeSynced: ep.shouldBeSynced,
//	}
//	jsonVMRequest, err := json.Marshal(jsonVMRequestWithOptionalParams)
//	if err != nil {
//		return nil, err
//	}
//
//	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetVmValues(), jsonVMRequest)
//	if err != nil || code != http.StatusOK {
//		return nil, createHTTPStatusError(code, err)
//	}
//
//	response := &ResponseVmValue{}
//	err = json.Unmarshal(buff, response)
//	if err != nil {
//		return nil, err
//	}
//	if response.Error != "" {
//		return nil, errors.New(response.Error)
//	}
//
//	return &response.Data, nil
//}

//func (ep *ElrondProxy) checkFinalState(ctx context.Context, address string) error {
//	if !ep.finalityCheck {
//		return nil
//	}
//
//	targetShardID, err := ep.GetShardOfAddress(ctx, address)
//	if err != nil {
//		return err
//	}
//
//	return ep.finalityProvider.CheckShardFinalization(ctx, targetShardID, uint64(ep.allowedDeltaToFinal))
//}

// GetNetworkEconomics retrieves the network economics from the proxy
func (Ep *ElrondProxy) GetNetworkEconomics(ctx context.Context) (*NetworkEconomics, error) {
	buff, code, err := Ep.GetHTTP(ctx, Ep.endpointProvider.GetNetworkEconomics())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &NetworkEconomicsResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Economics, nil
}

// GetDefaultTransactionArguments will prepare the transaction creation argument by querying the account's info
func (ep *ElrondProxy) GetDefaultTransactionArguments(
	ctx context.Context,
	address AddressHandler,
	networkConfigs *NetworkConfig,
) (ArgCreateTransaction, error) {
	if networkConfigs == nil {
		return ArgCreateTransaction{}, ErrNilNetworkConfigs
	}
	if IfNil(address) {
		return ArgCreateTransaction{}, ErrNilAddress
	}

	//account, err := ep.GetAccount(ctx, address)
	//if err != nil {
	//	return ArgCreateTransaction{}, err
	//}

	return ArgCreateTransaction{
		//Nonce:            account.Nonce,
		Value:     "",
		RcvAddr:   "",
		SndAddr:   address.AddressAsBech32String(),
		GasPrice:  networkConfigs.MinGasPrice,
		GasLimit:  networkConfigs.MinGasLimit,
		Data:      nil,
		Signature: "",
		ChainID:   networkConfigs.ChainID,
		Version:   networkConfigs.MinTransactionVersion,
		Options:   0,
		//AvailableBalance: account.Balance,
	}, nil
}

//GetAccount retrieves an account info from the network (nonce, balance)
//func (ep *ElrondProxy) GetAccount(ctx context.Context, address AddressHandler) (*Account, error) {
//	err := ep.checkFinalState(ctx, address.AddressAsBech32String())
//	if err != nil {
//		return nil, err
//	}
//
//	if IfNil(address) {
//		return nil, ErrNilAddress
//	}
//	if !address.IsValid() {
//		return nil, ErrInvalidAddress
//	}
//	endpoint := ep.endpointProvider.GetAccount(address.AddressAsBech32String())
//
//	buff, code, err := ep.GetHTTP(ctx, endpoint)
//	if err != nil || code != http.StatusOK {
//		return nil, createHTTPStatusError(code, err)
//	}
//
//	response := &AccountResponse{}
//	err = json.Unmarshal(buff, response)
//	if err != nil {
//		return nil, err
//	}
//	if response.Error != "" {
//		return nil, errors.New(response.Error)
//	}
//
//	return response.Data.Account, nil
//}
//func (ep *ElrondProxy) checkFinalState(ctx context.Context, address string) error {
//	if !ep.finalityCheck {
//		return nil
//	}
//
//	targetShardID, err := ep.GetShardOfAddress(ctx, address)
//	if err != nil {
//		return err
//	}
//
//	return ep.finalityProvider.CheckShardFinalization(ctx, targetShardID, uint64(ep.allowedDeltaToFinal))
//}

// SendTransaction broadcasts a transaction to the network and returns the txhash if successful
func (ep *ElrondProxy) SendTransaction(ctx context.Context, tx *Transaction) (string, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetSendTransaction(), jsonTx)
	if err != nil || code != http.StatusOK {
		return "", createHTTPStatusError(code, err)
	}

	response := &SendTransactionResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.Data.TxHash, nil
}

// SendTransactions broadcasts the provided transactions to the network and returns the txhashes if successful
func (ep *ElrondProxy) SendTransactions(ctx context.Context, txs []*Transaction) ([]string, error) {
	jsonTx, err := json.Marshal(txs)
	if err != nil {
		return nil, err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetSendMultipleTransactions(), jsonTx)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &SendTransactionsResponse{}
	fmt.Println("res:", response)
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return ep.postProcessSendMultipleTxsResult(response)
}

func (ep *ElrondProxy) postProcessSendMultipleTxsResult(response *SendTransactionsResponse) ([]string, error) {
	txHashes := make([]string, 0, len(response.Data.TxsHashes))
	indexes := make([]int, 0, len(response.Data.TxsHashes))
	for index := range response.Data.TxsHashes {
		indexes = append(indexes, index)
	}

	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	for _, idx := range indexes {
		txHashes = append(txHashes, response.Data.TxsHashes[idx])
	}

	return txHashes, nil
}

// GetTransactionStatus retrieves a transaction's status from the network
func (ep *ElrondProxy) GetTransactionStatus(ctx context.Context, hash string) (string, error) {
	endpoint := ep.endpointProvider.GetTransactionStatus(hash)
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return "", createHTTPStatusError(code, err)
	}

	response := &TransactionStatus{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return "", err
	}
	if response.Error != "" {
		return "", errors.New(response.Error)
	}

	return response.Data.Status, nil
}

// GetTransactionInfo retrieves a transaction's details from the network
func (ep *ElrondProxy) GetTransactionInfo(ctx context.Context, hash string) (*TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, false)
}

// GetTransactionInfoWithResults retrieves a transaction's details from the network with events
func (ep *ElrondProxy) GetTransactionInfoWithResults(ctx context.Context, hash string) (*TransactionInfo, error) {
	return ep.getTransactionInfo(ctx, hash, true)
}

func (ep *ElrondProxy) getTransactionInfo(ctx context.Context, hash string, withResults bool) (*TransactionInfo, error) {
	endpoint := ep.endpointProvider.GetTransactionInfo(hash)
	if withResults {
		endpoint += withResultsQueryParam
	}

	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &TransactionInfo{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response, nil
}

// RequestTransactionCost retrieves how many gas a transaction will consume
func (ep *ElrondProxy) RequestTransactionCost(ctx context.Context, tx *Transaction) (*TxCostResponseData, error) {
	jsonTx, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	buff, code, err := ep.PostHTTP(ctx, ep.endpointProvider.GetCostTransaction(), jsonTx)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &ResponseTxCost{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response.Data, nil
}

const MetachainShardId = uint32(0xFFFFFFFF)

// GetLatestHyperBlockNonce retrieves the latest hyper block (metachain) nonce from the network
func (ep *ElrondProxy) GetLatestHyperBlockNonce(ctx context.Context) (int64, error) {
	response, err := ep.GetNetworkStatus(ctx, MetachainShardId)
	fmt.Println("resp:", response)
	fmt.Println("shardid::", MetachainShardId)
	if err != nil {
		return 0, err
	}

	return response.Nonce, nil
}

// GetHyperBlockByNonce retrieves a hyper block's info by nonce from the network
func (ep *ElrondProxy) GetHyperBlockByNonce(ctx context.Context, nonce int64) (*HyperBlock, error) {
	endpoint := ep.endpointProvider.GetHyperBlockByNonce(nonce)

	return ep.getHyperBlock(ctx, endpoint)
}

// GetHyperBlockByHash retrieves a hyper block's info by hash from the network
func (ep *ElrondProxy) GetHyperBlockByHash(ctx context.Context, hash string) (*HyperBlock, error) {
	endpoint := ep.endpointProvider.GetHyperBlockByHash(hash)

	return ep.getHyperBlock(ctx, endpoint)
}

func (ep *ElrondProxy) getHyperBlock(ctx context.Context, endpoint string) (*HyperBlock, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &HyperBlockResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return &response.Data.HyperBlock, nil
}

// GetRawBlockByHash retrieves a raw block by hash from the network
func (ep *ElrondProxy) GetRawBlockByHash(ctx context.Context, shardId uint32, hash string) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawBlockByHash(shardId, hash)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawBlockByNonce retrieves a raw block by hash from the network
func (ep *ElrondProxy) GetRawBlockByNonce(ctx context.Context, shardId uint32, nonce int64) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawBlockByNonce(shardId, nonce)

	return ep.getRawBlock(ctx, endpoint)
}

// GetRawStartOfEpochMetaBlock retrieves a raw block by hash from the network
func (ep *ElrondProxy) GetRawStartOfEpochMetaBlock(ctx context.Context, epoch uint32) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawStartOfEpochMetaBlock(epoch)

	return ep.getRawBlock(ctx, endpoint)
}

func (ep *ElrondProxy) getRawBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &RawBlockRespone{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Block, nil
}

// GetRawMiniBlockByHash retrieves a raw block by hash from the network
func (ep *ElrondProxy) GetRawMiniBlockByHash(ctx context.Context, shardId uint32, hash string, epoch uint32) ([]byte, error) {
	endpoint := ep.endpointProvider.GetRawMiniBlockByHash(shardId, hash, epoch)

	return ep.getRawMiniBlock(ctx, endpoint)
}

func (ep *ElrondProxy) getRawMiniBlock(ctx context.Context, endpoint string) ([]byte, error) {
	buff, code, err := ep.GetHTTP(ctx, endpoint)
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &RawMiniBlockRespone{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.MiniBlock, nil
}

// GetNonceAtEpochStart retrieves the start of epoch nonce from hyper block (metachain)
func (ep *ElrondProxy) GetNonceAtEpochStart(ctx context.Context, shardId uint32) (uint64, error) {
	response, err := ep.GetNetworkStatus(ctx, shardId)
	if err != nil {
		return 0, err
	}

	return response.NonceAtEpochStart, nil
}

// GetRatingsConfig retrieves the ratings configuration from the proxy
func (ep *ElrondProxy) GetRatingsConfig(ctx context.Context) (*RatingsConfig, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetRatingsConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &RatingsConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetEnableEpochsConfig retrieves the ratings configuration from the proxy
func (ep *ElrondProxy) GetEnableEpochsConfig(ctx context.Context) (*EnableEpochsConfig, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetEnableEpochsConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &EnableEpochsConfigResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Config, nil
}

// GetGenesisNodesPubKeys retrieves genesis nodes configuration from proxy
func (ep *ElrondProxy) GetGenesisNodesPubKeys(ctx context.Context) (*GenesisNodes, error) {
	buff, code, err := ep.GetHTTP(ctx, ep.endpointProvider.GetGenesisNodesConfig())
	if err != nil || code != http.StatusOK {
		return nil, createHTTPStatusError(code, err)
	}

	response := &GenesisNodesResponse{}
	err = json.Unmarshal(buff, response)
	if err != nil {
		return nil, err
	}
	if response.Error != "" {
		return nil, errors.New(response.Error)
	}

	return response.Data.Nodes, nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ep *ElrondProxy) IsInterfaceNil() bool {
	return ep == nil
}
