package egld

import (
	"encoding/hex"
	"fmt"
	"github.com/ElrondNetwork/elrond-sdk-erdgo/data"
	"testing"
)

func Test(t *testing.T) {
	add := "erd1dj0yzmu3z2nnv463xpns8w3t99y894czz7t56xv4qdhc5we9039sm6aumm"
	GetBalance(add)
	//fmt.Println("hello")

}
func TestSign(t *testing.T) {
	arg := &data.Transaction{
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
	if err != nil {
		fmt.Println("err:", err)
	}
	a, err := SignTransaction(arg, sk)
	fmt.Println("signature:", a.Signature, err)
}
