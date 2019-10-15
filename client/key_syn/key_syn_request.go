package key_syn

import (
	"bytes"
	"encoding/json"
	"errors"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/util"
	"time"
)

type Request struct {
	Header      string
	ReqCode     []byte
	RequestBody []byte
}

//封包

// reqBody --请求体
// reqCode -- 请求命令代码

// []byte  ---封装好的请求报文

func PacketReq(reqBody interface{}, reqCode int) ([]byte, error) {
	t := time.Now()

	if reqBody == nil && reqCode == 0 {
		err := errorx.New(errors.New("请求参数不能为空！"))
		return nil, err
	}

	logger.Logger.Info("proxy PacketReq begin")
	//bytes
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}

	ReqCode, err := util.IntToBytes(reqCode)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	// resp success
	var req = Request{
		ConstHeader,
		ReqCode,
		reqBodyData,
	}
	logger.Logger.Debug("req:", req)

	var totalBuffer bytes.Buffer

	//总长度 = 报文头长度 + 请求命令代码长度 +请求体长度
	totalLength := len(req.Header) + ReqCodeLength + len(req.RequestBody)
	logger.Logger.Debugf("totalLength:%d", totalLength)
	//请求报文:报文头长度,报文头,请求命令代码,请求应体

	buff, err := util.IntToBytes(totalLength)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff)
	totalBuffer.WriteString(req.Header)
	totalBuffer.Write(req.ReqCode)
	totalBuffer.Write(req.RequestBody)
	logger.Logger.Debugf("totalBuffer(base64):%s", totalBuffer.Bytes())
	logger.Logger.Info("proxy PacketReq end totalTime:", time.Since(t))
	return totalBuffer.Bytes(), nil
}
