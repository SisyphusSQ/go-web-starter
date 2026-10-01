# 验证范围

- make test：生成器单元测试、输出保护、参数组合与同步保护。
- make integration：联网解析并验证无数据库、MySQL、MongoDB、双库、可选组件和完整示例；对产物运行 race/test/vet/build，并编译 integration 标签。
- 生成工程的 app/cmd 测试实际验证无外部依赖启动、Fx 构造阶段不连接外部服务、监听冲突返回失败。
- internal/http 测试覆盖 request_id、统一错误、鉴权保护、生命周期就绪状态。
- SQL 事务测试在数据库驱动边界验证语句、提交、回滚和取消。
- Docker 与真实隔离数据库使用生成工程 docs/test/integration.md，未提供环境不能算通过。

测试日志保存本地 .agents/runs/，脱敏的结果摘要放 docs/test/。不把仅编译 integration 标签写成真实数据库验收。

## 本次验证记录

- [业务子包规范 v2.0.1](2026-10-01-business-packages.md)
- [首轮现代化](2026-09-30-modernization.md)
- [六平台与开发规范](2026-09-30-platforms-and-conventions.md)

- [SQL 文件规范与最终范围](2026-09-30-sql-conventions.md)
