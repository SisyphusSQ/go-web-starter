# 两仓基底现代化实施记录

## 已确认范围

- 继续使用 Uber Fx 注入与生命周期；保留 app/cmd 和 controller/service/repository、models/do/dto/vo 分层。
- go-web-starter 为模板权威来源；go-starter 是固定配置生成的可克隆参考工程，默认外部组件关闭。
- Go 1.27.1，全部保留依赖升级至核对的稳定版本，核对大版本 API。
- 默认无数据库、无 Redis 也能启动；生成选项控制代码包含，配置控制组件启用。
- 清理默认 User CRUD、示例 SQL 和 cron；完整示例通过显式选项保留并验证。
- 保留 Basic/AK，可选 JWT+Redis；登录登出与 User 业务属于示例，不让删除示例削弱鉴权。
- 请求关联日志、统一错误和参数校验、存活/就绪、新环境初始化、同库事务、失败清理、隔离联调。
- 两仓及生成产物预置 Harness v0.7.0 和完整 docs；共享技能不重复分发。
- 不改已有业务项目，不部署；模块生成、OpenAPI、完整 OTel 后续另做。

## 基线与分支

- 当前开发机器 SQMC04；两仓分支 suqing/starter-modernization。
- go-web-starter 基线 bcbc29c；go-starter 基线 5496e63。
- Harness v0.7.0：b20e5e8ece6a529c7d74aa0a8b1de77bd06c374c，来源 SisyphusSQ/harness-template。
- 所有测试只使用本地临时目录/隔离资源；禁止连接业务数据库或发送真实飞书消息。

## 执行顺序和进度

1. [x] 生成选项、来源记录、无数据库默认、示例过滤和回归测试。
2. [x] 组件配置、Fx 装配和生命周期；鉴权、日志、错误、探针。
3. [x] SQL 文件规范与事务、可选 SDK、依赖升级。
4. [x] docs、Harness、Compose、CI、模板维护与参考仓同步。
5. [x] 四种数据库组合和组件组合验证、运行/失败/退出场景、自查与修复。
6. [ ] 环境验收剩余项：Docker、真实 MongoDB、MySQL 8.4 容器及实际 CI，当前缺少运行环境且尚未 push。

## 验收口径

- 生成器和生成项目的 build/test/race/vet 分开记录。
- 无数据库、MySQL、MongoDB、双库；Redis 开关、JWT 依赖约束、示例显式选择。
- 默认启动无外部连接；配置错误/端口冲突返回非零；鉴权拒绝、错误格式、请求 ID、ready 和退出准确。
- 新环境按当前表结构初始化；同库事务提交/回滚/取消可验证。
- Docker、lint、govulncheck 缺工具时如实报告，不能把未执行写成通过。
- 两仓同步可复现；已有本地文件不被静默覆盖；生成器拒绝非空目录和符号链接。

## 记录

- 2026-09-30：两仓干净 main 已各建分支；本机 Go 1.27.1 可用，未发现 Docker/golangci-lint/govulncheck。Modern Go 1.27 指南已读取。

- 生成边界三项回归先失败后通过；配置禁用组件回归先失败后通过；全组件模板编译/单元测试通过。404 使用 Echo.StatusCode 保持 v5 错误契约。完整示例仍在进一步验证。


## 当前结果

代码与文档已实现；真实 MySQL 8.0.46、Redis 8.6.2 和 JWT 撤销验证通过，临时资源已清理。两仓尚未提交或推送；Docker/MongoDB/CI 未执行，不能算完成这些验收。

完整范围、修复记录与验证证据见 docs/test/2026-09-30-modernization.md。之前记录中的“进一步验证”状态以该结果为准。


## 用户追加要求：六平台构建与开发规范

- Makefile 增加 windows/darwin/linux × amd64/arm64，Windows 使用 .exe，各组合分目录输出；build/release 不捆绑测试或发布。
- 保留无 examples 工程中的 models/do、dto、vo 说明；根 AGENTS、模型 AGENTS 和设计文档共同明确职责、依赖与新增业务步骤。
- 参考 tuyu-studio 的 architecture/details 和 codex-pulse 的主题索引，文档整理至 docs/design；同步索引与原有引用。
- 验收区分六平台交叉构建、产物格式和目标机器运行；不宣称 Windows/Linux 或其他架构已实机运行。

追加要求已落地，两仓六平台共 12 个交叉编译产物通过架构核对；生成/链接/同步检查通过。结果见 docs/test/2026-09-30-platforms-and-conventions.md。实机平台运行与远端 CI 未执行。


## 用户最终 SQL 约定

- 基底只维护代码及 SQL 文件规范，不携带业务/示例 DDL 或种子数据。
- docs/sqls/schema 保存后续实际业务完整结构；unreleased 保存待交付 SQL；真实版本确定后归档至 releases/vX.Y.Z，不预造版本目录。
- 参考 SQMC03 的 BaaS / DBIO SQL 组织原则，已发布文件不可变，版本归档与环境执行证据分开。

SQL 最终范围已实现；生成/构建/lint/文档链接与参考同步通过，详见 docs/test/2026-09-30-sql-conventions.md。本轮未执行真实数据库操作。

## v2.0.0 发布执行

用户已授权分别提 PR、合并 main、构建六平台发布包并依次发布生成器和参考工程。收尾不重复测试；保留已有证据和未覆盖项，完成后读取远端 PR、tag 和 Release 状态。
