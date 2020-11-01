[TOC]
# 界面展示

### 1、注册界面
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E6%B3%A8%E5%86%8C.png)

### 2、登录界面
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E7%99%BB%E5%BD%95.png)

### 3、管理界面
#### 管理员
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E7%AE%A1%E7%90%86%E5%91%98%E7%AE%A1%E7%90%86%E7%95%8C%E9%9D%A2.png)

#### 普通用户
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E6%99%AE%E9%80%9A%E7%94%A8%E6%88%B7%E7%AE%A1%E7%90%86%E7%95%8C%E9%9D%A2.png)

### 4、聊天界面
#### 发送信息
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E8%81%8A%E5%A4%A91.png)

#### 未读消息红点提示
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E6%9C%AA%E8%AF%BB%E6%B6%88%E6%81%AF%E7%BA%A2%E7%82%B9%E6%8F%90%E7%A4%BA.png)

#### 可发送图片
![](https://github.com/wasiiii/WaSi-webChat/blob/main/images/%E5%8F%AF%E5%8F%91%E9%80%81%E5%9B%BE%E7%89%87.png)

# 功能介绍

### 注册登录模块
1. 用户名按c语言规则
2. 手机邮箱不能重复

### 用户管理模块
1. 管理员可以查看所有人的邮箱和手机
2. 管理员可以修改除自己外所有人的邮箱和密码，手机不可更改
3. 普通用户只能修改自己的邮箱和密码，不能查看信息

### 聊天模块
1. 有未读消息红点提示
2. 发送图片
3. 即时聊天
4. 接收离线消息

# Docker部署
1. 先建立网络
```
$ docker network create wasinet
```
1. Golang：
```shell
$ sudo docker run --rm -it -v $PWD:/go golang:alpine go build main.go -p 9000:9000 --network wasinet
```
2. PostreSQL(要使用init.sql初始化数据库)：
```Dockerfile
Dockerfile

FROM postgres:alpine
ADD init.sql /docker-entrypoint-initdb.d
```
```shell
$ docker build -t myimages:postgres .

$ docker run -d -p 5432:5432 --name=wasipg myimages:postgres --network wasinet
```
3. Angular：
```Dockerfile
Dockerfile

FROM nginx:1.11-1.11-alpine
COPY . /usr/share/nginx/html/
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```
```shell
得到编译文件
$ ng build --prod
```
```shell
进入dist文件夹，把Dockerfile放进去
$ docker build -t myimages:angular .

$ docker run -d --name demo1 -p 4200:4200 myimages:angular --network wasinet
```
4. 注意run的时候要放在同一个net下

# 未解决问题
1. 未读消息红点：正在聊天的过程中，如果点去其他用户，刚刚聊天时的人会有红点，需要点击取消
2. 正则表达式：手机不能正确判断