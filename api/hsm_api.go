package api

import (
	"bytes"
	"errors"
	"fmt"
	"registryCenter/client/hsm"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/util"
	"time"
)

//获取秘钥索引 密钥类型： 0–对称密钥 1- RSA密钥 2–ECC密钥
func GetKeyIndex(hsmIp string, hsmPort int, hsmType string, indexType uint32) ([]int, error) {
	logger.Logger.Infof("hsmIp:%s", hsmIp)
	logger.Logger.Infof("hsmPort:%d", hsmPort)
	logger.Logger.Infof("hsmType:%s", hsmType)
	logger.Logger.Infof("indexType:%d", indexType)
	startTime := time.Now()
	var totalBuffer bytes.Buffer
	buff, err := util.Uint32ToBytes(indexType)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff)
	var hsmVo = hsm.NewHsm(hsmIp, hsmPort)
	//按照密码机类型选择密码机指令
	var reqCode int
	switch hsmType {
	case hsm.BaseHsmType:
		reqCode = hsm.HsmKeySjj1507
	}
	if reqCode == 0 {
		err := errorx.New(errors.New("密码机类型错误！"))
		return nil, err
	}
	resp, err := hsmVo.SendHsmRequest(totalBuffer.Bytes(), reqCode)
	var indexList []int
	if err != nil {
		logger.Logger.Error(err)
		logger.Logger.Debugf("%s方法执行时间：%s", "GetKeyIndex", time.Since(startTime))
		return indexList, err
	} else {
		//请求成功
		if resp[0:2][0] == 0 && resp[0:2][1] == 0 {

			buff, err := util.GetKeyIndexLen(resp[2:6])
			if err != nil {
				logger.Logger.Error(err)
				return nil, err
			}
			byteList := resp[6 : 6+buff]
			for i := 0; i < len(byteList); i++ {
				j := 1
				buff, err := util.ConvertToBin(int(byteList[i]))
				if err != nil {
					logger.Logger.Error(err)
					return nil, err
				}
				for _, s := range string(buff) {
					if string(s) == "1" {
						indexList = append(indexList, i*8+j)
					}
					j++
				}
			}
		} else {
			logger.Logger.Errorf("%s方法获取秘钥索引失败", "GetKeyIndex")
			err := errorx.New(errors.New("GetKeyIndex方法获取秘钥索引失败"))
			return indexList, err
		}
		logger.Logger.Debugf("%s方法执行时间：%s", "GetKeyIndex", time.Since(startTime))
		return indexList, err
	}
}

//获取ECC密钥的公钥
func GetECCPublicKey(hsmIp string, hsmPort int, hsmType string, index int, keyType uint32) ([]byte, error) {
	logger.Logger.Infof("hsmIp:%s", hsmIp)
	logger.Logger.Infof("hsmPort:%d", hsmPort)
	logger.Logger.Infof("hsmType:%s", hsmType)
	logger.Logger.Infof("index:%d", index)
	logger.Logger.Infof("keyType:%d", keyType)
	startTime := time.Now()
	var totalBuffer bytes.Buffer

	unit32Index, err := util.Uint32ToBytes(uint32(index))
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(unit32Index)
	//密钥类型（0 签名密钥or 1加密密钥 or 2密钥交换密钥）

	unit32KeyType, err := util.Uint32ToBytes(keyType)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(unit32KeyType)
	var hsmVo = hsm.NewHsm(hsmIp, hsmPort)
	//按照密码机类型选择密码机指令
	var reqCode int
	switch hsmType {
	case hsm.BaseHsmType:
		reqCode = hsm.HsmEccKeySjj1507
	}
	if reqCode == 0 {
		err := errorx.New(errors.New("密码机类型错误！"))
		return nil, err
	}
	resp, err := hsmVo.SendHsmRequest(totalBuffer.Bytes(), reqCode)
	if err != nil {
		logger.Logger.Error(err)
	}
	logger.Logger.Debugf("%s方法执行时间：%s", "GetECCPublicKey", time.Since(startTime))
	return resp, err
}

//获取RSA密钥的公钥
func GetRSAPublicKey(hsmIp string, hsmPort int, hsmType string, index int, keyType uint32) ([]byte, error) {
	logger.Logger.Infof("hsmIp:%s", hsmIp)
	logger.Logger.Infof("hsmPort:%d", hsmPort)
	logger.Logger.Infof("hsmType:%s", hsmType)
	logger.Logger.Infof("index:%d", index)
	logger.Logger.Infof("keyType:%d", keyType)
	startTime := time.Now()
	var totalBuffer bytes.Buffer
	buff, err := util.Uint32ToBytes(uint32(index))
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff)
	buff1, err := util.Uint32ToBytes(keyType)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff1)
	var hsmVo = hsm.NewHsm(hsmIp, hsmPort)

	//按照密码机类型选择密码机指令
	var reqCode int
	switch hsmType {
	case hsm.BaseHsmType:
		reqCode = hsm.HsmRsaKeySjj1507
	}
	if reqCode == 0 {
		err := errorx.New(errors.New("密码机类型错误！"))
		return nil, err
	}
	resp, err := hsmVo.SendHsmRequest(totalBuffer.Bytes(), reqCode)
	if err != nil {
		logger.Logger.Error(err)
	}
	logger.Logger.Debugf("%s方法执行时间：%s", "GetRSAPublicKey", time.Since(startTime))
	return resp, err
}

//向密码机发送随机数指令   指令正常返回false 反之 返回true
func SendRandomInstruction(hsmIp string, hsmPort int, hsmType string) (bool, error) {
	logger.Logger.Infof("hsmIp:%s", hsmIp)
	logger.Logger.Infof("hsmPort:%d", hsmPort)
	logger.Logger.Infof("hsmType:%s", hsmType)
	startTime := time.Now()
	var leng uint32 = 10
	var lengByte, err = util.Uint32ToBytes(leng)
	if err != nil {
		logger.Logger.Error(err)
		return true, err
	}
	var hsmVo = hsm.NewHsm(hsmIp, hsmPort)
	//按照密码机类型选择密码机指令
	var reqCode int
	switch hsmType {
	case hsm.BaseHsmType:
		reqCode = hsm.HsmGenRandomSjj1507
	}
	if reqCode == 0 {
		err := errorx.New(errors.New("密码机类型错误！"))
		return true, err
	}
	resp, err := hsmVo.SendHsmRequest(lengByte, reqCode)
	if err != nil {
		logger.Logger.Error(err)
		logger.Logger.Debugf("%s方法执行时间：%s", "SendRandomInstruction", time.Since(startTime))
		return true, err
	} else if resp[0:2][0] == 0 && resp[0:2][1] == 0 {
		fmt.Println(resp)
		logger.Logger.Debugf("%s方法执行时间：%s", "SendRandomInstruction", time.Since(startTime))
		return false, err
	}
	return true, err
}

func GetSymKeyInfo(hsmIp string, hsmPort int, hsmType string, index int) ([]byte, error) {
	logger.Logger.Infof("hsmIp:%s", hsmIp)
	logger.Logger.Infof("hsmPort:%d", hsmPort)
	logger.Logger.Infof("hsmType:%s", hsmType)
	logger.Logger.Infof("index:%d", index)
	startTime := time.Now()
	var totalBuffer bytes.Buffer
	buff, err := util.Uint32ToBytes(uint32(index))
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	totalBuffer.Write(buff)
	var hsmVo = hsm.NewHsm(hsmIp, hsmPort)
	//按照密码机类型选择密码机指令
	var reqCode int
	switch hsmType {
	case hsm.BaseHsmType:
		reqCode = hsm.HsmSymKeySjj1507
	}
	if reqCode == 0 {
		err := errorx.New(errors.New("密码机类型错误！"))
		return nil, err
	}
	resp, err := hsmVo.SendHsmRequest(totalBuffer.Bytes(), reqCode)
	if err != nil {
		logger.Logger.Error(err)
		return resp, err
	}
	logger.Logger.Debugf("%s方法执行时间：%s", "GetSymKeyInfo", time.Since(startTime))
	return resp, err
}
