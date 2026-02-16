# changeLog

### v1.0.0(20260216)
#### feature:
1. 发布 `go-web-starter` 1.0.0，提供 `new`/`init` 一键生成 Go Web 项目能力，开箱即可启动基础工程。
2. 支持多数据库模板按需生成，覆盖 `mysql`、`mongodb` 与双库组合三种场景，适配不同项目起步方式。
3. 完成模板参数化体系，统一支持模块名、二进制名等关键字段渲染，并提供 `version` 命令输出构建信息。
4. 完善跨库能力，`LarkService` 模板解除 MySQL 强依赖，mongodb-only 场景可直接生成并正常编译。
5. 强化工程可用性与稳定性，补全 `prometheus`/`lark` 默认配置、统一 Lark SDK 依赖管理，并补齐关键边界与集成验证。
