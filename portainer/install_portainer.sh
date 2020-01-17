#!/bin/bash

path="/opt/portainer"

info()
{
    echo -e '\033[32m[INFO]\033[0m ' "$@"
}
error()
{
    echo -e '\033[31m[INFO]\033[0m ' "$@"
}

[ -d $path ] ||  mkdir $path
cd $path
curl -so public.zip https://raw.githubusercontent.com/absonggit/tools/master/portainer/public.zip && unzip -o public.zip > /dev/null
curl -so docker-compose.yml  https://raw.githubusercontent.com/absonggit/tools/master/portainer/docker-compose.yml
docker-compose --version &> /dev/null || {
    error "docker-compose 没有安装,开始自动执行安装"
    curl -s https://raw.githubusercontent.com/absonggit/tools/master/docker_compose/install_docker_compose.sh | sh
}
info "开始安装portainer"
docker-compose up -d && info "portainer 安装完成" && docker-compose ps
