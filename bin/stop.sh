#!/bin/sh
# $0 表shell命令第一个参数，即脚本本身
# $(dirname $0) 定位脚本文件所在的目录
BASE_DIR=$(dirname $0)
# $(basename $0 .sh) 定位脚本名称，.sh表示去除.sh后缀
PID=$(dirname $0)/start.pid

#USER=root
  if [ -f $PID ]
  then
    # 执行kill的几种命令方式
    # kill -9 `ps -ef | grep -v grep | grep $CMD | grep -w $USR | awk '{print $2}'`
    # echo `cat $PID` | xargs kill -9
    kill -9 `cat $PID`
	/bin/rm $PID
    sleep 1
    echo "Stop success"
  else
    echo "---No pid file---"
  fi
exit 0