import random
import time
import hashlib
from flask import Flask, jsonify, redirect
import threading
import requests
import json
import sys
app = Flask(__name__)
with open('config.txt', 'r') as f:
    contents = f.read()
    password = contents.split('password = "')[1].split('"')[0]
    key = contents.split('key = "')[1].split('"')[0]
    iv = contents.split('iv = "')[1].split('"')[0]
    ip = contents.split('ip = "')[1].split('"')[0]
def nonce_creat():
    type_var = 0
    deviceId = '2a:45:6f:11:0b:a5'
    time_var = int(time.time())
    random_var = random.randint(0, 9999)
    return f"{type_var}_{deviceId}_{time_var}_{random_var}"

def hash_password(pwd, nonce, key):
    pwd_key_hash = hashlib.sha1((pwd + key).encode('utf-8')).hexdigest()
    nonce_pwd_key_hash = hashlib.sha1((nonce + pwd_key_hash).encode('utf-8')).hexdigest()
    return nonce_pwd_key_hash


def upstok():
    global token
    nonce = nonce_creat()
    hashed_password = hash_password(password, nonce, key)
    url = f'http://{ip}/cgi-bin/luci/api/xqsystem/login'
    data = {
        'username': 'admin',
        'password': hashed_password,
        'logtype': '2',
        'nonce': nonce
    }

    response = requests.post(url, data=data)
    if response.json()['code'] != 0:
        print("登录失败，请检查配置或路由器状态")
        sys.exit()
    token = response.json()['token']
    # print(token)

@app.route('/api/<path:apipath>', methods=['GET'])
def api_proxy(apipath):
    url = f'http://{ip}/cgi-bin/luci/;stok={token}/api/{apipath}'
    try:
        response = requests.get(url)
        response.raise_for_status() # 如果响应码不是 2XX，抛出异常
        data = json.loads(response.text)
        return data
    except requests.exceptions.RequestException as e:
        return jsonify({'code':1101,'msg': 'MiRouterのapi调用出错，请检查配置或路由器状态'}), 200

@app.route('/')
def home():
    return redirect('/index.html')

@app.route('/<path:path>')
def static_file(path):
    return app.send_static_file(path)
   
def run_timer():
    while True:                                                                                                                                                                              
        # 每30分钟执行一次
        time.sleep(30 * 60)
        # 在后台线程中执行需要定时执行的函数
        threading.Thread(target=upstok).start()

if __name__ == "__main__":
    upstok()
    timer_thread = threading.Thread(target=run_timer)
    timer_thread.daemon = True
    timer_thread.start()
    app.run(host='0.0.0.0',port=6789)