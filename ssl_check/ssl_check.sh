#!/bin/bash
while getopts h:p: opts
do
    case $opts in
        h)
        host=$OPTARG
        ;;
        p)
        port=$OPTARG
        ;;
        *)
        echo "ERROR：-h 或 -p 有误"
        echo "OPTIONS:
    -h:域名(必选)
    -p:端口(必选)"
        exit
        ;;
    esac
done

if [ ! $host -o $port ]
then
    echo "ERROR: -h 或 -p 有误"
    echo "OPTIONS:
    -h:域名(必选)
    -p:端口(必选)"
else
    end_date=`openssl s_client -host $host -port $port -showcerts </dev/null 2>/dev/null |
          sed -n '/BEGIN CERTIFICATE/,/END CERT/p' |
          openssl x509 -text 2>/dev/null |
          sed -n 's/ *Not After : *//p'`
    if [ -n "$end_date" ]
    then
        end_date_seconds=`date '+%s' --date "$end_date"`
        now_seconds=`date '+%s'`
        let  ssl_date=($end_date_seconds-$now_seconds)/24/3600
        echo "$host的证书还有$ssl_date天过期"
    else
        echo "ERROR: -h 或 -p 有误"
        echo "OPTIONS:
    -h:域名(必选)
    -p:端口(必选)"
    fi
fi
