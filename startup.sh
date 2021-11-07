#!/usr/bin/env bash

process="goblog"
echo ${process} $1

start(){
      pid=`pgrep ${process}`
      if [ "${pid}"x = ""x ];then
          echo "start new process..."
          nohup ./${process} -c conf/prod.yml &
      else
          for i in ${pid}
          do
              echo "reload the process [ $i ]"
              kill -9 $i
          done
      fi
      sleep 1
      nohup ./${process} -c conf/prod.yml &
      pid=`pgrep ${process}`
      echo "new process id: ${pid}"
}

stop(){
    pid=`pgrep ${process}`
    echo ${pid}
    for i in ${pid}
    do
        echo "kill the process [ $i ]"
	kill -9 $i
    done
}

status(){
    ps aux | grep -w ${process} | grep -v 'grep'
}


case "$1" in
    start)
	start $1;;
    stop)
	stop ;;
    status)
	status ;;
    *)
	echo "Usage: $0 {start|stop|reload|status}"

esac
