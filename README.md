# kube-vue-admin

Kubernetes Web 管理平台后端服务，基于 Go 语言和 Gin 框架构建，为 Kubernetes 集群提供可视化管理界面的 API 支持。

## 项目简介

kube-vue-admin 是一个轻量级、高性能的 Kubernetes 管理后台后端服务，提供完整的用户认证、集群管理、资源监控等核心功能。该项目采用现代化的 Go 语言开发模式，集成 RESTful API 设计，与前端 Vue.js 界面完美配合，为用户提供直观、高效的 Kubernetes 集群管理体验。

## 功能特性

### 认证与授权
- 用户登录与注册
- JWT 令牌管理
- 基于角色的访问控制

### Kubernetes 集群管理
- 集群资源概览
- 节点、Pod、Service 管理
- 命名空间与资源配额管理
- 部署与状态监控

### 核心服务特性
- 基于 Gin 框架的高性能 HTTP 服务
- MySQL 数据库集成与连接池管理
- 双通道日志记录 (文件 + 控制台)
- 自动异常恢复机制
- RESTful API 设计规范
- 完整的 API 文档

### 技术架构
- **编程语言**: Go
- **Web框架**: Gin
- **数据库**: MySQL
- **API风格**: RESTful
- **认证方式**: JWT
- **日志系统**: 自定义日志记录

## 快速开始

### 环境要求
- Go 1.16+
- MySQL 5.7+
- Kubernetes 集群访问权限 (kubeconfig 文件)

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

#### 环境变量
| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `DB_HOST` | 数据库主机地址 | `localhost` |
| `DB_PORT` | 数据库端口 | `3306` |
| `DB_USER` | 数据库用户名 | `root` |
| `DB_PASSWORD` | 数据库密码 | - |
| `DB_NAME` | 数据库名称 | `kube_vue_admin` |
| `JWT_SECRET` | JWT 密钥 | 从 `.jwt_key` 文件读取 |
| `SERVER_PORT` | 服务监听端口 | `9000` |
| `KUBECONFIG` | Kubernetes 配置文件路径 | `~/.kube/config` |

#### 默认配置
- 默认监听端口: `9000`
- 日志文件: `gin.log`
- JWT 密钥文件: `.jwt_key`

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

### 更多 API 文档
完整的 API 文档请查看 [api_documentation.md](api_documentation.md) 文件。

## 项目结构

```
kube-vue-admin/
├── api/          # API 接口实现
├── database/     # 数据库连接和初始化
├── kube-test/    # Kubernetes 测试相关代码
├── models/       # 数据模型定义
├── pkg/          # 公共包和工具函数
├── sql/          # SQL 迁移和初始化脚本
├── tokens/       # JWT 令牌管理
├── .gitignore    # Git 忽略文件
├── .jwt_key      # JWT 密钥文件
├── api_documentation.md # API 文档
├── gin.log       # 日志文件
├── go.mod        # Go 模块定义
├── go.sum        # 依赖版本锁定
├── LICENSE       # 许可证文件
└── main.go       # 程序入口
```

### 目录说明
- **api/**: 包含所有 API 接口的路由和处理逻辑
- **database/**: 数据库连接配置和初始化代码
- **models/**: 定义数据库表结构和数据模型
- **pkg/**: 公共工具包，包括 Kubernetes 客户端、日志等
- **tokens/**: JWT 令牌生成、验证和管理
- **sql/**: 数据库迁移脚本和初始化 SQL

## 开发指南

### 本地开发
1. 配置本地 MySQL 数据库
2. 设置必要的环境变量或修改配置
3. 运行 `go run main.go` 启动服务
4. 访问 `http://localhost:9000` 验证服务
5. 使用 API 测试工具（如 Postman）测试接口

### 日志查看
- 实时日志输出到控制台
- 完整日志保存在 `gin.log` 文件中
- 日志包含请求信息、错误详情和操作记录

### 代码结构
- 每个 API 接口对应一个处理函数
- 数据模型与数据库表一一对应
- 业务逻辑与 API 处理分离
- 统一的错误处理和响应格式

## 部署方式

### 本地开发环境
```bash
go run main.go
```

### 编译部署
```bash
# 编译为当前系统可执行文件
go build -o kube-vue-admin main.go

# 编译为 Linux 系统可执行文件（Windows PowerShell）
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o kube-vue-admin main.go

# 编译为 Linux 系统可执行文件（Linux/Mac）
export GOOS=linux
export GOARCH=amd64
go build -o kube-vue-admin main.go
```

## 贡献指南

欢迎提交 Issue 和 Pull Request 来改进项目。贡献前请确保：

1. 代码符合 Go 语言规范
2. 添加必要的测试用例
3. 更新相关文档
4. 遵循项目的代码风格

## 许可证

本项目采用 MIT 许可证，详情请查看 [LICENSE](LICENSE) 文件。

## 联系方式

- 项目地址: <repository-url>
- 问题反馈: [Issues](<repository-url>/issues)

---
*该项目为 Kubernetes 可视化管理平台的后端服务组件*