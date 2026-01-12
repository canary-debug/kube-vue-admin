# kube-vue-admin 后端部署文档

## 目录

1. [概述](#概述)
2. [环境要求](#环境要求)
3. [Docker 镜像构建](#docker-镜像构建)
4. [Kubernetes 部署准备](#kubernetes-部署准备)
5. [Kubernetes 部署步骤](#kubernetes-部署步骤)
6. [部署验证](#部署验证)
7. [常见问题](#常见问题)

## 概述

本文档描述了如何使用 Docker 和 Kubernetes 部署 kube-vue-admin 后端服务。部署流程包括 Docker 镜像构建和 Kubernetes 资源部署两个主要阶段。

## 环境要求

- Docker 19.03+ 
- Kubernetes 1.18+ 
- kubectl 命令行工具
- Go 1.24+ (仅用于源码构建)

## Docker 镜像构建

### 1. 准备构建环境

确保当前目录在项目根目录下：

```bash
git clone https://github.com/canary-debug/kube-vue-admin.git
cd kube-vue-admin
```

### 2. 构建 Docker 镜像

使用示例目录中的 Dockerfile 构建镜像：

```bash
# 复制 Dockerfile 到项目根目录
cp example/backend/Dockerfile ./

# 构建镜像
docker build -t kube-vue-admin:v1.0 .
```

### 3. 验证镜像

构建完成后，验证镜像是否成功创建：

```bash
docker images | grep kube-vue-admin
```

### 4. 推送镜像到仓库（可选）

如果需要将镜像推送到远程仓库，执行以下命令：

```bash
# 标记镜像
docker tag kube-vue-admin:v1.0 <registry-url>/kube-vue-admin:v1.0

# 推送镜像
docker push <registry-url>/kube-vue-admin:v1.0
```

## Kubernetes 部署准备

### 1. 确保 Kubernetes 集群可用

验证 kubectl 能否正常连接到 Kubernetes 集群：

```bash
kubectl cluster-info
```

### 2. 创建命名空间

创建用于部署应用的命名空间：

```bash
kubectl create ns kube-vue-admin
```

### 3. 创建 ConfigMap 存储 kubeconfig

将 kubeconfig 文件创建为 ConfigMap，以便 Pod 能够访问 Kubernetes 集群：

```bash
# 使用当前用户的 kubeconfig
kubectl -n kube-vue-admin create configmap kube-config --from-file=config=$HOME/.kube/config

# 或使用指定路径的 kubeconfig
kubectl -n kube-vue-admin create configmap kube-config --from-file=config=/path/to/your/kubeconfig
```

## Kubernetes 部署步骤

### 1. 修改部署配置文件

编辑 `example/backend/k8s-deploy.yaml` 文件，根据实际情况修改以下内容：

- `image`: 替换为你的 Docker 镜像地址
- `nodeName`: 替换为你想要部署的节点名称（可选）
- 环境变量：根据实际情况修改数据库和 Redis 连接信息
- `nodePort`: 根据需要修改服务暴露的节点端口（可选）

### 2. 部署应用

使用 kubectl 部署应用：

```bash
kubectl apply -f example/backend/k8s-deploy.yaml
```

### 3. 查看部署状态

检查 Deployment 状态：

```bash
kubectl -n kube-vue-admin get deployment kube-vue-backend
```

检查 Pod 状态：

```bash
kubectl -n kube-vue-admin get pods -l app=kube-vue-backend
```

检查 Service 状态：

```bash
kubectl -n kube-vue-admin get service kube-vue-backend-service
```

## 部署验证

### 1. 访问服务

使用以下方式访问部署的服务：

```bash
# 使用 NodePort 访问
curl http://<node-ip>:30080

# 预期响应
# {"message": "pong"}
```

### 2. 查看 Pod 日志

查看应用运行日志，确认服务正常启动：

```bash
kubectl -n kube-vue-admin logs -f <pod-name>
```

## 常见问题

### 1. Pod 启动失败

检查 Pod 事件和日志：

```bash
kubectl -n kube-vue-admin describe pod <pod-name>
kubectl -n kube-vue-admin logs <pod-name>
```

### 2. 服务无法访问

检查 Service 配置和节点端口：

```bash
kubectl -n kube-vue-admin describe service kube-vue-backend-service
```

### 3. 数据库连接失败

确保环境变量中的数据库连接信息正确，并验证数据库服务是否可用。

### 4. Kubernetes API 访问失败

检查 ConfigMap 中的 kubeconfig 是否正确，以及 Pod 是否能够访问 Kubernetes API 服务器。

## 附录：部署架构图

```
┌───────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster                   │
│ ┌─────────────────────────────────────────────────────┐ │
│ │                      Namespace                      │ │
│ │                kube-vue-admin                       │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │                   Deployment                   │ │ │
│ │ │              kube-vue-backend                  │ │ │
│ │ │ ┌─────────────────┐ ┌─────────────────┐        │ │ │
│ │ │ │      Pod        │ │      Pod        │        │ │ │
│ │ │ │ ┌─────────────┐ │ │ ┌─────────────┐ │        │ │ │
│ │ │ │ │   Container │ │ │ │   Container │ │        │ │ │
│ │ │ │ │ kube-vue-admin │ │ │ kube-vue-admin │ │        │ │ │
│ │ │ │ └─────────────┘ │ │ └─────────────┘ │        │ │ │
│ │ │ │   Port: 9000     │ │   Port: 9000     │        │ │ │
│ │ │ └─────────────────┘ └─────────────────┘        │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │                   Service                      │ │ │
│ │ │          kube-vue-backend-service              │ │ │
│ │ │    Type: NodePort, Port: 9000, NodePort: 30080 │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ │ ┌─────────────────────────────────────────────────┐ │ │
│ │ │                   ConfigMap                     │ │ │
│ │ │                  kube-config                    │ │ │
│ │ │       Contains kubeconfig for cluster access    │ │ │
│ │ └─────────────────────────────────────────────────┘ │ │
│ └─────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────┘
```

## 更新日志

| 日期       | 版本 | 更新内容               |
|------------|------|------------------------|
| 2026-01-12 | v1.0 | 初始版本部署文档       |
