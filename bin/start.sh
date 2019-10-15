#!/bin/sh
# $0 表shell命令第一个参数，即脚本本身
# $(dirname $0) 定位脚本文件所在的目录
BASE_DIR=$(dirname $0)
# $(basename $0 .sh) 定位脚本名称，.sh表示去除.sh后缀
PID=$(dirname $0)/bin/$(basename $0 .sh).pid

#USER=root

if [ -f $PID ]
    then
        echo "---Already started. PID: [$( cat $PID )]---"
    else
      # 项目所在路径
	  cd ..
      ./main &
      # $? 为上条命令执行结果，成功执行则返回0
      if [[ "$?" -eq 0 ]];
      then
        # $! 为上条shell命令的进程号，如执行成功，则将进程号写入pid文件
        echo $!>$PID
        echo "START success"
      else
        echo "---START failure, please check----"
      fi
    fi

exit 0