package hsm

import (
	"bytes"
	"errors"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/util"
)

type Request struct {
	ReqCode     []byte
	RequestBody []byte
}

//封包
func PacketReq(reqByte []byte, reqCode int) ([]byte, error) {
	logger.Logger.Debugf("reqByte:%s", reqByte)
	if len(reqByte) < 1 && reqCode == 0 {
		err := errorx.New(errors.New("请求参数不能为空！"))
		return nil, err
	}

	var totalBuffer bytes.Buffer
	var reqCodeByte, err = util.IntToBytes(reqCode)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalLength := ReqCodeLength + len(reqByte)
	buff, err := util.IntToBytes(totalLength)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff)
	totalBuffer.Write(reqCodeByte)
	totalBuffer.Write(reqByte)

	return totalBuffer.Bytes(), nil
}
