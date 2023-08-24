![mrui-logo](https://github.com/Mirouterui/mirouter-ui/assets/63234268/da737f28-e8b6-42d7-a21e-70be2d53fb78)

## Mirouter-ui | 基于小米路由器API的展示面板

将本程序部署在小米路由器的网络环境中，配置完成即可食用

后端基于Golang，多平台兼容

已在小米路由器r1d,r4a上测试通过

部分新路由无法获取cpu占用，如红米AX6000,AX1800。可在路由器上运行解决

### 图片展示

#### 首页

![index](https://github.com/thun888/mirouter-ui/assets/63234268/48bbf554-ec03-41dc-b5fd-42b5faeba466)

#### 设备列表

![Snipaste_2023-08-24_14-53-25](https://github.com/Mirouterui/mirouter-ui/assets/63234268/47309e3a-cc02-479c-a9d3-29cfca235a83)


#### 设备详情

![device_index](https://github.com/thun888/mirouter-ui/assets/63234268/20c465e1-660b-41bf-a200-973423057d31)

#### 路由器状态

![router_index](https://github.com/thun888/mirouter-ui/assets/63234268/1ddce346-7abd-4816-bc55-fe55d3dc70c9)



### 部署

#### 下载

从[Release](https://github.com/thun888/mirouter-ui/releases/)下载二进制文件

> 可访问[镜像站](https://mrui-api.hzchu.top/down/)以获取更快的速度

如果路由器有足够（内存）空间可以下载对应架构版本的部署在路由器上（ps:使用`uname -m`查看，若为armv7l,请使用armv5版本）

![image](https://github.com/Mirouterui/mirouter-ui/assets/63234268/5dfa3deb-0aab-4198-9170-5af1141b3746)



#### 获取key和iv

> 自动获取：[Mirouterui/MiKVIVator](https://github.com/Mirouterui/MiKVIVator)

打开路由器登录页面，右键，点击`查看页面源代码`，按下`CTRL + F`组合键打开搜索框，搜索`key:`，不出意外你能看见以下结果

![image](https://github.com/thun888/mirouter-ui/assets/63234268/87dd59bd-dc9f-4a9f-b22f-d5fd9a9d047a)

复制双引号里的内容粘贴到`config.json`对应栏目中，并填上密码（路由器后台密码）

![image](https://github.com/Mirouterui/mirouter-ui/assets/63234268/56edd993-2119-4979-bb2d-6f822f32059b)


> config.json 会在初次运行时自动下载
> ip可以根据实际情况修改

配置项：

| 配置名     | 默认值       | 解释                                    |
| ---------- | ------------ | --------------------------------------- |
| password   |              | 路由器管理后台密码                      |
| key        |              | 路由器管理后台key                       |
| iv         |              | 路由器管理后台iv                        |
| ip         | 192.168.31.1 | 路由器IP                                |
| tiny       | false        | 启用后，不再下载静态文件，需搭配[在线前端](http://mrui.hzchu.top:8880/)使用|
| routerunit | false        | 启用后，程序通过`gopsutil`库获取CPU占用 |
| port       | 6789         | 网页页面端口号                          |
| debug      | true         | debug模式，建议在测试正常后关闭         |

然后运行程序

如果遇到防火墙提示请勾上两个勾并确定

![image](https://github.com/thun888/mirouter-ui/assets/63234268/fc6a7515-6e65-48be-9bbd-1de1eac41146)

此时命令窗口中会显示网页的访问端口，使用设备的`ip地址+端口号(6789)`访问面板

### 后台运行

自行参考：

[Linux命令后台运行_后台运行命令_拉普拉斯妖1228的博客-CSDN博客](https://blog.csdn.net/caesar1228/article/details/118853871)

[windows守护进程工具--nssm详解 - 与f - 博客园 (cnblogs.com)](https://www.cnblogs.com/fps2tao/p/16433588.html)
