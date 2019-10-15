package hsm

import (
	"errors"
	"fmt"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/util"
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
func UnpackResp(buffer []byte, reqCode int) (int, []byte, error) {
	fmt.Println(buffer)
	length := len(buffer)

	if length < 1 {
		err := errorx.New(errors.New("请求无响应！"))
		return 0, nil, err
	}

	// respcode ,respbody
	var messagge []byte
	var messageLength int

	if length >= TotalLength {

		totalLength, err := util.BytesToInt(buffer[:TotalLength])
		if err != nil {
			logger.Logger.Error(err)
			return 0, nil, err
		}

		if length >= totalLength {
			//verify respCode
			var respCodeBegin = TotalLength
			var respCodeEnd = TotalLength + RespCodeLength
			respCode, err := util.BytesToInt(buffer[respCodeBegin:respCodeEnd])
			if err != nil {
				logger.Logger.Error(err)
				return 0, nil, err
			}

			if reqCode+1 != respCode {
				return 0, make([]byte, 0), errorx.New(errors.New("wrong respcode"))
			}
			// respcode ,respbody
			messagge = buffer[TotalLength : totalLength+TotalLength]
			messageLength = totalLength
		}
	}
	return messageLength, messagge, nil
}
