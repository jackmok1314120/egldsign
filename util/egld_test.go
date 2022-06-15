package util

import (
	"encoding/hex"
	"fmt"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/blockchain"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/builders"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/core"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/interactors"
	"github.com/ethereum/go-ethereum/log"
	"testing"
	"time"
)

func Test(t *testing.T) {
	args := blockchain.ArgsElrondProxy{
		ProxyURL:            "https://testnet-gateway.elrond.com",
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}
	ep, err := blockchain.NewElrondProxy(args)
	if err != nil {
		return
	}
	txBuilder, err := builders.NewTxBuilder(blockchain.NewTxSigner())
	if err != nil {
		return
	}
	a, err := interactors.NewTransactionInteractor(ep, txBuilder)
	if err != nil {
	}
	arg := data.ArgCreateTransaction{
		Nonce:     1,
		Value:     "100000000000000000",
		RcvAddr:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		SndAddr:   "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		GasPrice:  1000000000,
		GasLimit:  50000,
		Data:      []byte(""),
		Signature: "",
		ChainID:   "T",
		Version:   uint32(1),
		Options:   1,
	}
	privateKeys := "413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9"
	fmt.Println("[]byte(privateKeys)", []byte(privateKeys))
	fmt.Println("hex.EncodeToString([]byte(privateKeys))", hex.EncodeToString([]byte(privateKeys)))

	sk, err := hex.DecodeString(privateKeys)
	fmt.Println("sk", sk)
	tx, err := a.ApplySignatureAndGenerateTx([]byte(sk), arg)
	if err != nil {
		fmt.Println("err:", err)
		log.Error("tx", err)
	}

	fmt.Println("signture:", tx)

}
