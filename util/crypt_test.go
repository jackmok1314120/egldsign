package util

import (
	"testing"
)

func TestAesEncode(t *testing.T) {
	aesKey := RandBase64Key()
	t.Logf("随机生成密钥：%s", aesKey)
	crypter, _ := AesBase64Crypt([]byte("KyAevMVhccueJq9c4FiUUsXjQJudXWHHCLDeXaCMbx4JAnTY3VgK"), aesKey, true)
	t.Logf("加密结果：%s", string(crypter))

}

func TestAesDecode(t *testing.T) {
	ciphertext := "uiV1IcdYQ+W1GHhlpfqlRi5dB6tsWLpoOcb448tMbZdCXNBx0PIQDg2gry2VVlJXUoZHOQ=="
	key := "/+t4wzKOeS2qqcLTcbTQJi/ubtPQUt1b"
	result := "KyAevMVhccueJq9c4FiUUsXjQJudXWHHCLDeXaCMbx4JAnTY3VgK"

	text, _ := AesBase64Crypt([]byte(ciphertext), []byte(key), false)
	t.Logf("解密结果：%s", string(text))
	t.Logf("解密结果对比：%t", result == string(text))
}
func TestCreatAddr(t *testing.T) {
	//res := new(model.ReqCreateAddressParamsV2)
	//res.Count = 2
	//res.CoinCode = "egld"
	//res.Mch = "hoo"
	//res.BatchNo = "egld_usb_20220601001"
	//v1.EgldService.CreateAddressService()
	//fmt.Println(v1.CreateAddressInfo())
}
