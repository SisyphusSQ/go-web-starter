# go-web-starter

Go 1.27.1 + Echo v5 + Uber Fx 脚手架。当前发布：v2.0.0，Go module 路径使用 /v2。生成工程自带严格配置校验、请求关联日志、统一错误、鉴权、存活/就绪探针、资源生命周期、文档和 Harness。

```sh
make build
./bin/go-web-starter new demo
cd demo
go mod tidy
make run
```

默认无数据库、无 Redis、无示例业务；监听本机且仅 debug 允许无鉴权。已选入的组件默认关闭，运行时在配置里开启。

```sh
go-web-starter new my-service --module example.com/my-service --db mysql --with redis,cron
go-web-starter new account-example --db mysql --with redis,jwt --examples
go-web-starter init --module example.com/my-service --db none
```

| 选项 | 说明 |
|---|---|
| --db | none（默认）、mysql、mongodb、mysql,mongodb |
| --with | redis、cron、lark、prometheus-query、jwt；jwt 要求 redis |
| --examples | 显式生成 User CRUD，至少选择一种数据库 |
| --module / --binary | 自定义模块路径和二进制名 |
| --issue-provider / --issue-prefix | Harness 项目元数据；仅 repo provider 生成 docs/issues |

new/init 只接受不存在、空目录或仅有 .git 的目录，拒绝符号链接和覆盖。生成前完成所有模板渲染及 Go 格式检查，避免模板错误留下半成品。磁盘写入失败仍可能部分完成，应检查输出后恢复。

[维护文档](docs/README.md) · [组件与测试](docs/test/README.md)

make test 验证生成器；make integration 验证实际生成工程；make build/release 只构建。go-starter 参考仓通过固定配置和来源清单同步，见维护文档。

安装此版本：`go install github.com/SisyphusSQ/go-web-starter/v2@v2.0.0`。也可下载对应平台的 Release 压缩包；Windows 使用 .exe。
