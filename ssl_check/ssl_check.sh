#!/bin/bash

func_usage() {
    cat <<- EOF
    ERROR：-h 或 -p 有误
    OPTIONS:
        -h:域名(必选)
        -p:端口(必选)"
    EOF
}

func_opts() {
    while getopts h:p: opts
    do
        case $opts in
            h)host=$OPTARG;;
            p)port=$OPTARG;;
            *)
            func_err
            exit;;
        esac
    done
    func_check
}
func_check() {
    if [ ! $host -o ! $port ]
    then
        func_err
    else
        end_date=`openssl s_client -host $host -port $port -showcerts </dev/null 2>/dev/null|
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
            func_err
        fi
    fi
}


until [ $# -eq 0 ]
do
   str="$str $1"
   shift
done
func_opts $str
