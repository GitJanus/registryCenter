package conf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
	Id          string `yaml:"id"`
	ServiceName string `yaml:"serviceName"`
	Domain      string `yaml:"domain"`
	Ip          string `yaml:"ip"`
	Port        int    `yaml:"port"`
	TimeOut     int    `yaml:"timeOut"`
	ClusterPort int    `yaml:"clusterPort"`
	Communication  int    `yaml:"communication"`
}

type RemoteServer struct {
	KeySynServer Server `yaml:"keySynServer"`
}

type DataSource struct {
	DbType string `yaml:"dbType"`
	Url    string `yaml:"url"`
}

type Log struct {
	LastDays   int    `yaml:"lastdays"`
	SizeOfFile int    `yaml:"filesize"`
	Status     int    `yaml:"status"`
	FileName   string `yaml:"filename"`
	FilePath   string `yaml:"filepath"`
}

type FileLog struct {
	InfoLog  Log `yaml:"info"`
	ErrorLog Log `yaml:"error"`
	DebugLog Log `yaml:"debug"`
	TempLog  Log `yaml:"temp"`
}
type HsmServer struct {
	TimeOut int `yaml:"timeout"`
}

//Config   系统配置配置
type Config struct {
	Server            Server            `yaml:"server"`
	DataSource        DataSource        `yaml:"dataSource"`
	RemoteServer      RemoteServer      `yaml:"remoteServer"`
	FileLog           FileLog           `yaml:"fileLog"`
	HsmServer         HsmServer         `yaml:"hsmServer"`
	CertificateConfig CertificateConfig `yaml:"certificateConfig"`
}

type CertificateConfig struct {
	ServerCert string `yaml:"serverCert"`
	ServerPass string `yaml:"serverPass"`
	ClientCert string `yaml:"clientCert"`
}

var ConfigData Config

func InitConfig() {
	config, err := ioutil.ReadFile("./conf/config.yaml")
	if err != nil {
		fmt.Print(err)
	}
	yaml.Unmarshal(config, &ConfigData)
}
