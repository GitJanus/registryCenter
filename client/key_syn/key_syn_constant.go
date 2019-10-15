package key_syn

import "registryCenter/util"

const (
	ConstHeader  = "tassproxy"
	TotalLength  = 2
	HeaderLength = 9

	ResultCodeLength = 2
	RespCodeLength   = 2

	ReqCodeLength = 2
)

var KeyNoticeReqCode, _ = util.BytesToInt([]byte{0x03, 0x27})
