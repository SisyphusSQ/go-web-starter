# go-web-starter

`go-web-starter` 是一个 Go Web 脚手架工具，用于快速生成可运行的 Web
项目基础结构。工具本身基于 `cobra` 构建，模板通过 `embed.FS` 内置。

## 功能概览

- 提供 `new` 命令在目标目录生成新项目
- 提供 `init` 命令在当前目录初始化项目
- 支持 `--db` 按需选择数据库模板：
  - `mysql`
  - `mongodb`
  - `mysql,mongodb`（默认）
- 支持自定义模块名（`--module`）与二进制名（`--binary`）

## 环境要求

- Go `1.26.0` 或更高版本

## 构建与运行

```bash
go mod tidy
make build
```

构建产物默认在 `bin/go-web-starter`。

## 使用方式

### 1) 生成到新目录

```bash
go-web-starter new <output-dir> [flags]
```

### 2) 在当前目录初始化

```bash
go-web-starter init [flags]
```

`init` 仅允许当前目录为空或仅包含 `.git` 目录。

### 常用参数

- `-m, --module`：Go module 路径（默认 `example.com/<directory-name>`）
- `-b, --binary`：二进制名（默认从目录名推导）
- `--db`：数据库选择（`mysql` / `mongodb` / `mysql,mongodb`）

## 示例

```bash
# 生成 MySQL + MongoDB 双数据库项目（默认）
go-web-starter new demo-web

# 仅生成 MySQL 相关代码
go-web-starter new demo-web --db mysql

# 仅生成 MongoDB 相关代码
go-web-starter new demo-web --db mongodb

# 在当前目录初始化并指定模块名
go-web-starter init --module github.com/acme/demo-web --db mysql
```

## 生成后建议步骤

```bash
cd <output-dir>
go mod tidy
# 按需修改 config/config.yml
go run ./app/main.go http
```

## 开发与验证

```bash
# 默认测试（离线稳定）
go test ./...

# integration 测试（包含联网构建校验）
go test -tags integration ./internal/scaf_fold -run TestGenerateE2EDBCombosIntegration -count=1

# 构建与静态检查
go build ./...
go vet ./...
```
