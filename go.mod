module github.com/group-coldwallet/trxsign

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Dipper-Labs/go-sdk v1.0.3
	github.com/JFJun/arweave-go v0.0.0-20200525082925-be2aa616e219
	github.com/JFJun/trx-sign-go v1.0.3
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/ethereum/go-ethereum v1.10.8
	github.com/fbsobreira/gotron-sdk v0.0.0-20201030191254-389aec83c8f9
	github.com/gin-gonic/gin v1.7.7
	github.com/go-redis/redis/v8 v8.11.4
	github.com/kr/text v0.2.0 // indirect
	github.com/mendsley/gojwk v0.0.0-20141217222730-4d5ec6e58103
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.6.0
	go.uber.org/zap v1.19.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/ElrondNetwork/elrond-go-crypto v1.0.1
	github.com/ElrondNetwork/elrond-sdk-erdgo v1.0.22
)

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_2 v1.2.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.2.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_3 v1.3.35 => github.com/ElrondNetwork/arwen-wasm-vm v1.3.35

replace github.com/ElrondNetwork/arwen-wasm-vm/v1_4 v1.4.40 => github.com/ElrondNetwork/arwen-wasm-vm v1.4.40
