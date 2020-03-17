#!/bin/bash

# 检测的日志路径
log_path=/home/wwwlogs/app.com.log
# 阈值（每分钟IP的访问量）
vpt=5
# 解封时间
expire_time=1
# 被拒过的IP
black_ip=/tmp/black_ip_list
# 获取到达到阈值的IP
ip_list=$(grep "$(date +'%d/%b/%Y:%H:%M' -d "1 minute ago")" $log_path  | awk '{print $3}' | sort  | uniq -c | awk '$1>"'$vpt'"{print $2}' |xargs)

# 遍历IP数组执行iptables规则
function main() {
    while [ $# != 0 ]
    do
       iptables -A INPUT -s $1 -j DROP
       # 定义临时任务，达到解封时间执行iptables规则
       export ban_ip=$1
       at now +$expire_time minutes << EOF
echo $ban_ip >> $black_ip
iptables -D INPUT -s $ban_ip -j DROP
EOF
    shift
    done
}

main ${ip_list[@]}
