package cert

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"registryCenter/conf"
	"registryCenter/logger"
	"registryCenter/logger/errorx"
	"registryCenter/model"
	"time"
)

func InitCert(certificate conf.CertificateConfig) error {
	t := time.Now()
	logger.Logger.Info("InitCert begin")
	logger.Logger.Debug("InitCert-certificate:", certificate)
	defer logger.Logger.Info("InitCert end totalTime:", time.Since(t))
	var cert model.Certificate
	var err error

	/*---服务端证书读取---begin---*/
	//读取证书
	var serverCert = certificate.ServerCert
	var serverPass = certificate.ServerPass
	serverBytes, err := ioutil.ReadFile(serverCert)
	if err != nil {
		logger.Logger.Error(err)
		return errorx.New(err)
	}
	//解析证书
	_, x509Cert, err := pkcs12.Decode(serverBytes, serverPass)
	if err != nil {
		logger.Logger.Error(err)
		return errorx.New(err)
	}

	//组装证书
	cert = TransformData(x509Cert, serverBytes, serverPass, "server", "use")

	//查询数据库
	dbCert, _ := model.GetServerCertByStatus("use")
	if dbCert != nil {
		logger.Logger.Debug("cert.SN:", cert.SN)
		logger.Logger.Debug("dbCert.SN:", dbCert.SN)

		if cert.SN != dbCert.SN {
			//将数据库中的记录状态更为 revoke,注销此证书
			dbCert.Status = "revoke"
			if err := model.UpdateCertStatus(*dbCert); err != nil {
				logger.Logger.Error(err)
			}
			//添加新的证书信息
			if err := model.InsertCert(&cert); err != nil {
				logger.Logger.Error(err)
			}
		}
	} else {
		if err := model.InsertCert(&cert); err != nil {
			logger.Logger.Error(err)
			return errorx.New(err)
		}
	}
	/*---服务端证书读取---end---*/

	/*---客户端证书读取---begin---*/
	//读取证书
	var clientCert = certificate.ClientCert
	clientCertCertPEMBlock, err := ioutil.ReadFile(clientCert)
	if err != nil {
		logger.Logger.Error(err)
		return errorx.New(err)
	}
	//解析客户端证书
	x509Cert, err = ParseCert(clientCertCertPEMBlock)
	if err != nil {
		logger.Logger.Error(err)
		return errorx.New(err)
	}
	cert = TransformData(x509Cert, clientCertCertPEMBlock, "", "client", "use")
	//添加新的证书信息
	if err := model.InsertCert(&cert); err != nil {
		logger.Logger.Error(err)
		return errorx.New(err)
	}
	return nil
}

func ParseCert(certBytes []byte) (*x509.Certificate, error) {
	t := time.Now()
	logger.Logger.Info("ParseCert begin")
	logger.Logger.Debug("ParseCert-certBytes:", certBytes)
	defer logger.Logger.Info("ParseCert end totalTime:", time.Since(t))
	//获取下一个pem格式证书数据 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
	certDERBlock, _ := pem.Decode(certBytes)
	var x509Cert *x509.Certificate
	var err error
	if certDERBlock == nil {
		x509Cert, err = x509.ParseCertificate(certBytes) //原内容为二进制编码
	} else {
		x509Cert, err = x509.ParseCertificate(certDERBlock.Bytes) //原内容为base64编码
	}
	if err != nil {
		logger.Logger.Error("x509证书解析失败")
		return nil, errorx.New(errors.New("x509证书解析失败"))
	}
	if x509Cert.NotAfter.Unix() < time.Now().Unix() {
		return nil, errorx.New(errors.New("证书已过期，不支持导入"))
	}

	return x509Cert, errorx.New(err)
}

//数据组装
func TransformData(x509Cert *x509.Certificate, certInfo []byte, password, certType, status string) model.Certificate {

	t := time.Now()
	logger.Logger.Info("transformData begin")
	logger.Logger.Debug("transformData-x509Cert:", x509Cert)
	logger.Logger.Debug("transformData-certInfo:", certInfo)
	logger.Logger.Debug("transformData-password:", password)
	logger.Logger.Debug("transformData-certType:", certType)
	logger.Logger.Debug("transformData-status:", status)
	defer logger.Logger.Info("transformData end totalTime:", time.Since(t))
	var dbCert model.Certificate
	if x509Cert != nil {
		dbCert.SN = hex.EncodeToString(x509Cert.SerialNumber.Bytes()) //序列号
		dbCert.SignAlg = fmt.Sprint(x509Cert.SignatureAlgorithm)      //签名算法
		dbCert.Issuer = fmt.Sprint(x509Cert.Issuer)                   //颁发者
		dbCert.Subject = fmt.Sprint(x509Cert.Subject)                 //使用者
		dbCert.CertInfo = certInfo                                    //证书信息
		dbCert.NotBefore = x509Cert.NotBefore                         //生效时间
		dbCert.NotAfter = x509Cert.NotAfter                           //失效时间
		dbCert.Status = status                                        //证书状态
		dbCert.Type = certType                                        //证书类型
		dbCert.CreateTime = time.Now()                                //导入时间
		//dbCert.PrivateKey = privateKey                                 //密钥信息
		dbCert.Password = password                                     //密钥信息
		dbCert.Length = x509Cert.PublicKey.(*rsa.PublicKey).N.BitLen() //公钥长度 ，单位Bits
		//获取证书指纹
		h := sha1.New()
		h.Write(x509Cert.Raw)
		dbCert.Fingerprint = hex.EncodeToString(h.Sum(nil)) //指纹
	}
	return dbCert
}
