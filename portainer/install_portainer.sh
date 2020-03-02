#!/bin/bash

path="/opt/portainer"

info() {
    echo -e '\033[32m[INFO]\033[0m ' "$@"
}

warn() {
    echo -e '\033[33m[WARN]\033[0m ' "$@"
}

error() {
    echo -e '\033[31m[ERROR]\033[0m ' "$@"
}

init() {
    info "初始化..."
    docker-compose --version &> /dev/null || {
        warn "docker-compose 没有安装,开始自动执行安装"
        curl -s https://raw.githubusercontent.com/absonggit/tools/master/docker_compose/install_docker_compose.sh | sh
    }
    if [ ! -d $path ]
    then
        mkdir $path
    fi   
    cd $path
}

get_docker_compose() {
    info "下载docker-compose.yml"
    curl -so docker-compose.yml  https://raw.githubusercontent.com/absonggit/tools/master/portainer/docker-compose.yml
}
get_public() {
    info "下载汉化包"
    yum install unzip -y &> /dev/null
    curl -so public.zip https://raw.githubusercontent.com/absonggit/tools/master/portainer/public.zip && unzip -o public.zip > /dev/null
}
install_portainer() {
    info "开始安装portainer"
    docker-compose up -d && info "安装完成" && docker-compose ps
    }
init
get_docker_compose
get_public
install_portainer
