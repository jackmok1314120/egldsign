package v1

import (
	"fmt"
	"github.com/group-coldwallet/trxsign/util/egld"
	"github.com/shopspring/decimal"
	"testing"
)

func TestTransfer(t *testing.T) {
	//privateKeys := "46604e82b9143c1540f339ba7679aacfcbad017ad4038c974324dc5eaa1dfe62"
	//fmt.Println("[]byte(privateKeys)", []byte(privateKeys))
	//fmt.Println("hex.EncodeToString([]byte(privateKeys))", hex.EncodeToString([]byte(privateKeys)))
	//addr := "erd1jfattve6azvfxl2ke5684kdv792c4assvzv8dhhfkrm264apym5quevft4"
	value := "1000000000000000000"
	values := "1000000000000000000000"
	//sk, err := hex.DecodeString(privateKeys)
	//if err != nil {
	//	fmt.Println("err:", err)
	//}
	//a, err := cs.Transfer(sk, addr, value)
	//fmt.Println("a", a)
	toAmount, err := decimal.NewFromString(value)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("toAmount", toAmount)
	valueAmount := toAmount.BigInt()
	fmt.Println("valueAmount", valueAmount)
	fmt.Println("valueAmount.Int64", valueAmount.Int64())
	toAmounts, err := decimal.NewFromString(values)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("toAmount", toAmounts)
	valueAmounts := toAmounts.BigInt()
	fmt.Println("valueAmounts", valueAmounts)
	fmt.Println("valueAmounts.Int64", valueAmounts.Int64())

}

func TestGetData(t *testing.T) {
	mathod := "hyperblock/by-nonce"
	nonce := int64(240138)
	block, err := GetBlocks(mathod, nonce)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("valueAmounts", block)
}

func TestGetBalance(t *testing.T) {
	addr := "erd1nwnr0vpzxkvcxnd2e4xjzwzqx7n8keu76athkunvvyc0nmjydyns7dey82"
	block, err := egld.GetBalance(addr)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("Amounts", block)
}
