#!/bin/bash
# shell 命令行执行安装 
# curl -s https://raw.githubusercontent.com/absonggit/tools/master/docker_compose/install_docker_compose.sh |sh
info()
{
    echo -e '\033[32m[INFO]\033[0m ' "$@"
}
install_docker_compose() {
    latest=`curl -s  https://github.com/docker/compose/tags |grep "<a href=\"\/docker\/compose\/releases\/tag\/"  |awk -F'/|"' '{print $(NF-1)}' | grep -v "-" | head -1`
    info "开始安装 docker-compose $latest"
    sudo curl -L --progress-bar "https://github.com/docker/compose/releases/download/${latest}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/sbin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    info "安装完成"
    docker-compose --version
}

install_docker_compose
