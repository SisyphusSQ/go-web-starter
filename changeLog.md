# changeLog

### v1.1.0(20260911)
#### optimization:
1. 工具与生成项目统一升级到 Go 1.27.1；生成器固定稳定版本并在写盘前格式化 Go 源码，避免因本机 toolchain 不同产生漂移。
2. 依赖统一升级到当前稳定 module path，完成 Echo v5、MongoDB Driver v2、`lumberjack.v2` 与 `go-hashids/v2` 迁移，并移除无稳定 tag 的 `golib`、`qmgo` 和第三方 UUID 依赖。
3. 配置改为独立 Viper 实例、严格解码、duration 字符串和 `APP_` 环境变量覆盖；示例配置不再内置密码或 JWT secret。
4. MySQL、MongoDB、Redis、HTTP 与 cron 接入 Fx 生命周期，启动阶段校验连通性或监听地址，关闭阶段释放资源。
5. Docker 镜像升级到 Alpine 3.24.1 并使用非 root 用户；新增 `.dockerignore`、健康检查、GitHub Actions CI 与 Dependabot 配置。

#### bugFix:
1. 修复 GORM 默认连接参数只修改副本、日志级别设置未生效、更新与删除不存在记录仍返回成功的问题。
2. 修复 cron provider 未被消费导致任务不启动、锁 TTL 过短且 key 未按应用隔离的问题。
3. 修复 AES-CBC 固定 IV、非法密文可能触发 panic，以及并发工具共享错误变量产生数据竞争的问题。
4. 收紧 CORS、Basic/AK 常量时间比较、JWT 算法与 issuer 校验、请求体大小和外部 HTTP 超时，避免泄露查询参数或把业务失败记录为成功。

#### note:
1. 本次变更影响后续新生成的项目，不会自动迁移已有项目；新模板仅在 debug 模式允许 `key.type: none`，非 debug 环境必须配置明确鉴权与允许的 CORS Origin。

### v1.0.0(20260216)
#### feature:
1. 发布 `go-web-starter` 1.0.0，提供 `new`/`init` 一键生成 Go Web 项目能力，开箱即可启动基础工程。
2. 支持多数据库模板按需生成，覆盖 `mysql`、`mongodb` 与双库组合三种场景，适配不同项目起步方式。
3. 完成模板参数化体系，统一支持模块名、二进制名等关键字段渲染，并提供 `version` 命令输出构建信息。
4. 完善跨库能力，`LarkService` 模板解除 MySQL 强依赖，mongodb-only 场景可直接生成并正常编译。
5. 强化工程可用性与稳定性，补全 `prometheus`/`lark` 默认配置、统一 Lark SDK 依赖管理，并补齐关键边界与集成验证。
