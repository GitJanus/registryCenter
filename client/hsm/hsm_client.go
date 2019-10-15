//Package hsm send request to hsm client
package hsm

import (
	"bufio"
	"fmt"
	"net"
	"registryCenter/conf"
	"registryCenter/logger"
	"strings"
	"time"
)

type Hsm struct {
	HsmType string //密码机类型(SJJ1507.....)
	HsmIp   string
	HsmPort int
}

func NewHsm(hsmIp string, hsmPort int) Hsm {
	return Hsm{
		HsmIp:   hsmIp,
		HsmPort: hsmPort,
	}
}

//发送请求给密码机

// reqByte --请求体
// reqCode -- 请求命令代码

// []byte ---响应结果：包含响应结果和具体响应结果数据(resultCode,respBody)
func (hsm Hsm) SendHsmRequest(reqByte []byte, reqCode int) ([]byte, error) {
	t := time.Now()
	logger.Logger.Infof("SendHsmRequest begin")
	logger.Logger.Debugf("SendHsmRequest-base64reqByte:%s,reqCode:%d", reqByte, reqCode)

	//连接密码机
	var conn net.Conn
	var err error

	// timeout repeat try conn to it 3 times
	for i := 0; i < 3; i++ {

		conn, err = establishConn(hsm.HsmIp, hsm.HsmPort)
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
	allReqByte, err1 := PacketReq(reqByte, reqCode)
	if err1 != nil {
		logger.Logger.Error(err1)
		return nil, err1
	}
	logger.Logger.Debug("allReqByte(base64):", allReqByte)
	//send request
	conn.Write(allReqByte)

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
	logger.Logger.Debug("allresp(base64):", resp)
	//处理响应数据
	var respBody []byte
	if errscan := scanner.Err(); errscan != nil {

		logger.Logger.Error("无效数据包", err)
		return nil, errscan
	} else {

		if resp != nil {
			//去掉响应命令，包含响应结果和具体响应结果数据(resultCode,respBody)
			respBody = resp[RespCodeLength:]
			logger.Logger.Debugf("respBody(base64):%s", respBody)
		}
	}
	logger.Logger.Info("SendHsmRequest end totalTime:", time.Since(t))
	return respBody, nil

}

//构建连接

//ip--hsm ip
//port--hsm port

//net.Conn---hsm conn

func establishConn(ip string, port int) (net.Conn, error) {

	var second = conf.ConfigData.HsmServer.TimeOut
	timeout := time.Duration(second) * time.Second

	logger.Logger.Debugf("ip:%s,port:%d,timeOut:%d", ip, port, second)

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	return conn, nil
}
