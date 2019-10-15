#!/bin/bash
cd ..

echo $(date +%Y/%m/%d/%H:%M:%S) 密钥同步服务进程监控开始……
echo 监控程序进程id是: $$ >> logs/watch.log
num=1
iNum=1

while(( $num < 5))
	do
	sn=`ps -ef | grep ./main | grep -v grep |awk '{print $2}'`
	#echo $sn
	if [ "$sn"x = ""x ]
	then
	    echo $(date +%Y/%m/%d/%H:%M:%S) 没有发现密钥同步服务进程 >> logs/watch.log
		echo $(date +%Y/%m/%d/%H:%M:%S) 正在重新启动密钥同步服务 >> logs/watch.log
		let "iNum++"
		#echo $iNum
		nohup ./main >> logs/watch.log 2>&1 &
		echo $(date +%Y/%m/%d/%H:%M:%S)  密钥同步服务启动成功，请查看logs/watch.log
		echo $(date +%Y/%m/%d/%H:%M:%S)  密钥同步服务启动成功 >> logs/watch.log
	else
		echo $(date +%Y/%m/%d/%H:%M:%S) 密钥同步服务正在运行……
	fi
	sleep 5
done
