# 模板维护和参考仓同步

所有公共行为修改 internal/scaf_fold/_template，再验证生成器及生成工程。

固定 go-starter 参考配置：module=github.com/SisyphusSQ/go-starter/v2，appVersion=v2.0.0，binary=go-starter，MySQL+MongoDB+Redis+cron+Lark+Prometheus query+JWT，examples=false，issue-provider=linear。所有外部组件默认 disabled。

```sh
go run . reference --target /abs/path/to/go-starter
go run . reference --target /abs/path/to/go-starter --write
```

比较过程在系统临时目录生成工程并执行 go mod tidy，因此可能访问模块代理和写 Go 缓存，不修改目标业务文件。--write 更新管理清单内的文件；已有本地修改或与未管理文件碰撞会整体拒绝。清单文件 .starter.json 记录版本摘要、Harness 来源、选项及逐文件 SHA-256。LICENSE、changeLog.md 和未管理文件由各仓独立维护。

首次接管旧参考工程需要人工审阅全部已跟踪文件并建立初始清单；不提供对任意现存业务仓库的强制接管入口。基底只用于新项目初始化，已有业务项目保持独立。

依赖升级时同时检查根 go.mod、go.mod.tmpl、生成工程解析出的间接依赖和对应 API。确定版本后固定到模板；不把 latest 写成运行时解析策略。执行 make integration 验证组合，随后同步参考仓并复核差异。

来源清单额外记录 generator_version；生成业务工程默认版本为 dev，不继承基底发布版本。
