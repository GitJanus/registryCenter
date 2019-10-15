package key_syn

import (
	"github.com/pkg/errors"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/util"
	"time"
)

type Response struct {
	Header       string
	RespCode     []byte
	ResultCode   []byte
	ResponseBody interface{}
}

type ErrorBody struct {
	ErrorCode []byte `json:"code"`
	Message   string `json:"message"`
}

//解包

//buffer -- 接收到的socket响应报文
//reqCode -- 发送给服务端的请求命令代码

//int --- 整个报文长度
//[]byte --- 整个响应报文

func UnpackResp(buffer []byte, reqCode int) (int, []byte, error) {
	t := time.Now()
	logger.Logger.Info("proxy UnpackResp begin")
	length := len(buffer)

	// respcode ,respbody
	var messagge []byte
	var messageLength int

	if length >= TotalLength+HeaderLength {

		//只处理header为tassproxy的响应报文
		if string(buffer[TotalLength:HeaderLength+TotalLength]) == ConstHeader {

			//响应报文的总长度
			totalLength, err := util.BytesToInt(buffer[:TotalLength])
			if err != nil {
				logger.Logger.Error(err)
				return 0, nil, err
			}
			logger.Logger.Debugf("totalLength:%d", totalLength)
			//有可截取的报文
			if length >= totalLength {
				//verify respCode
				var respCodeBegin = HeaderLength + TotalLength
				var respCodeEnd = HeaderLength + TotalLength + RespCodeLength
				respCode, err := util.BytesToInt(buffer[respCodeBegin:respCodeEnd])
				if err != nil {
					logger.Logger.Error(err)
					return 0, nil, err
				}
				//响应命令代码：一定是请求命令代码加1
				if reqCode+1 != respCode {
					var errnew = "wrong respcode:" + string(respCode)
					return 0, make([]byte, 0), errorx.New(errors.New(errnew))
				}

				// respcode ,respbody
				messagge = buffer[HeaderLength+TotalLength : totalLength+TotalLength]
				messageLength = totalLength - HeaderLength
			}

		}
	}

	logger.Logger.Info("proxy UnpackResp end totalTime:", time.Since(t))
	logger.Logger.Debugf("messageLength:%d,messagge(base64):%s", messageLength, messagge)
	return messageLength, messagge, nil
}
