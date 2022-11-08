# GIN-TEMPLATE

**GIN-TEMPLATE是一款基于GIN的后台框架，支持对接vue-element-admin、vue-admin-template、vue-admin-beautiful等前端框架**

<img align="right" width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">


## 演示地址
#### - [🚀 演示地址：gin-template](http://188.225.25.5:8001/v1)

## 安装
```bash
# 克隆项目
git clone https://github.com/MrLeeang/gin-template.git
# 进入项目目录
cd gin-template
# 安装依赖
go mod tidy
go mod vendor
# 本地开发 启动项目
go run main.go
# 打包
go build
```
## 友情链接

#### - [Element UI 表单设计及代码生成器（可视化表单设计器，一键生成 element 表单）](https://github.com/JakHuang/form-generator/)
#### - [vue-admin-better](https://github.com/chuzhixin/vue-admin-better)
#### - [vue-element-admin](https://github.com/PanJiaChen/vue-element-admin)

## gin-template golang学习交流群-377948518
不管您加或者不加，您都可以享受到开源的代码，感谢您的支持和信任

## config.ini 配置
```
[server]
; 服务端口
serverPort=8001
; 文件上传目录
uploadDir=upload
; 每秒最大访问量
maxRequest=100
; debug开关
debug=false
; 接口加密，返回值加密
encrypt=false

[mysql]
; 数据库地址
host=localhost
; 数据库端口
port=3306
; 数据库用户名
username=root
; 数据库用户密码
password=123456
; 数据库名称
dbname=gintemplate
```

## 付费技术支持

### 联系：VX I-AM-Lihw
![img](https://gitee.com/MrLeeang/image/raw/master/15051057867ab195181e5127ee5c479.jpg)

## 捐赠
![img](https://gitee.com/MrLeeang/image/raw/master/a440e7423e8730f9fa18f95e59dfe6e.jpg)

