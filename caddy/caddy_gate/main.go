---
- name: 部署 caddy 及其配置
  hosts: all
  become: true

  tasks:
    - name: 创建 /opt/caddy 目录
      file:
        path: /opt/caddy
        state: directory
        mode: '0755'

    - name: 上传 caddy 二进制
      copy:
        src: caddy
        dest: /opt/caddy/caddy
        mode: '0755'

    - name: 上传 caddy_gate 二进制
      copy:
        src: caddy_gate
        dest: /opt/caddy/caddy_gate
        mode: '0755'

    - name: 上传 Caddyfile
      copy:
        src: Caddyfile
        dest: /opt/caddy/Caddyfile
        mode: '0644'

    - name: 上传 config.yaml
      copy:
        src: config.yaml
        dest: /opt/caddy/config.yaml
        mode: '0644'

    - name: 安装 supervisor
      package:
        name: supervisor
        state: present

    - name: 创建 /etc/supervisord.d 目录
      file:
        path: /etc/supervisord.d
        state: directory
        mode: '0755'

    - name: 上传 caddy_node.ini
      copy:
        src: caddy_node.ini
        dest: /etc/supervisord.d/caddy_node.ini
        mode: '0644'

    - name: 重载 supervisor 配置
      shell: supervisorctl reread && supervisorctl update

    - name: 启动 caddy 服务
      shell: supervisorctl start caddy || true # 已启动不报错
