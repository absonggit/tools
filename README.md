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
- caddy caddy通过api动态管理站点
  - 
