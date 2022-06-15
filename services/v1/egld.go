package v1

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/examples"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/group-coldwallet/trxsign/conf"
	"github.com/group-coldwallet/trxsign/model"
	"github.com/group-coldwallet/trxsign/util"
	"github.com/group-coldwallet/trxsign/util/egld"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sync"
	"time"
)

type EgldService struct {
	*BaseService
	client              *util.RpcClient
	nonceCtl, noncePool sync.Map
}

func (bs *BaseService) EGLDService() *EgldService {
	cs := new(EgldService)
	cs.BaseService = bs
	// 初始化连接
	client := util.New(conf.Config.EgldCfg.NodeUrl, conf.Config.EgldCfg.User, conf.Config.EgldCfg.Password)
	cs.client = client
	cs.nonceCtl = sync.Map{}
	// 新增nonce维护池
	cs.noncePool = sync.Map{}
	return cs
}

/*
接口创建地址服务
	无需改动
*/
func (cs *EgldService) CreateAddressService(req *model.ReqCreateAddressParamsV2) (*model.RespCreateAddressParams, error) {
	if req.Count == 0 {
		req.Count = 1000
	}
	if req.BatchNo == "" {
		req.BatchNo = util.GetTimeNowStr()
	}

	var (
		result *model.RespCreateAddressParams
		err    error
	)
	if conf.Config.IsStartThread {
		result, err = cs.BaseService.multiThreadCreateAddress(req.Count, req.CoinCode, req.Mch, req.BatchNo, cs.createAddressInfo)
	} else {
		result, err = cs.BaseService.createAddress(req, cs.createAddressInfo)
	}
	if err == nil {
		log.Infof("CreateAddressService 完成，共生成 %d 个地址，准备重新加载地址", len(result.Address))
		cs.InitKeyMap()
		log.Info("重新加载地址完成")
	}
	return result, err
}

/*
离线创建地址服务，通过多线程创建
	无需改动
*/
func (cs *EgldService) MultiThreadCreateAddrService(nums int, coinName, mchId, orderId string) error {
	fmt.Println("start create cph address")
	_, err := cs.BaseService.multiThreadCreateAddress(nums, coinName, mchId, orderId, cs.createAddressInfo)
	return err
}

/*
创建地址实体方法
*/
/*
签名服务
*/
func (cs *EgldService) SignService(req *model.ReqSignParams) (interface{}, error) {
	reqData, err := json.Marshal(req.Data)
	if err != nil {
		return nil, err
	}
	var tp *model.EgldSignParams
	if err := json.Unmarshal(reqData, &tp); err != nil {
		return nil, err
	}
	if &tp == nil {
		return nil, errors.New("transfer params is null")
	}
	if tp.Sender == "" || tp.Receiver == "" || tp.Value == "" {
		return nil, fmt.Errorf("params is null,from=[%s],to=[%s],amount=[%s]", tp.Sender, tp.Receiver, tp.Value)
	}
	if tp.Nonce < 0 {
		return nil, fmt.Errorf("nonce is less 0: nonce=%d", tp.Nonce)
	}
	var gasPrice *big.Int
	if tp.GasPrice <= 0 {
		gasPrice = big.NewInt(conf.Config.EgldCfg.GasPrice)
	} else {
		gasPrice = big.NewInt(tp.GasPrice)
	}
	var gasLimit int64
	if tp.GasLimit <= 0 {
		gasLimit = conf.Config.EgldCfg.GasLimit
	} else {
		gasLimit = tp.GasLimit
	}
	nonce := tp.Nonce
	toAmount, err := decimal.NewFromString(tp.Value)
	if err != nil {
		return nil, fmt.Errorf("parse amount error,err=%v", err)
	}
	amount := toAmount.BigInt()
	//toAddress := common.HexToAddress(tp.Sender)
	log.Printf("出账金额为： %d,手续费为： %d,Nonce: %d", amount.Int64(), gasPrice.Int64()*gasLimit, nonce)
	//if strings.Compare(toAddress.String(), tp.Receiver[:]) != 0 {
	//	return nil, fmt.Errorf("to address is not equal,address1=[%s],address2=[%s]", toAddress.String(),
	//		tp.Receiver[:])
	//}
	from := tp.Sender
	privateKeys, err := cs.BaseService.addressOrPublicKeyToPrivate(from)
	hexPrivateKey, err := hex.DecodeString(privateKeys)
	//fmt.Println("hexPrivateKey", hexPrivateKey)
	//log.Info("hexPrivateKey", hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("get private key error,Err=%v", err)
	}
	signtx, err := cs.getSignaturetx(hexPrivateKey, tp)
	return signtx, err
}
func (cs *EgldService) ValidAddress(address string) error {
	if !common.IsHexAddress(address) {
		return errors.New("valid ETH address error")
	}
	return nil

}
func (cs *EgldService) TransferService(req interface{}) (interface{}, error) {
	var tp *model.EgldSignParams
	if err := cs.BaseService.parseData(req, &tp); err != nil {
		return nil, err
	}

	//var hexPrivateKey string

	privateKeys, err := cs.BaseService.addressOrPublicKeyToPrivate(tp.Sender)
	hexPrivateKey, err := hex.DecodeString(privateKeys)
	fmt.Println("hexPrivateKey", hexPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("get private key error,Err=%v", err)
	}
	tx, err := cs.Transfer(hexPrivateKey, tp.Receiver, tp.Value)
	if err != nil {
		log.Error("unable to get signature", "error", err)
	}
	return tx, err
}

func (cs *EgldService) createAddressInfo() (util.AddrInfo, error) {
	w := interactors.NewWallet()
	mnemonic, err := w.GenerateMnemonic()
	if err != nil {
		log.Error("error generating mnemonic", "error", err)
	}
	//log.Info("generated mnemonics", "mnemonics", string(mnemonic))
	index0 := uint32(0)
	privkey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, index0)
	//fmt.Println("privkey", string(privkey))
	if err != nil {
		panic(err)
	}
	var (
		addrInfo util.AddrInfo
		//address  string
	)
	// 避免priv的len不是32
	if len(privkey) != 32 {
		for true {
			privkey := w.GetPrivateKeyFromMnemonic(mnemonic, 0, index0)
			if err != nil {
				// if have some error ,cut this exe
				continue
			}

			if len(privkey) == 32 {
				break
			}
		}
	}
	if privkey == nil {
		return addrInfo, errors.New("privKey is nil ptr")
	}
	//wif := hex.EncodeToString(privkey.D.Bytes())
	address, err := w.GetAddressFromPrivateKey(privkey)
	if err != nil {
		log.Error("error getting address from private key", "error", err)
	}
	addrInfo.PrivKey = hex.EncodeToString(privkey)
	//fmt.Println("addrInfo.PrivKey", addrInfo.PrivKey)
	addrInfo.Mnemonic = string(mnemonic)
	addrInfo.Address = address.AddressAsBech32String()
	return addrInfo, nil
}

func (cs *EgldService) GetBalance(req *model.ReqGetBalanceParams) (interface{}, error) {
	Balance, err := egld.GetBalance(req.Address)
	if err != nil {
		log.Error("unable to compute balance", "error", err)
	}
	return Balance.Data.Account.Balance, err
}
func (cs *EgldService) getSignaturetx(prikey []byte, req *model.EgldSignParams) (string, error) {
	args := blockchain.ArgsElrondProxy{
		ProxyURL:            conf.Config.NodeUrl,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewElrondProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
	}
	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
	}
	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
	}
	arg := data.ArgCreateTransaction{
		Nonce:     1,
		Value:     req.Value,
		RcvAddr:   req.Receiver,
		SndAddr:   req.Sender,
		GasPrice:  uint64(req.GasPrice),
		GasLimit:  uint64(req.GasLimit),
		Signature: "",
		ChainID:   req.ChainId,
		Version:   req.Version,
	}

	tx, err := ti.ApplySignatureAndGenerateTx([]byte(prikey), arg)
	if err != nil {
		log.Error("error creating transaction", "error", err)
	}
	//ti.AddTransaction(tx)
	//hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	//if err != nil {
	//	log.Error("error sending transaction", "error", err)
	//}
	//return hashes, err
	return tx.Signature, err
}

func (cs *EgldService) getSignature(prikey []byte, req *model.EgldSignParams) (*data.Transaction, error) {
	tx := &data.Transaction{
		Nonce:     uint64(req.Nonce),
		Value:     req.Value,
		RcvAddr:   req.Receiver,
		SndAddr:   req.Sender,
		GasPrice:  uint64(req.GasPrice),
		GasLimit:  uint64(req.GasLimit),
		Data:      req.Data,
		Signature: req.Signature,
		ChainID:   req.ChainId,
		Version:   req.Version,
	}
	sign, err := egld.SignTransaction(tx, prikey)
	return sign, err
}
func (cs *EgldService) Transfer(privateKey []byte, rev, value string) ([]string, error) {

	args := blockchain.ArgsElrondProxy{
		ProxyURL:            examples.TestnetGateway,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewElrondProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
	}

	w := interactors.NewWallet()

	//privateKeys, err := hex.DecodeString(hex.EncodeToString(privateKey))
	////privateKeysb := []byte(privateKeys)
	//fmt.Println("privateKey1", privateKeys)
	if err != nil {
		log.Error("unable to load alice.pem", "error", err)
	}
	// Generate address from private key
	address, err := w.GetAddressFromPrivateKey(privateKey)
	fmt.Println("address", address.AddressAsBech32String())
	if err != nil {
		log.Error("unable to load the address from the private key", "error", err)
	}

	// netConfigs can be used multiple times (for example when sending multiple transactions) as to improve the
	// responsiveness of the system
	netConfigs, err := ep.GetNetworkConfig(context.Background())
	if err != nil {
		log.Error("unable to get the network configs", "error", err)
	}
	transactionArguments, err := ep.GetDefaultTransactionArguments(context.Background(), address, netConfigs)
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
	}

	transactionArguments.RcvAddr = rev // send to self
	fmt.Println("addr", transactionArguments.RcvAddr)
	transactionArguments.Value = value // 1EGLD

	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		log.Error("unable to prepare the transaction creation arguments", "error", err)
	}

	ti, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
		log.Error("error creating transaction interactor", "error", err)
	}

	tx, err := ti.ApplySignatureAndGenerateTx(privateKey, transactionArguments)
	if err != nil {
		log.Error("error creating transaction", "error", err)
	}
	ti.AddTransaction(tx)

	hashes, err := ti.SendTransactionsAsBunch(context.Background(), 100)
	if err != nil {
		log.Error("error sending transaction", "error", err)
	}
	log.Info("transactions sent", "hashes", hashes)
	return hashes, nil
}
