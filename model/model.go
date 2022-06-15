package model

//========================create address==================//
type ReqBaseParams struct {
	OrderId    string `json:"orderId,omitempty"`
	MchId      string `json:"mchId,omitempty"`
	CoinName   string `json:"coinName,omitempty"`
	Sign       string `json:"sign"`
	CreateTime string `json:"createTime"` // time.Now().Unix()  长度为10位的时间戳
}

type ApiSignParams struct {
	ToAddress     string `json:"to_address"`
	ToAmount      string `json:"to_amount"` //精度转换后的amount
	ChangeAddress string `json:"change_address"`
	ChangeAmount  string `json:"change_amount"` //精度转换后的amount
}

type ReqCreateAddressParams struct {
	Num int `json:"num"`
	ReqBaseParams
}

type ReqCreateAddressParamsV2 struct {
	// 商户名称,如：hoo
	Mch string `json:"mch,omitempty"`
	// 币种，如：trx
	CoinCode string `json:"coinCode,omitempty"`
	// 生成地址数量
	Count int `json:"count,omitempty"`
	// 本次生成地址对应的编号，如：trx_usb_20220601001
	BatchNo string `json:"batchNo,omitempty"`
}

type RespCreateAddressParams struct {
	// 商户名称,如：hoo（直接使用入参的值）
	Mch string `json:"mch,omitempty"`
	// 币种，如：trx（直接使用入参的值）
	CoinCode string `json:"coinCode,omitempty"`
	// 本次生成地址对应的编号，如：trx_usb_20220601001（直接使用入参的值）
	BatchNo string `json:"batchNo,omitempty"`
	// 本次生成的地址列表
	Address []string `json:"address"`
}

//==========================================================//

type ReqSignParams struct {
	ReqBaseParams
	Data interface{} `json:"data"`
}

type RespSignParams struct {
	ReqBaseParams
	Result string `json:"result"`
}

//=============================================================//

//-----------------------getBalance----------------------------//
type ReqGetBalanceParams struct {
	CoinName        string      `json:"coin_name"`        // 币种主链的名字
	Address         string      `json:"address"`          //	需要获取余额的地址
	Token           string      `json:"token"`            // 	token的名字
	ContractAddress string      `json:"contract_address"` //合约地址
	Params          interface{} `json:"params"`           //特殊参数（如果有特殊参数，传入到这里面）
}

//------------------valid address-------------------
type ReqValidAddressParams struct {
	Address string `json:"address"`
}

//----------egldtransfer--------
type EgldSignParams struct {
	Version   uint32 `json:"version"`
	ChainId   string `json:"chainId"`
	Nonce     int64  `json:"nonce"`
	Value     string `json:"value"`
	Receiver  string `json:"receiver"`
	Sender    string `json:"sender"`
	GasPrice  int64  `json:"gasPrice"`
	GasLimit  int64  `json:"gasLimit"`
	Data      []byte `json:"data"`
	Signature string `json:"signature"`
}
