package egld

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ElrondNetwork/elrond-go-crypto/signing"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go-crypto/signing/ed25519/singlesig"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"github.com/ethereum/go-ethereum/log"
	"github.com/group-coldwallet/trxsign/conf"
	"io/ioutil"
	"net/http"
	"time"
)

func Getnonce(url string) (int64, error) {
	args := ArgsElrondProxy{
		ProxyURL:            url,
		Client:              nil,
		SameScState:         false,
		ShouldBeSynced:      false,
		FinalityCheck:       false,
		CacheExpirationTime: time.Minute,
		EntityType:          Proxy,
	}
	ep, err := NewElrondProxy(args)
	if err != nil {
		log.Error("error creating proxy", "error", err)
		return 0, err
	}

	// Get latest hyper block (metachain) nonce
	nonce, err := ep.GetLatestHyperBlockNonce(context.Background())
	fmt.Println("nonce:", nonce)
	if err != nil {
		log.Error("error retrieving latest block nonce", "error", err)
		return 0, err
	}
	log.Info("latest hyper block", "nonce", nonce)
	return nonce, err
}
func (ep *ElrondProxy) Getblock(nonce int64) (block *HyperBlock, err error) {
	block, errGet := ep.GetHyperBlockByNonce(context.Background(), nonce)
	if errGet != nil {
		log.Error("error retrieving hyper block", "error", err)
		return
	}
	return block, err
}

type EgldBalance struct {
	Data struct {
		Account struct {
			Address         string      `json:"address"`
			Nonce           int         `json:"nonce"`
			Balance         string      `json:"balance"`
			Username        string      `json:"username"`
			Code            string      `json:"code"`
			CodeHash        interface{} `json:"codeHash"`
			RootHash        interface{} `json:"rootHash"`
			CodeMetadata    interface{} `json:"codeMetadata"`
			DeveloperReward string      `json:"developerReward"`
			OwnerAddress    string      `json:"ownerAddress"`
		} `json:"account"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

func GetBalance(address string) (Balance *EgldBalance, err error) {

	url := fmt.Sprintf("%s/%s/%s", conf.Config.NodeUrl, "address", address)
	//fmt.Println("Node:", conf.Config.NodeUrl)
	//fmt.Println("url", url)
	log.Info("url", url)
	rep, err := Get(url)
	//fmt.Println("rep", string(rep))
	log.Info("rep:", string(rep))
	b := EgldBalance{}
	err = json.Unmarshal(rep, &b)
	//fmt.Println("Balance", b)
	log.Info("Balance:", Balance)
	if err != nil {

		fmt.Println("Umarshal failed:", err)
		return
	}

	return &b, err
}
func Get(url string) ([]byte, error) {
	// 超时时间：60秒
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("resp status code is not equal 200 ,Code=[%d]", resp.StatusCode)
	}
	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}
func Post(url string) {

}
func SignTransaction(tx *data.Transaction, privateKey []byte) (*data.Transaction, error) {
	tx.Signature = ""
	txSingleSigner := &singlesig.Ed25519Signer{}
	suite := ed25519.NewEd25519()
	keyGen := signing.NewKeyGenerator(suite)
	txSignPrivKey, err := keyGen.PrivateKeyFromByteArray(privateKey)
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(&tx)
	if err != nil {
		return nil, err
	}
	signature, err := txSingleSigner.Sign(txSignPrivKey, bytes)
	if err != nil {
		return nil, err
	}
	tx.Signature = hex.EncodeToString(signature)

	return tx, err
}

//SendTransaction
func SendTransaction(tx *data.Transaction) (string, error) {
	return "", nil
}
