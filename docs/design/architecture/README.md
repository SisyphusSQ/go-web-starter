# 生成器与参考工程架构

- cmd 解析 CLI 参数；internal/scaf_fold 校验选项、渲染、格式化和受保护地写入文件。
- internal/scaf_fold/_template 是运行工程的权威来源。默认工程保留架构目录与说明，业务示例由 --examples 控制。
- go-starter 是固定配置参考工程；reference 通过逐文件摘要比较/同步，禁止覆盖本地业务修改或未知文件。
- 通用目录、模型规范、Makefile 与文档也属于模板契约，新增生成选项不能使必要协作入口消失。

运行工程的规范源文件：

- [分层与 Fx](../../../internal/scaf_fold/_template/docs/design/architecture/README.md.tmpl)
- [DO / DTO / VO](../../../internal/scaf_fold/_template/docs/design/architecture/models.md.tmpl)
- [Go 代码规范](../../../internal/scaf_fold/_template/docs/design/details/development/code-style.md.tmpl)

生成器自身不采用业务 DO/DTO/VO；CLI/渲染配置留在现有包，修改模板中的业务代码必须服从上述分层。
