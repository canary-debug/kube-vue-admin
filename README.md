# kube-vue-admin

Kubernetes Web 管理平台后端服务，基于 Go 语言和 Gin 框架构建，为 Kubernetes 集群提供可视化管理界面的 API 支持。

## 项目简介

kube-vue-admin 是一个轻量级的 Kubernetes 管理后台后端服务，提供用户认证、集群管理等核心功能。该项目采用现代化的 Go 语言开发模式，集成 RESTful API 设计，与前端 Vue.js 界面完美配合。

## 功能特性

### 认证管理
- 用户登录接口 (`/api/auth/login`)
- 用户注册接口 (`/api/auth/register`)

### 核心服务
- 基于 Gin 框架的高性能 HTTP 服务
- MySQL 数据库集成 ([database.ConnectDatabase](file://E:\code\go\kube-vue-admin\database\db.go#L10-L27))
- 双通道日志记录 (文件 + 控制台)
- 自动异常恢复机制

### 技术架构
- **编程语言**: Go
- **Web框架**: Gin
- **数据库**: MySQL
- **API风格**: RESTful

## 快速开始

### 环境要求
- Go 1.16+
- MySQL 5.7+
- Kubernetes 集群访问权限

### 安装部署
```bash
# 克隆项目
git clone <repository-url>

# 进入项目目录
cd kube-vue-admin

# 下载依赖
go mod tidy

# 编译运行
go run main.go
```


### 服务配置
- 默认监听端口: `9000`
- 日志文件: [gin.log](file://E:\code\go\kube-vue-admin\gin.log)
- 数据库配置: 通过环境变量或配置文件设置

## API 接口

### 健康检查
```
GET /
响应: {"message": "pong"}
```


### 认证接口
```
POST /api/auth/login    # 用户登录
POST /api/auth/register # 用户注册
```


## 项目结构
```
kube-vue-admin/
├── api/          # API 接口实现
├── database/     # 数据库连接和模型
├── main.go       # 程序入口
└── go.mod        # 依赖管理
```


## 开发指南

### 本地开发
1. 配置本地 MySQL 数据库
2. 设置必要的环境变量
3. 运行 `go run main.go` 启动服务
4. 访问 `http://localhost:9000` 验证服务

### 日志查看
- 实时日志输出到控制台
- 完整日志保存在 [gin.log](file://E:\code\go\kube-vue-admin\gin.log) 文件中

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目。

## 许可证

[待添加许可证信息]

---
*该项目为 Kubernetes 可视化管理平台的后端服务组件*

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o kube-vue-admin . 