# sgin
sgin是基于gin框架的一层封装，集成了一些常用的业务组件，可以基于此框架快速开发应用。

或者也可以把sgin框架放在业务服务的前面，这样子sgin也可以作为一个网关来使用。不过这个网关比较特殊，有一些业务组件，比如用户管理、权限管理等。这样子业务服务可以不用再去实现sgin已经有的基础功能服务，只需要实现自己的业务服务即可。

关于安全方面，业务服务配置只可以sgin访问即可。

![](doc/sgin.png)

### 为什么弄这个框架
- 1. 集成一些常用的业务组件
- 2. 为了快速开发
- 3. 为了学习

### 特性
- 1. 支持swagger文档生成
- 2. 集成gorm
- 3. 集成redis
- 4. jwt认证
- 5. 集成zap日志
- 6. 集成viper配置文件
- 7. 集成casbin权限管理
- 8. 集成邮件服务
- 9. 丰富的中间件（请求和响应日志hook、用户认证、签名校验、api请求权限等）
- 10. 路由转发
- 11. API接口限流

### 功能
- 1. 用户管理
- 2. 角色管理
- 3. 权限管理
- 4. 菜单管理
- 5. 邮件验证码
- 6. 文件上传
- 7. 团队管理
- 8. APP调用方管理
- 9. API接口权限管理



### swagger
[swagger操作文档](https://github.com/swaggo/swag/blob/master/README_zh-CN.md)

- 生成swagger文档
> swag init 

### 运行与配置
- 启动
```powershell
go build
./sgin.exe
```

- 关键环境变量
	- `SERVER_PORT`: 服务端口（如 8080）
	- `LOG_FILE`: 日志文件路径，对应 `LogConfig.Filename`
	- `MYSQL_HOST`/`MYSQL_PORT` 等：数据库连接
	- `ALLOWED_ORIGINS`: 允许跨域来源，逗号分隔（如 `https://foo.com,https://bar.com`）
	- `PASSWD_KEY`: 用于密码与 JWT 签名的密钥

- 健康检查
	- `GET /ping`: 基础存活检查
	- `GET /healthz`: 应用健康检查
	- `GET /readyz`: 就绪检查（会探测 DB/Redis 可用性）

### 安全与稳定性
- CORS: 增加白名单控制，仅回显允许的 `Origin`
- 日志脱敏: MySQL 连接信息在日志中隐藏密码
- 依赖更新: JWT 迁移到 `github.com/golang-jwt/jwt/v4`
- 服务稳态: 配置了 HTTP `Read/Write/Idle` 超时，降低慢连接风险
