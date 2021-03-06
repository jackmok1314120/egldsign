package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/group-coldwallet/trxsign/conf"
	v1 "github.com/group-coldwallet/trxsign/routers/apis/v1"
)

type Apis interface {
	CreateAddress(c *gin.Context)
	Sign(c *gin.Context)
	Transfer(c *gin.Context)
	GetBalance(c *gin.Context)
	ValidAddress(c *gin.Context)
}

func CreateApis() Apis {
	var apis Apis
	switch conf.Config.Version {
	case "v1":
		if conf.Config.CoinType == "gxc" {
			// 由于之前版本gxc已经上线，所以单独抽离出来
			apis = v1.NewGxcApi()
		} else {
			apis = v1.NewBaseApi()
		}
	default:
		//默认使用v1版本
		apis = v1.NewBaseApi()
	}
	return apis
}
