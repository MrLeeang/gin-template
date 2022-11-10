# GIN-TEMPLATE

**GIN-TEMPLATE是一款基于GIN的后台框架，支持对接vue-element-admin、vue-admin-template、vue-admin-beautiful等前端框架**

<img align="right" width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">

## 🎉 特性

- 💪 AES加密
- 💅 RBAC 模型
- 🌍 JWT 权限控制
- 📦️ 接口流量控制
- 💪 日志管理
- 📦️ 微服务架构
- 🌍 短信服务
- 💪 邮件服务
- 💅 配置管理
- 👏 良好的类型定义
- 🥳 开源版本支持免费商用

## 演示地址
#### - [🚀 演示地址：gin-template](https://documenter.getpostman.com/view/7717980/2s8YYPGKZR)

- [🌐 github 仓库地址](https://github.com/MrLeeang/gin-template)

- [🌐 码云仓库地址](https://gitee.com/MrLeeang/gin-template)

## 安装
```bash
# 克隆项目
git clone https://github.com/MrLeeang/gin-template.git
# 进入项目目录
cd gin-template
# 安装依赖
go mod tidy
# 本地开发 启动项目
go run cmd/appv1/main.go
go run cmd/service/main.go
# 打包
go build -o app cmd/appv1/main.go
go build -o srv cmd/service/main.go
```
## 友情链接

#### - [Element UI 表单设计及代码生成器（可视化表单设计器，一键生成 element 表单）](https://github.com/JakHuang/form-generator/)
#### - [Gin Web Framework](https://github.com/gin-gonic/gin)
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

[service]
; 微服务地址
address=localhost:8090

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

[consul]
address=localhost:8500

[mail]
; 登录地址
host=smtp.163.com
; 登录账号
username=xxx@163.com
; 登录密码
password=xxx
; 邮件服务
address=smtp.163.com:25
; 发件人邮箱地址
from=xxx@163.com

; 短信服务
[alibaba]
accessKeyId=
accessKeySecret=
signName=阿里云短信测试
templateCode=SMS_154950909
```

## 付费技术支持

### 联系：VX I-AM-Lihw
<table>
<tr>
<td>
<img align="left" width="200px" src="https://gitee.com/MrLeeang/image/raw/master/15051057867ab195181e5127ee5c479.jpg">
</td>
</tr>
</table>

## 捐赠
<table>
<tr>
<td>
<img align="left" width="200px" src="https://gitee.com/MrLeeang/image/raw/master/a440e7423e8730f9fa18f95e59dfe6e.jpg">
</td>
</tr>
</table>
