# gin-template

`gin-template` 是基于 [Gin](https://github.com/gin-gonic/gin) 进行模块化设计的 API 框架，封装了大部分常用的功能，致力于进行快速的业务研发。

本项目基于 [project-layout](https://github.com/golang-standards/project-layout) 项目结构进行设计，参考了 [go-gin-api](https://github.com/xinliangnote/go-gin-api)

项目完成后会拆分为三部分，迷你版、基础版、完整版，迷你版只包含基础功能，基础版包含基础功能，完整版包含所有有用没用的功能。

完成后不一定可以进行完整的上线测试，所以该项目仅供参考，如需线上使用，请进行完整的上线测试。

## 功能特性

- [x] 支持 [rate](https://golang.org/x/time/rate) 接口限流
- [x] 支持 panic 异常时邮件/企业微信通知
- [x] 支持 [cors](https://github.com/rs/cors) 接口跨域
- [x] 支持 [Prometheus](https://github.com/prometheus/client_golang) 指标记录
- [x] 支持 [Swagger](https://github.com/swaggo/gin-swagger) 接口文档生成
- [ ] 支持 [GraphQL](https://github.com/99designs/gqlgen) 查询语言
- [x] 支持 trace 项目内部链路追踪
- [x] 支持 [pprof](https://github.com/gin-contrib/pprof) 性能剖析
- [x] 支持 errno 统一定义错误码
- [x] 支持 [zap](https://go.uber.org/zap) 日志收集
- [x] 支持 [viper](https://github.com/spf13/viper) 配置文件解析
- [x] 支持 [gorm](https://gorm.io/gorm) 数据库组件
- [x] 支持 [go-redis](https://github.com/go-redis/redis/v7) 组件
- [x] 支持 RESTful API 返回值规范
- [x] 支持 生成数据表 CURD、控制器方法 等代码生成器
- [x] 支持 [cron](https://github.com/jakecoffman/cron) 定时任务，在后台可界面配置
- [x] 支持 [websocket](https://github.com/gorilla/websocket) 实时通讯，在后台有界面演示
- [ ] 支持 多租户
- [ ] 支持 多语言
- [ ] 支持 [casbin](https://github.com/casbin/casbin) RBAC 访问控制
- [ ] 支持 多点登陆拦截
- [ ] 支持 Github 第三方登录
- [ ] 支持 微信第三方登录
- [ ] 支持 QQ第三方登陆
- [ ] 支持 微信支付
- [ ] 支持 支付宝支付
- [ ] Vue前端
- [ ] ...




