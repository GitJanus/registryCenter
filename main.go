package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/tidwall/evio"
	"log"
	"registryCenter/cert"
	"registryCenter/conf"
	"registryCenter/db"
	"registryCenter/logger"
	"registryCenter/logger/logutil"
	"registryCenter/model"
	"runtime/debug"
	"strconv"
	"time"
)

func init() {
	conf.InitConfig()
	var dbType string = conf.ConfigData.DataSource.DbType
	var url string = conf.ConfigData.DataSource.Url
	db.InitDB(dbType, url)
	model.Init()
	logger.LogFileConfig()
	var communication = conf.ConfigData.Server.Communication
	if communication != 0 { //0：明文 1：单向ssl 2：双向ssl
		err := cert.InitCert(conf.ConfigData.CertificateConfig)
		logger.Logger.Error("初始化证书失败，错误原因：", err)
	}
	//hsm_server.InitCron()
}

func main() {
	//当程序崩溃时，进入该函数，记录日志
	defer TryE()
	var port int
	var loops int
	var stdlib bool

	var serverPort int = conf.ConfigData.Server.Port
	var clusterPort int = conf.ConfigData.Server.ClusterPort

	flag.IntVar(&port, "port", serverPort, "server port")
	flag.IntVar(&loops, "loops", 2, "num loops")
	flag.BoolVar(&stdlib, "stdlib", false, "use stdlib")
	flag.Parse()

	var events evio.Events
	events.NumLoops = loops

	events.Data = func(conn evio.Conn, inputStream []byte) (out []byte, action evio.Action) {
		t := time.Now()
		logger.Logger.Info("events.Data begin")
		logger.Logger.Debugf("events.Data-inputStream:%s", inputStream)

		//read req stickyBag
		buf := new(bytes.Buffer)
		buf.Write(inputStream)
		scanner := bufio.NewScanner(buf)
		//服务请求端口
		ipFlag, err1 := strconv.Atoi(conn.LocalAddr().String()[len(conn.LocalAddr().String())-4 : len(conn.LocalAddr().String())])
		if err1 != nil {
			logger.Logger.Error("获取服务请求端口失败！")
			return
		}
		//resp data
		var resp []byte
		var err error

		scanner.Split(func(data []byte, atEOF bool) (messageLength int, messgae []byte, err error) {
			if !atEOF {
				switch ipFlag {
				//本地服务指令调用
				case serverPort:
					//调用本地服务指令解析方法
					//messageLength, messgae, err = base_service.UnpackReq(data)
				case clusterPort:
					//调用第三方服务指令解析方法
					//messageLength, messgae, err = len(data), data, nil
				}
			}
			return
		})
		if scanner.Scan() {
			switch ipFlag {
			case serverPort:
				//调用本地服务业务处理入口0
				//resp, err = base_service.DoBusiness(scanner.Bytes())
			case clusterPort:
				//调用第三方服务业务处理入口
				//resp, err = cluster_service.DoClusterBusiness(scanner.Bytes())
			}
		}
		if errscan := scanner.Err(); errscan != nil {
			logger.Logger.Error("invalid data packet")
			return
		}
		if err != nil {
			logger.Logger.Error(err)
			return
		}
		out = resp
		//out = []byte{0,14,1,3,0,0,252,19,97,114,233,199,25,187,71,102}
		logger.Logger.Debugf("out:%s", out)
		logger.Logger.Info("events.Data end totalTime:", time.Since(t))
		return
	}

	events.Serving = func(srv evio.Server) (action evio.Action) {
		logger.Logger.Infof("proxyServer server started on localhost_port %d (loops: %d)", port, srv.NumLoops)
		logger.Logger.Infof("proxyServer server started on cluster_port %d (loops: %d)", clusterPort, srv.NumLoops)
		return
	}

	log.Fatal(evio.Serve(events, fmt.Sprintf("tcp://:%d", serverPort), fmt.Sprintf("tcp://:%d", clusterPort)))

}
func TryE() {
	if errs := recover(); errs != nil {
		logger.Logger.WithFields(logutil.Fields{
			"stacktrace": string(debug.Stack()), //记录堆栈信息
		}).Error(errs) //errs记录错误信息
	}
}
