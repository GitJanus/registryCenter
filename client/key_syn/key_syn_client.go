package key_syn

import (
	"bufio"
	"fmt"
	"net"
	"registryCenter/conf"
	"registryCenter/logger"
	"strings"
	"time"
)

//发送请求给代理服务

// reqBody --请求体
// reqCode -- 请求命令代码

// []byte ---响应结果：包含响应结果和具体响应结果数据(resultCode,respBody)
func SendRequest(reqBody interface{}, reqCode int) ([]byte, error) {
	t := time.Now()
	logger.Logger.Info("ProxySendRequest begin")
	logger.Logger.Debug("ProxySendRequest-reqBody:", reqBody)
	logger.Logger.Debugf("ProxySendRequest-reqCode:%d", reqCode)

	var conn net.Conn
	var err error

	// timeout repeat try conn to it 3 times
	for i := 0; i < 3; i++ {

		conn, err = establishConn()

		if err != nil {
			logger.Logger.Error(err)
			if strings.Index(err.Error(), "timeout") > -1 {
				logger.Logger.Error("error times:"+string(i), err)
			} else {
				break
			}
		} else {
			break
		}

	}

	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}

	//关闭连接
	defer conn.Close()

	//封装请求
	allReqBodyByte, err1 := PacketReq(reqBody, reqCode)
	if err1 != nil {
		logger.Logger.Error(err1)
		return nil, err1
	}
	logger.Logger.Debug("allReqBodyByte(base64):", allReqBodyByte)
	//send request
	conn.Write(allReqBodyByte)

	//deal resp
	scanner := bufio.NewScanner(conn)
	scanner.Split(func(data []byte, atEOF bool) (messageLength int, message []byte, err error) {
		if !atEOF {
			messageLength, message, err = UnpackResp(data, reqCode)
		}
		return
	})
	var resp []byte
	if scanner.Scan() {
		resp = scanner.Bytes()
	}
	logger.Logger.Debug("resp(base64):", resp)

	//处理响应数据
	var respBody []byte
	if errscan := scanner.Err(); errscan != nil {

		logger.Logger.Error("无效数据包", err)
		return nil, errscan
	} else {

		//去掉响应命令，包含响应结果和具体响应结果数据(resultCode,respBody)
		if resp != nil {
			respBody = resp[RespCodeLength:]
			logger.Logger.Debugf("respBody(base64):%s", respBody)
		}
	}
	logger.Logger.Info("ProxySendRequest end totalTime:", time.Since(t))
	return respBody, nil

}

func establishConn() (net.Conn, error) {
	t := time.Now()
	logger.Logger.Info("establishConn begin")

	var port = conf.ConfigData.RemoteServer.KeySynServer.Port
	var ip = conf.ConfigData.RemoteServer.KeySynServer.Ip
	var second = conf.ConfigData.RemoteServer.KeySynServer.TimeOut

	logger.Logger.Debugf("ip:%s,port:%d,timeOut:%d", ip, port, second)

	timeout := time.Duration(second) * time.Second

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}

	logger.Logger.Info("establishConn end totalTime:", time.Since(t))
	return conn, nil
}
