## Mirouter-ui |基于小米路由器API的展示面板

将本程序部署在小米路由器的网络环境中，配置完成即可食用

后端基于Golang，多平台兼容

已在小米路由器r1d,r4a,ra71上测试通过

部分新路由无法获取cpu占用，如红米AX6000

### 图片展示

#### 首页

![index](https://github.com/thun888/mirouter-ui/assets/63234268/48bbf554-ec03-41dc-b5fd-42b5faeba466)

#### 设备详情

![device_index](https://github.com/thun888/mirouter-ui/assets/63234268/20c465e1-660b-41bf-a200-973423057d31)

#### 路由器状态

![router_index](https://github.com/thun888/mirouter-ui/assets/63234268/1ddce346-7abd-4816-bc55-fe55d3dc70c9)



### 部署

#### 下载

从[Release](https://github.com/thun888/mirouter-ui/releases/)下载二进制文件和`config.txt`

如果路由器有足够空间可以下载`mipsle`版本的部署在路由器上（理论上）

> 如果自己新建config.txt千万不要新建成config.txt.txt

#### 获取key和iv

打开路由器登录页面，右键，点击`查看页面源代码`，按下`CTRL + F`组合键打开搜索框，搜索`key:`，不出意外你能看见以下结果

![image](https://github.com/thun888/mirouter-ui/assets/63234268/87dd59bd-dc9f-4a9f-b22f-d5fd9a9d047a)

复制双引号里的内容粘贴到`config.txt`对应栏目中，并填上密码（路由器后台密码）
![image](https://github.com/thun888/mirouter-ui/assets/63234268/b581d6b9-c56e-4ce4-a356-167c6856cdf9)

> ip可以根据实际情况修改

然后运行程序

如果遇到防火墙提示请勾上两个勾并确定

![image](https://github.com/thun888/mirouter-ui/assets/63234268/fc6a7515-6e65-48be-9bbd-1de1eac41146)

此时命令窗口中会显示网页的访问端口，使用设备的`ip地址+端口号(6789)`访问面板

### 后台运行

自行参考：

[Linux命令后台运行_后台运行命令_拉普拉斯妖1228的博客-CSDN博客](https://blog.csdn.net/caesar1228/article/details/118853871)

[windows守护进程工具--nssm详解 - 与f - 博客园 (cnblogs.com)](https://www.cnblogs.com/fps2tao/p/16433588.html)
