package util

import (
	"bytes"
	"registryCenter/logger"
	"strconv"
)

func BytesToInt(bytes []byte) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	s1 := (int)(bytes[0] & 0xff)
	s2 := (int)(bytes[1] & 0xff)
	s1 <<= 8
	var result = (int)(s1 | s2)
	return result, nil

}

func IntToBytes(num int) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	var bytes2 = make([]byte, 2)
	bytes2[0] = (byte)(num >> 8 & 0xff)
	bytes2[1] = (byte)(num & 0xff)
	return bytes2, nil
}

//4B
func Uint32ToBytes(num uint32) ([]byte, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	var bytes2 = make([]byte, 4)
	bytes2[0] = (byte)(num >> 24 & 0xFF)
	bytes2[1] = (byte)(num >> 16 & 0xFF)
	bytes2[2] = (byte)(num >> 8 & 0xFF)
	bytes2[3] = (byte)(num >> 0 & 0xFF)
	return bytes2, nil
}

//4B
func BytesToUint32(bytes []byte) uint32 {
	s1 := (int)(bytes[0] & 0xFF)
	s2 := (int)(bytes[1] & 0xFF)
	s3 := (int)(bytes[2] & 0xFF)
	s4 := (int)(bytes[3] & 0xFF)
	s1 <<= 24
	s2 <<= 16
	s3 <<= 8
	var result = (uint32)(s1 | s2 | s3 | s4)
	return result
}

// 将十进制数字转化为二进制字符串
func ConvertToBin(num int) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	s := ""
	if num == 0 {
		return "00000000", nil
	}
	// num /= 2 每次循环的时候 都将num除以2  再把结果赋值给 num
	for ; num > 0; num /= 2 {
		lsb := num % 2
		// strconv.Itoa() 将数字强制性转化为字符串
		s = strconv.Itoa(lsb) + s
	}
	var buffer bytes.Buffer
	if len(s) < 8 {
		for i := 0; i < 8-len(s); i++ {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
		return buffer.String(), nil
	}
	return s, nil
}

//获取秘钥索引长度
func GetKeyIndexLen(lenByte []byte) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	var totalBuffer bytes.Buffer
	for i := 0; i < len(lenByte); i++ {
		data, err := ConvertToBin(int(lenByte[i]))
		if err == nil {
			for _, s := range data {
				totalBuffer.WriteString(string(s))
			}
		}
	}
	num, err := Str2DEC(totalBuffer.String())
	if err != nil {
		logger.Logger.Error(err)
		return 0, err
	}
	return num, nil
}

//将二进制string字符串转换成int
func Str2DEC(s string) (num int, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("捕获到的错误：", r)
			return
		}
	}()
	l := len(s)
	for i := l - 1; i >= 0; i-- {
		num += (int(s[l-i-1]) & 0xf) << uint8(i)
	}
	return
}
