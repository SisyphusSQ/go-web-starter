# go-web-starter Agent Guide

- 本仓维护 Go 1.27.1 / Echo v5 / Uber Fx 生成器；模板权威来源为 internal/scaf_fold/_template。
- app/service 等运行能力修改必须落在模板，再通过固定 reference 配置同步 go-starter。
- 构建：make build；生成器测试：make test；生成矩阵：make integration；只读格式/vet/build：make verify。
- 模板中的 *_test.go.tmpl 是生成产物行为测试，生成器本身通过不会替代产物测试。
- reference --target /abs/path 默认只比较；--write 仅更新清单管理的文件，拒绝覆盖本地修改。法律文件和 changeLog 由各仓维护。
- 依赖更新需覆盖根 go.mod、模板 go.mod 和生成解析出的间接依赖；普通生成不联网解析 @latest。
- 构造器不访问外部服务；Fx OnStart 打开资源，失败 hook 清理自身资源，OnStop 释放成功启动的资源。
- 新业务不默认进入基底；组件通过生成选项和运行配置显式选择；JWT 要求 Redis。
- 文档见 docs/README.md；当前跨仓计划见 .agents/plans/2026-09-30-starter-modernization.md。
- Harness v0.7.0 来源 b20e5e8ece6a529c7d74aa0a8b1de77bd06c374c，模板内嵌且保留来源说明，通用技能使用共享插件。
- .agents/runs 和 state 的真实记录、凭据与日志不提交；模板可提交。
- 本地验证、生成产物测试、Docker 和真实环境验收分别报告。提交/发版收尾不重复测试；推送、合并与发布以当前用户授权为准。

## 开发与文档入口

- 修改前阅读 docs/design/architecture/README.md 和 docs/design/details/development/README.md。
- 生成工程模型与代码规则的权威文件位于 internal/scaf_fold/_template/docs/design/architecture/models.md.tmpl 和 internal/scaf_fold/_template/docs/design/details/development/code-style.md.tmpl。
- 业务子包权威规则位于 internal/scaf_fold/_template/docs/design/architecture/packages.md.tmpl；模板中的业务 controller/service/repository、DO/DTO/VO 必须按域分包，MySQL DO 一个表一个文件，层根包仅保留装配及明确公共能力。
- 不生成 examples 也必须保留 models/do、dto、vo 的说明与规范入口；不要用恢复演示业务代替架构说明。
- docs/README.md 只做导航；architecture 放长期边界，details/<topic>/ 放实现细节，test 放验收；调整目录时同步所有链接。
- make build-all / release-all 覆盖 Windows/macOS/Linux 的 amd64 和 arm64；详见 docs/design/details/build/README.md。

## SQL 文件规范

新增持久化模型或交付 SQL 前阅读 docs/sqls/README.md。当前完整结构放 docs/sqls/schema，待发布 SQL 放 docs/sqls/unreleased，真实发布时归档到 docs/sqls/releases/vX.Y.Z。无实际业务 SQL 时只保留说明，不生成 examples DDL、种子数据、空 SQL 或虚构版本目录；已发布文件不可改写。归档不等于数据库已经执行，执行证据单独记录。
