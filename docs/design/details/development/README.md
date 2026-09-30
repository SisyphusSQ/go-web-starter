# 生成器开发规范

修改前阅读根 AGENTS.md、设计架构和相邻代码。使用 Go 1.27.1、gofmt，import 按标准库/第三方/本仓分组，导出类型写中文契约注释。

- CLI 负责参数与错误展示，scaf_fold 负责生成/同步，模板承担运行工程代码，避免三处重复实现同一行为。
- 保留错误链和上下文，不记录用户凭据；先渲染验证再写入，保持非空目录、符号链接和本地修改保护。
- 修改模型或业务分层先读模板 docs/design/architecture/models.md；不得在模板 service 中放业务 DTO/DO/VO。
- 文档调整同步索引、AGENTS、README 和相对链接；保留不带 examples 工程的 models/do、dto、vo 说明。
- make test 验证生成器，make integration 验证产物；跨平台修改执行 make build-all。build/release 本身不捆绑测试或发布。
- 构建、测试、目标平台运行、CI 和发布是不同证据；提交收尾复用已有证据。
