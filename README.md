# tools

- docker和docker-compose
  - 安装脚本
- portainer
  - 安装脚本(汉化)
- ssl_check
  - 证书过期时间脚本
- cc
  - banIP脚本
- files_to_tg
  - 多文件压缩发送TG 
  - 配置文件 config.yaml
- file_watch_to_tg
  - 监控文件内容，匹配关键字，发送内容到TG
  - 配置文件 config.yaml
- tgwebhook
  - tg 发消息webhook
  - HTTP方法: POST
  - ```json
    {
    "chatid": "接收的ChatID",
    "token": "TG Bot Token",
    "text": "要发送的消息内容"
    }
    ```
- ModSecurity 规则生成器
  - [规则生成器](https://absonggit.github.io/tools)
- go-cron 支持秒级的定时器可以执行shell命令（替代crontab）
- caddy 通过api动态管理站点
  - autosave.json 持久化配置文件，初始json
  - caddy.service systemd 守护配置
  - caddy_manager.py python管理脚本
    - ```
      usage: caddy_manage.py [-h] [--list] [--add DOMAIN URL] [--update DOMAIN NEW_URL] [--delete DOMAIN]
      optional arguments:
      -h, --help              show this help message and exit
      --list                  列出所有配置
      --add DOMAIN URL        添加配置
      --update DOMAIN NEW_URL 更新配置
      --delete DOMAIN         删除配置
      ```
