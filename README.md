# 项目介绍
一款仿 TikTok 的视频网站服务端，包含视频流，视频上传，好友管理和聊天等特色。设计采用拆分服务的微服务架构，分为用户、视频、互动、社交、聊天、群组等子模块。

# 技术栈
- web：hertz by cwgo
- 微服务：kitex by cwgo
- 数据库：MySQL
- 缓存：Redis
- 消息队列：RabbitMQ
- 服务注册与发现：etcd
- 快速查询：ElasticSearch
- 监控指标拉取：Prometheus & Exporters
- 链路追踪：jaeger
- 可观测：grafana
- CI/CD: Github Action
- 单测：mockey + GoConvey

# 项目结构
```txt
.
├── api                  # 网关
│   ├── chat               # 聊天服务
│   ├── handler            # http
│   ├── middleware         # 中间件
│   ├── model
│   ├── pack               # resp封装
│   ├── router             # http路由
│   └── rpc                # rpc封装
├── cmd                  # 程序入口
├── config               # 配置
├── docker               # docker相关
├── docs                 # 文档
├── idl                  # 接口定义
├── internal             # 核心业务逻辑&基础设施
│   ├── chat               # 聊天
│   ├── group              # 群组
│   ├── interaction        # 互动
│   ├── social             # 社交
│   ├── user               # 用户
│   └── video              # 视频
│   ├── pkg                # infra
│   │   ├── base
│   │   ├── cache            # 缓存
│   │   ├── constants        # 宏
│   │   ├── db               # db
│   │   ├── domain           # 业务层obj
│   │   ├── es               # elasticsearch
│   │   ├── mq               # 消息队列
│   │   └── pack             # 业务层obj封装
├── kitex_gen            # kitex生成
├── pkg                  # 公有包
│   ├── mq
│   └── util
│       └── jwt
└── uploads              # 保存上传文件
    └── static
        └── videos
```

# 构建 & 部署
## 本地部署
### 环境准备
- Docker
    - 参照Docker文档安装 Docker
    - 用于启动 MySQL、Redis 等基础设施

### 启动基础设施
```bash
make env-up
```

### 启动 api网关
```bash
make api
```
### 启动必要的服务
```bash
make <service>  # <service>为要启动的服务名
```
以下为支持的服务名：
- user
- video
- interaction
- social
- chat
- group

### （可选）关闭基础设施
```bash
make env-down
```

# 技术细节
## 子模块架构
见 [子模块架构](./docs/arch.md)。

# 传送门
- apifox接口文档：https://hhhb4f6wui.apifox.cn/
- 单测踩坑和心得：[Golang 测试踩坑与心得](https://becks723.github.io/2026/03/14/golang%E6%B5%8B%E8%AF%95%E8%B8%A9%E5%9D%91/)

# 不足 & 规划
- http路由应统一加 `/api/v1` 前缀
- swagger 文档
- 自己应维护一套错误码系统 errno
- 日志，包括业务中返回err太多，
- 单测覆盖率
- 参数校验不应耦合在业务逻辑中
