cd ..

@echo off
rem 后续命令使用的是：UTF-8编码
chcp 65001

rem 定义需监控程序的进程名和程序路径，可根据需要进行修改

set AppName=registryCenter.exe

set AppPath=%cd%/

title 进程监控

cls

echo.

echo 进程监控开始……

echo.

rem 定义循环体

:loop

   rem 从进程列表中查找指定进程
   rem  下面语句也可写成 qprocess %AppName% >nul （经验发布后补充）

    tasklist | findstr %AppName% >nul

   rem 变量errorlevel的值等于0表示查找到进程，否则没有查找到进程

   if %errorlevel%==0 (

         echo ^>%date:~0,10% %time:~0,8% 程序正在运行……

    )else (

           echo ^>%date:~0,10% %time:~0,8% 没有发现程序进程 >> logs/watch.log

           echo ^>%date:~0,10% %time:~0,8% 正在重新启动程序 >> logs/watch.log

           start %AppPath%%AppName% 2>nul && echo ^>%date:~0,10% %time:~0,8% 启动程序成功 >> logs/watch.log

   )

   rem 用ping命令来实现延时运行(10秒)
   for /l %%i in (1,1,10) do ping -n 1 -w 1000 127.0.0.1>nul

   goto loop

echo on