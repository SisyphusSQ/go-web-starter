# changeLog

### Unreleased

### v2.0.1(20261001)

#### optimization:

1. 明确 controller、service、repository 与 DO/DTO/VO 按业务域分子包，层根包保留装配和明确公共能力；补充各层 Agent 入口和新增模块步骤。
2. 明确 MySQL DO 业务子包内一表一文件，DO、TableName() 与列映射同文件，禁止多表模型集中平铺。
3. 示例用户 VO 移至 models/vo/example_vo，Lark DTO 移至 models/dto/lark_dto，同步调用方和组件过滤规则。

#### note:

1. 已有生成工程需自行同步规范并迁移相关 package/import；HTTP、JSON、表名及鉴权契约不因目录迁移而改变。
2. 本轮开发验证及未覆盖范围见 [业务子包验证记录](docs/test/2026-10-01-business-packages.md)；本版本无 SQL 变更。

### v2.0.0(20260930)

#### feature:

1. 新增 docs/sqls/schema、unreleased、releases 的 SQL 文件规范，明确版本归档及执行证据要求。

2. 新增 Windows、macOS、Linux 的 amd64/arm64 六平台构建及独立产物目录；保留本机构建入口。
3. 文档按 docs/design/architecture 与 docs/design/details/<topic> 分层，固化 DO/DTO/VO、Go 代码规范和新增模块步骤；默认工程保留模型目录说明。

4. 默认生成无数据库、无业务示例的工程；新增 --with、--examples 和 Harness Issue Provider 配置。
5. 新增 reference 比较和同步命令、生成来源与文件摘要清单，拒绝覆盖参考工程的本地修改。
6. 预置 Harness v0.7.0、docs、隔离联调 Compose 和产物组合验证流程。
#### optimization:


1. 保持 Go 1.27.1、Echo v5 和 Uber Fx，更新保留依赖及间接依赖到核对的稳定版本。
2. 组件按配置装配；新增请求关联日志、统一错误、健康与就绪探针及同库事务入口。
3. JWT 与 Redis 显式关联，metrics 受鉴权保护；完整 CRUD 仅随 --examples 生成。

#### note:

- 模块路径采用 /v2；version 与 --version 使用同一版本来源，正式产物固定注入 v2.0.0。

1. 本基底用于新项目初始化；生成默认值、组件配置和示例开关以当前文档为准。
2. 不附带业务或 examples SQL；docs/sqls 预置完整结构、unreleased 和版本归档规范。

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
