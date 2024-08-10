import time
import yaml
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
import os
import logging
import requests
import re
import json

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class FileChangeHandler(FileSystemEventHandler):
    def __init__(self, target_files, keywords, bot_token, chat_id):
        self.target_files = set(target_files)
        self.keywords = keywords
        self.bot_token = bot_token
        self.chat_id = chat_id
        self.last_modified = {}
        self.file_offsets = {file: self.get_file_size(file) for file in target_files}

    def get_file_size(self, file_path):
        try:
            return os.path.getsize(file_path)
        except OSError:
            return 0

    def on_modified(self, event):
        if event.is_directory:
            return

        if event.src_path in self.target_files:
            if (event.src_path not in self.last_modified):
                self.print_new_file_content(event.src_path)
    def print_new_file_content(self, file_path):
        try:
            with open(file_path, 'r', encoding='utf-8') as file:
                file.seek(self.file_offsets[file_path])
                new_content = file.read().strip()
                if new_content:
                    for keyword in self.keywords:
                        if keyword in new_content:
                            message = {}
                            message['date'] = extractors['date'](new_content)
                            message['file'] = extractors['file_path'](new_content)
                            message['line'] = extractors['line_number'](new_content) 
                            message['id'] = extractors['rule_id'](new_content)
                            message['ip'] = extractors['client_ip'](new_content)
                            message['host'] = extractors['host'](new_content)
                            message['request'] = extractors['request'](new_content)
                            message['msg'] = extractors['message'](new_content)
                            text = f"""
**WAF 通知**
```
时间: {message['date']}
规则ID: {message['id']}
规则文件: {message['file']} 第 {message['line']} 行
请求IP: {message['ip']}
请求域名: {message['host']}
请求内容: {message['request']}
拦截消息: {message['msg']}
```
"""
                            logging.info(f"{file_path} 匹配到关键字{keyword},发送telegram 通知")
                            send_message_to_telegram(self.bot_token, self.chat_id, text)
                            break
                self.file_offsets[file_path] = file.tell()
        except Exception as e:
            logging.error(f"Could not read file {file_path}: {e}")

def load_config(config_path):
    with open(config_path, 'r') as file:
        config = yaml.safe_load(file)
    return config

def send_message_to_telegram(bot_token, chat_id, message):
    url = f'https://api.telegram.org/bot{bot_token}/sendMessage'
    data = {
        'chat_id': chat_id,
        'text': message,
        "parse_mode": "MarkdownV2"
    }
    try:
        response = requests.post(url, data=data)
        if response.status_code != 200:
            logging.error(f"Failed to send message: {response.text}")
    except requests.exceptions.RequestException as e:
        logging.error(f"Request to Telegram failed: {e}")

def create_extractor(pattern, group_index=1):
    def extractor(log_line):
        match = re.search(pattern, log_line)
        return match.group(group_index) if match else None
    return extractor

extractors = {
    'date': create_extractor(r'\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}', 0),
    'file_path': create_extractor(r'\[file "([^"]+)"\]'),
    'line_number': create_extractor(r'\[line "(\d+)"\]'),
    'rule_id': create_extractor(r'\[id "(\d+)"\]'),
    'client_ip': create_extractor(r'\[client ([\d\.]+)\]'),
    'request': create_extractor(r'request: "([^"]+)"'),
    'host': create_extractor(r'host: "([^"]+)"'),
    'message': create_extractor(r'\[msg "([^"]+)"\]')
}

def main():
    config = load_config('config.yaml')
    target_files = config['files']
    keywords = config['keywords']
    bot_token = config['telegram']['bot_token']
    chat_id = config['telegram']['chat_id']
    logging.info(f"开始监听文件: {target_files}")
    
    observer = Observer()
    event_handler = FileChangeHandler(target_files, keywords, bot_token, chat_id)

    directories_to_watch = {file_path.rsplit('/', 1)[0] if '/' in file_path else '.' for file_path in target_files}
    
    for directory in directories_to_watch:
        observer.schedule(event_handler, path=directory, recursive=False)
    observer.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt as e:
        print("检测到用户中断，正在安全退出...")
        observer.stop()
    observer.join()

if __name__ == "__main__":
    main()
