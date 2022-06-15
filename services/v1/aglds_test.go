package v1

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestTransfer(t *testing.T) {
	privateKeys := "46604e82b9143c1540f339ba7679aacfcbad017ad4038c974324dc5eaa1dfe62"
	fmt.Println("[]byte(privateKeys)", []byte(privateKeys))
	fmt.Println("hex.EncodeToString([]byte(privateKeys))", hex.EncodeToString([]byte(privateKeys)))
	addr := "erd1jfattve6azvfxl2ke5684kdv792c4assvzv8dhhfkrm264apym5quevft4"
	value := "1000000000000000000"
	sk, err := hex.DecodeString(privateKeys)
	if err != nil {
		fmt.Println("err:", err)
	}
	a, err := Transfer(sk, addr, value)
	fmt.Println("a", a)

}
