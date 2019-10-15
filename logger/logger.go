package logger

import (
	"io"
	"path"
	"registryCenter/conf"
	rotatelogs "registryCenter/logger/file-rotatelogs"
	"registryCenter/logger/lfshook"
	"registryCenter/logger/logutil"
	"time"
)

var Logger *logutil.Logger

func LogFileConfig() {
	Logger = logutil.New()
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	Logger.SetFormatter(&logutil.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999", //
	})

	//是否开启了临时日志，开启了则输出到临时日志，否则输出到控制台
	isTemp := conf.ConfigData.FileLog.TempLog.Status == 1
	if isTemp {
		temp := conf.ConfigData.FileLog.TempLog
		maxAge := time.Duration(temp.LastDays) * time.Hour //临时缓存，按小时清理
		w, _ := createWriter(temp.FilePath, temp.FileName, maxAge, temp.SizeOfFile)
		Logger.SetOutput(w)
	}

	//打印调用位置
	Logger.SetReportCaller(true)

	//设置最低loglevel
	Logger.SetLevel(logutil.ErrorLevel)

	//创建 io.Writer map
	writerMap := make(lfshook.WriterMap, 3)

	if conf.ConfigData.FileLog.DebugLog.Status == 1 {
		//创建debug io.writer
		debug := conf.ConfigData.FileLog.DebugLog
		maxAge := time.Duration(debug.LastDays) * 24 * time.Hour
		writer, _ := createWriter(debug.FilePath, debug.FileName, maxAge, debug.SizeOfFile)
		writerMap[logutil.DebugLevel] = writer
	}
	if conf.ConfigData.FileLog.InfoLog.Status == 1 {
		//创建info io.writer
		info := conf.ConfigData.FileLog.InfoLog
		maxAge := time.Duration(info.LastDays) * 24 * time.Hour
		writer, _ := createWriter(info.FilePath, info.FileName, maxAge, info.SizeOfFile)
		writerMap[logutil.InfoLevel] = writer
	}
	if conf.ConfigData.FileLog.ErrorLog.Status == 1 {
		//创建error io.writer
		errLog := conf.ConfigData.FileLog.ErrorLog
		maxAge := time.Duration(errLog.LastDays) * 24 * time.Hour
		writer, _ := createWriter(errLog.FilePath, errLog.FileName, maxAge, errLog.SizeOfFile)
		writerMap[logutil.ErrorLevel] = writer
	}
	//初始化Hook
	lfsHook := lfshook.NewHook(writerMap, &logutil.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.999",
	})
	Logger.AddHook(lfsHook)

}

// param:
// filepath：文件路径;
// filename:文件名称;
// maxAge:文件保留时间;
// rotationSize :单个文件大小;
//返回值：ioWriter;error
func createWriter(filepath string, filename string, maxAge time.Duration, rotationSize int) (io.Writer, error) {

	return rotatelogs.New(
		// 分割后的文件名称
		path.Join(filepath, conf.ConfigData.Server.ServiceName+"-"+filename)+"-%Y%m%d.log",

		// 生成软链，指向最新日志文件
		rotatelogs.WithLinkName(path.Join(filepath, filename)),

		// 设置最大保存时间(7天)
		rotatelogs.WithMaxAge(maxAge),

		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(24*time.Hour),

		// 设置日志切割大小(单位M)
		rotatelogs.WithRotationSize(rotationSize),
	)
}
