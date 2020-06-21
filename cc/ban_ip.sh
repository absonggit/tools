#!/bin/bash
################################## 脚本说明 ####################################
## 统计nginx访问日志前一分钟之内，将请求数大于阈值的IP添加至防火墙定义规则，禁止其访问
## 并增加at任务在指定过期时间后清除相应防火墙规则                                   
## nginx日志格式                                                                
## 第一个字段为 [17/Mar/2020:15:01:43 +0800]
## 第二个字段为 客户端的真实IP
###############################################################################
# 检测的日志路径
log_path=/home/wwwlogs/app.com.log
# 阈值（每分钟IP的访问量）
vpt=1000
# 解封时间
expire_time=30
# 被拒过的IP
black_ip=/tmp/black_ip_list
# 获取到达到阈值的IP
ip_list=$(grep "$(date +'%d/%b/%Y:%H:%M' -d "1 minute ago")" $log_path  | awk '{print $3}' | sort  | uniq -c | awk '$1>'$vpt'{print $2}' |xargs)

# 遍历IP数组执行iptables规则
function main() {
    while [ $# != 0 ]
    do
       if ! iptables -L INPUT | grep $1 &> /dev/null
       then
       iptables -A INPUT -s $1 -j DROP
       echo "$(date +'%D-%T') IP: $1 封禁" >> $black_ip
       # 定义临时任务，达到解封时间执行iptables规则
       export ban_ip=$1
       at now +$expire_time minutes << EOF
echo "$(date +'%D-%T') IP: $ban_ip 解封" >> $black_ip
iptables -D INPUT -s $ban_ip -j DROP
EOF
       fi
    shift
    done
}

main ${ip_list[@]}
