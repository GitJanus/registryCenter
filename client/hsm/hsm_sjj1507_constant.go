package hsm

import "registryCenter/util"

const (
	TotalLength = 2

	RespCodeLength = 2

	ReqCodeLength = 2
)

var BaseHsmType = "SJJ1507"

//获取随机数指令
var HsmGenRandomSjj1507, _ = util.BytesToInt([]byte{0x01, 0x02})

//获取ECC秘钥公钥指令
var HsmEccKeySjj1507, _ = util.BytesToInt([]byte{0x02, 0x16})

//获取RSA秘钥公钥指令
var HsmRsaKeySjj1507, _ = util.BytesToInt([]byte{0x02, 0x01})

//获取对称秘钥信息指令
var HsmSymKeySjj1507, _ = util.BytesToInt([]byte{0x08, 0x09})

//获取全部索引指令
var HsmKeySjj1507, _ = util.BytesToInt([]byte{0x08, 0x61})

//获取设备版本信息
var HsmVersionMessageSjj1507, _ = util.BytesToInt([]byte{0x01, 0x04})
