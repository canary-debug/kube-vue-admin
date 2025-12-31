# Kube-Vue-Admin 后端API接口文档

## 1. 认证相关接口

### 1.1 用户登录接口

**接口名称**: 用户登录
**请求方法**: POST
**请求URL**: `/api/auth/login`
**请求参数**:
```json
{
  "username": "admin",     // 用户名(必填，长度3-20)
  "password": "password123" // 密码(必填，长度6-30)
}
```

**响应示例**:
- 登录成功:
```json
{
  "code": 200,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

- 登录失败:
```json
{
  "code": 400,
  "msg": "用户不存在"
}
```

```json
{
  "error": "用户名或密码错误"
}
```

**错误码说明**:
- 200: 登录成功
- 400: 用户不存在或参数错误
- 401: 用户名或密码错误
- 500: 服务器内部错误

### 1.2 用户注册接口

**接口名称**: 用户注册
**请求方法**: POST
**请求URL**: `/api/auth/register`
**请求参数**:
```json
{
  "username": "newuser",   // 用户名(必填，长度3-20)
  "password": "password123", // 密码(必填，长度6-30)
  "email": "user@example.com" // 邮箱(必填，格式正确)
}
```

**响应示例**:
- 注册成功:
```json
{
  "message": "注册成功",
  "user": {
    "id": 1,
    "username": "newuser",
    "email": "user@example.com"
  }
}
```

- 注册失败:
```json
{
  "error": "用户名已存在"
}
```

```json
{
  "error": "邮箱已被注册"
}
```

```json
{
  "error": "请求参数无效: xxx"
}
```

**错误码说明**:
- 200: 注册成功
- 400: 请求参数无效
- 409: 用户名或邮箱已存在
- 500: 服务器内部错误

## 2. K8s相关接口（需Token认证）

### 2.1 获取指定命名空间下的控制器资源

**接口名称**: 获取控制器资源
**请求方法**: GET
**请求URL**: `/api/k8s/namespaces/:namespace/controllers`
**请求参数**:
- Path参数: `namespace` - 命名空间名称

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "namespace": "default",
  "deployments": [
    {
      "name": "nginx-deployment",
      "replicas": 3,
      "images": ["nginx:latest"],
      "ready": 3,
      "updated": 3,
      "available": 3,
      "created_at": "2023-10-01T12:00:00Z",
      "update_at": "2023-10-01T12:00:00Z",
      "port": 80
    }
  ],
  "statefulsets": [
    {
      "name": "mysql-statefulset",
      "replicas": 1,
      "images": ["mysql:8.0"],
      "ready": 1,
      "updated": 1,
      "available": 1,
      "created_at": "2023-10-01T12:00:00Z",
      "update_at": "2023-10-01T12:00:00Z",
      "port": 3306
    }
  ],
  "daemonsets": [
    {
      "name": "fluentd-daemonset",
      "images": ["fluentd:v1.14"],
      "desired": 3,
      "current": 3,
      "ready": 3,
      "updated": 3,
      "available": 3,
      "created_at": "2023-10-01T12:00:00Z",
      "update_at": "2023-10-01T12:00:00Z",
      "port": 24224
    }
  ]
}
```

**错误码说明**:
- 200: 请求成功
- 400: 参数错误
- 401: 未授权（Token无效或过期）
- 500: 服务器内部错误

### 2.2 获取指定命名空间下的Pod资源

**接口名称**: 获取Pod资源
**请求方法**: GET
**请求URL**: `/api/k8s/:namespace/:resourcename/pods`
**请求参数**:
- Path参数: `namespace` - 命名空间名称
- Path参数: `resourcename` - 资源名称

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "resourcename": "nginx-deployment",
  "namespace": "default"
}
```

**错误码说明**:
- 200: 请求成功
- 400: 参数错误
- 401: 未授权（Token无效或过期）

### 2.3 获取所有节点资源

**接口名称**: 获取所有节点
**请求方法**: GET
**请求URL**: `/api/k8s/get/nodes`

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
[
  {
    "name": "node-1",
    "status": "Ready",
    "taints": [],
    "role": "master",
    "cpu": "2",
    "memory": "4Gi"
  },
  {
    "name": "node-2",
    "status": "Ready",
    "taints": [],
    "role": "worker",
    "cpu": "4",
    "memory": "8Gi"
  }
]
```

**错误码说明**:
- 200: 请求成功
- 401: 未授权（Token无效或过期）
- 500: 服务器内部错误

### 2.4 获取节点详细信息

**接口名称**: 获取节点详细信息
**请求方法**: GET
**请求URL**: `/api/k8s/get/nodename`

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
[
  {
    "name": "node-1",
    "status": "Ready",
    "role": "master",
    "ip": "192.168.1.100",
    "osType": "linux",
    "osVersion": "Ubuntu 20.04.4 LTS",
    "kubeletVersion": "v1.24.0",
    "kube-proxy": "v1.24.0",
    "dockerVersion": "docker://20.10.16",
    "coreVersion": "5.4.0-124-generic",
    "nodecreatetime": "2023-10-01T12:00:00Z",
    "taints": []
  }
]
```

**错误码说明**:
- 200: 请求成功
- 401: 未授权（Token无效或过期）
- 500: 服务器内部错误

### 2.5 获取容器组

**接口名称**: 获取容器组
**请求方法**: GET
**请求URL**: `/api/k8s/get/container_group/:namespace`
**请求参数**:
- Path参数: `namespace` - 命名空间名称

**请求头**:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**:
```json
{
  "pod": {
    "metadata": {
      "name": "nginx-pod",
      "namespace": "default",
      "labels": {
        "app": "nginx"
      }
    },
    "spec": {
      "containers": [
        {
          "name": "nginx",
          "image": "nginx:latest",
          "ports": [
            {
              "containerPort": 80
            }
          ]
        }
      ]
    },
    "status": {
      "phase": "Running",
      "conditions": [
        {
          "type": "Ready",
          "status": "True"
        }
      ]
    }
  }
}
```

**错误码说明**:
- 200: 请求成功
- 401: 未授权（Token无效或过期）
- 500: 服务器内部错误

## 2. 通用说明

### 2.1 认证Token

- 除了登录和注册接口外，所有K8s相关接口都需要在请求头中携带Token
- Token格式: `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- Token有效期: 默认24小时

### 2.2 错误处理

所有接口的错误响应格式统一为:
```json
{
  "error": "错误描述"
}
```

或者:
```json
{
  "code": 错误码,
  "msg": "错误描述"
}
```

### 2.3 请求头

- Content-Type: application/json
- Authorization: Bearer Token（需要认证的接口）

## 3. 技术栈说明

- 后端框架: Gin (Go语言)
- 数据库: MySQL + Redis
- K8s客户端: client-go
- 认证: JWT Token

## 4. 接口变更记录

| 日期 | 变更内容 | 版本 |
|------|---------|------|
| 2023-10-01 | 初始版本 | v1.0 |
| 2023-10-15 | 添加容器组接口 | v1.1 |
| 2023-10-20 | 优化节点信息接口 | v1.2 |