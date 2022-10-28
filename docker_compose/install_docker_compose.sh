#!/bin/bash
info()
{
    echo -e '\033[32m[INFO]\033[0m ' "$@"
}
install_docker_compose() {
    latest=`curl -s https://github.com/docker/compose/tags |grep "<a href=\"\/docker\/compose\/releases\/tag\/" | awk '{print $NF}'| awk -F"<|>" 'NR==1{print $2}'`
    info "开始安装 docker-compose $latest"
    sudo curl -L --progress-bar "https://github.com/docker/compose/releases/download/${latest}/docker-compose-$(uname -s)-$(uname -m)" -o /bin/docker-compose
    sudo chmod +x /bin/docker-compose
    info "安装完成"
    docker-compose --version
}
install_docker_compose
