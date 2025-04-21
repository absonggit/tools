import requests
import argparse
import json

CADDY_API = "http://127.0.0.1:2019/config/"

def get_config():
    resp = requests.get(CADDY_API)
    if resp.ok:
        return resp.json()
    else:
        print(f"获取配置失败: {resp.status_code} {resp.text}")
        return None

def save_config(cfg):
    resp = requests.post(CADDY_API, json=cfg)
    if resp.ok:
        print("配置已成功保存")
    else:
        print(f"保存失败: {resp.status_code} {resp.text}")

def find_routes(cfg):
    try:
        return cfg["apps"]["http"]["servers"]["custom"]["routes"]
    except KeyError:
        return []

def list_sites():
    cfg = get_config()
    if not cfg:
        return
    routes = find_routes(cfg)
    for idx, route in enumerate(routes):
        host = route.get("match", [{}])[0].get("host", [""])[0]
        body = route.get("handle", [{}])[0].get("body", "")
        target = body.split("url=")[-1].split('"')[0]
        print(f"[{idx}] {host} => {target}")

def add_site(domain, url):
    cfg = get_config()
    if not cfg:
        return

    # 检查域名是否已经存在
    routes = find_routes(cfg)
    for route in routes:
        hosts = route.get("match", [{}])[0].get("host", [])
        if domain in hosts:
            print(f"站点 {domain} 已经存在，跳过添加")
            return

    # 如果没有找到相同的域名，继续添加
    route = {
        "match": [{"host": [domain]}],
        "handle": [{
            "handler": "static_response",
            "body": f'<html><head><meta http-equiv="refresh" content="0;url={url}"></head></html>',
            "headers": {"Content-Type": ["text/html"]}
        }]
    }

    cfg.setdefault("apps", {}).setdefault("http", {}).setdefault("servers", {}).setdefault("custom", {}).setdefault("routes", []).append(route)
    save_config(cfg)
    print(f"已添加: {domain} => {url}")

def update_site(domain, new_url):
    cfg = get_config()
    if not cfg:
        return

    routes = find_routes(cfg)
    for route in routes:
        hosts = route.get("match", [{}])[0].get("host", [])
        if domain in hosts:
            route["handle"][0]["body"] = f'<html><head><meta http-equiv="refresh" content="0;url={new_url}"></head></html>'
            save_config(cfg)
            print(f"已更新: {domain} => {new_url}")
            return
    print("未找到指定域名")

def delete_site(domain):
    cfg = get_config()
    if not cfg:
        return

    routes = find_routes(cfg)
    new_routes = [r for r in routes if domain not in r.get("match", [{}])[0].get("host", [])]
    if len(routes) == len(new_routes):
        print("未找到指定域名")
        return

    cfg["apps"]["http"]["servers"]["custom"]["routes"] = new_routes
    save_config(cfg)
    print(f"已删除: {domain}")

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--list', action='store_true', help='列出所有配置')
    parser.add_argument('--add', nargs=2, metavar=('DOMAIN', 'URL'), help='添加配置')
    parser.add_argument('--update', nargs=2, metavar=('DOMAIN', 'NEW_URL'), help='更新配置')
    parser.add_argument('--delete', metavar='DOMAIN', help='删除配置')
    args = parser.parse_args()

    if args.list:
        list_sites()
    elif args.add:
        add_site(args.add[0], args.add[1])
    elif args.update:
        update_site(args.update[0], args.update[1])
    elif args.delete:
        delete_site(args.delete)
    else:
        parser.print_help()

if __name__ == "__main__":
    main()
