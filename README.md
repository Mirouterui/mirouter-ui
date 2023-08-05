## Mirouter-ui |基于小米路由器API的展示面板

将本程序部署在小米路由器的网络环境中，配置完成即可食用

### 图片展示

#### 首页

![index](https://github.com/thun888/mirouter-ui/assets/63234268/48bbf554-ec03-41dc-b5fd-42b5faeba466)

#### 设备详情

![device_index](https://github.com/thun888/mirouter-ui/assets/63234268/20c465e1-660b-41bf-a200-973423057d31)

#### 路由器状态

![router_index](https://github.com/thun888/mirouter-ui/assets/63234268/1ddce346-7abd-4816-bc55-fe55d3dc70c9)



### 部署

#### 下载

WIN：从[release](https://github.com/thun888/mirouter-ui/releases/tag/zip)下载压缩包
Linux: fork项目仓库，并安装相应依赖

#### 获取key和iv

打开路由器登录页面，右键，点击`查看页面源代码`，按下`CTRL + F`组合键打开搜索框，搜索`key:`，不出意外你能看见以下结果

![image](https://github.com/thun888/mirouter-ui/assets/63234268/87dd59bd-dc9f-4a9f-b22f-d5fd9a9d047a)

复制双引号里的内容粘贴到配置文件中，并填上密码
![image](https://github.com/thun888/mirouter-ui/assets/63234268/3aed12f2-255e-4b30-a765-eb8b4a963995)

> ip可以根据实际情况修改

然后双击`mirouterui.exe`运行（linux执行`python3 ./mirouter.py`）

如果遇到防火墙提示请勾上两个勾并确定

![image](https://github.com/thun888/mirouter-ui/assets/63234268/fc6a7515-6e65-48be-9bbd-1de1eac41146)

此时命令窗口中会显示网页的访问地址，但只有以`路由器分配的IP地址开头的`才能被其他设备访问

![image](https://github.com/thun888/mirouter-ui/assets/63234268/5e05fde4-a62d-4f92-93af-6020242c36e3)

### 后台运行

自行参考：

[Linux命令后台运行_后台运行命令_拉普拉斯妖1228的博客-CSDN博客](https://blog.csdn.net/caesar1228/article/details/118853871)

[windows守护进程工具--nssm详解 - 与f - 博客园 (cnblogs.com)](https://www.cnblogs.com/fps2tao/p/16433588.html)
