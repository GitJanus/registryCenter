package model

import (
	"errors"
	"registryCenter/logger/errorx"
	"time"
)

type Certificate struct {
	Id          int       `json:"id"`          //主键（自增）
	SN          string    `json:"sn"`          //证书序列号
	SignAlg     string    `json:"sign_alg"`    //签名算法
	Issuer      string    `json:"issuer"`      //颁发者
	Subject     string    `json:"subject"`     //主题
	NotBefore   time.Time `json:"not_before"`  //生效时间
	NotAfter    time.Time `json:"not_after"`   //失效时间
	Length      int       `json:"length"`      //密钥长度
	Fingerprint string    `json:"fingerprint"` //指纹
	CertInfo    []byte    `json:"cert_info"`   //证书信息
	PrivateKey  []byte    `json:"private_key"` //私钥信息
	Password    string    `json:"password"`    //私钥密码
	Status      string    `json:"status"`      // use:使用中； revoke:已注销
	Type        string    `json:"type"`        // client:客户端； server:服务端
	CreateTime  time.Time `json:"create_time"` //创建时间
}

func (Certificate) TableName() string {
	return "certificate"
}

//创建证书
func InsertCert(cert *Certificate) error {
	var err error
	if cert != nil {
		if CheckCertBySN(cert.SN) {
			err = errorx.New(errors.New("使用中的证书序列号不能重复"))
		} else {
			err = sqlDB.Create(cert).Error
		}
	} else {
		err = errorx.New(errors.New("证书不能为空"))
	}
	return err
}

//获取证书
func GetCert(param *Certificate) ([]Certificate, error) {
	var cert []Certificate
	result := sqlDB.Where(&param).Find(&cert)
	return cert, result.Error
}

//根据Id获取证书
func GetCertById(id int) (*Certificate, error) {
	var cert Certificate
	result := sqlDB.First(&cert, id)
	return &cert, result.Error
}

//根据证书状态获取服务端证书
func GetServerCertByStatus(status string) (*Certificate, error) {
	var cert Certificate
	result := sqlDB.Where("type=? AND status = ?", "server", status).First(&cert)
	return &cert, result.Error
}

//根据序列号更新证书状态
func UpdateCertStatus(cert Certificate) error {
	err := sqlDB.Model(&cert).Where("sn = ?", cert.SN).Update("status", cert.Status).Error
	return err
}

//根据id删除证书
func DelCertById(id int) error {
	var cert Certificate
	err := sqlDB.Delete(&cert, id).Error
	return err
}

func DeleteAllCert() error {
	err := sqlDB.Delete(&Certificate{}).Error
	return err
}

//检验状态为“use"的证书，sn（序列号）唯一
func CheckCertBySN(sn string) bool {
	var cert []Certificate
	sqlDB.Where("status =? AND sn=?", "use", sn).Find(&cert)
	return len(cert) > 0
}
