#! /bin/sh
# if [ $# -ne 2 ]
# then
#    echo "缺少必要的参数：-h [域名] -p [端口]"
#    exit
# fi

while getopts :h:p: opts
do
    case $opts in
        h)
        host=$OPTARG
        ;;
        p)
        port=$OPTARG
        ;;
        *)
        echo "-$OPTARG 参数有误"
        exit
        ;;
    esac    
done
end_date=`openssl s_client -host $host -port $port -showcerts </dev/null 2>/dev/null |
      sed -n '/BEGIN CERTIFICATE/,/END CERT/p' |
      openssl x509 -text 2>/dev/null |
      sed -n 's/ *Not After : *//p'`

if [ -n "$end_date" ]
then
    end_date_seconds=`date '+%s' --date "$end_date"`
    now_seconds=`date '+%s'`
    let  ssl_date=($end_date_seconds-$now_seconds)/24/3600
    echo "$1的证书还有$ssl_date天过期"
fi
