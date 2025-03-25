# XCode

XCode 后端代码采用 Gin 框架作为 Web 后端，使用 GRPC 微服务架构实现业务逻辑，完成用户注册、登陆、鉴权、搜索一系列功能。

## 服务启动方式

按照以下顺序启动服务：

1. 启动 `auth-service`
2. 启动 `search-service`
3. 最后启动 `gateway-service`

### auth-service

`auth-service` 负责用户注册、登录和鉴权功能。启动方式如下：

```bash
go build
./auth-service -redishost <REDIS_HOST> -redisport <REDIS_PORT> -redispass <REDIS_PASSWORD> -mysqlhost <MYSQL_HOST> -mysqlport <MYSQL_PORT> -mysqluser <MYSQL_USER> -mysqlpassword <MYSQL_PASSWORD>
```
参数说明：
- `<REDIS_HOST>`: Redis 服务器地址
- `<REDIS_PORT>`: Redis 服务器端口
- `<REDIS_PASSWORD>`: Redis 服务器密码
- `<MYSQL_HOST>`: MySQL 服务器地址
- `<MYSQL_PORT>`: MySQL 服务器端口
- `<MYSQL_USER>`: MySQL 用户名
- `<MYSQL_PASSWORD>`: MySQL 用户密码
### search-service

`search-service` 负责搜索功能。启动方式如下：

```bash
go build
./search-service -port <PORT> -esaddr <ES_ADDRESS> -esport <ES_PORT>
```

参数说明：
- `<PORT>`: 服务监听端口
- `<ES_ADDRESS>`: Elasticsearch 服务器地址
- `<ES_PORT>`: Elasticsearch 服务器端口

### gateway-service

`gateway-service` 作为网关服务，负责请求路由和负载均衡。启动方式如下：

```bash
go build
./gateway-service
```

## 项目结构

```
.
├── README.md
├── auth-service
│   ├── auth_service.go
│   ├── entity
│   │   └── config
│   │       └── cfg.go
│   ├── main.go
│   └── repo
│       ├── userdb
│       │   └── api.go
│       └── usertoken
│           └── api.go
├── gateway-service
│   ├── config
│   │   └── config.go
│   └── main.go
├── go.mod
├── go.sum
├── proto
│   ├── auth
│   │   ├── auth.pb.go
│   │   ├── auth.proto
│   │   └── auth_grpc.pb.go
│   ├── auth.proto
│   ├── search
│   │   ├── search.pb.go
│   │   ├── search.proto
│   │   └── search_grpc.pb.go
│   └── search.proto
└── search-service
    ├── entity
    │   └── config
    │       └── cfg.go
    ├── main.go
    ├── repo
    │   └── searchdb
    │       └── api.go
    └── search_service.go
```

## 依赖安装

在启动服务之前，请确保已经安装了所有依赖。你可以使用以下命令安装依赖：

```bash
go mod tidy
```

## 环境配置
请根据实际情况配置环境变量或配置文件，以确保服务能够正确连接到 Redis、MySQL 和 Elasticsearch。