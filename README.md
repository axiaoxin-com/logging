# logging

[![Go Report Card](https://goreportcard.com/badge/github.com/axiaoxin-com/logging)](https://goreportcard.com/report/github.com/axiaoxin-com/logging)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/3f0bf6abb7504f2f8734f094bb65e0d6)](https://app.codacy.com/gh/axiaoxin-com/logging?utm_source=github.com&utm_medium=referral&utm_content=axiaoxin-com/logging&utm_campaign=Badge_Grade_Settings)
[![Build Status](https://travis-ci.org/axiaoxin-com/logging.svg?branch=master)](https://travis-ci.org/axiaoxin-com/logging)


logging 简单封装了在日常使用 [zap](https://github.com/uber-go/zap) 打日志时的常用方法。

- 提供快速使用 zap 打印日志的方法，除 zap 的 DPanic 、 DPanicf 方法外所有日志打印方法开箱即用
- 提供多种快速创建 logger 的方法
- 集成 **Sentry**，设置 DSN 后可直接使用 Sentry ，支持在使用 Error 及其以上级别打印日志时自动将该事件上报到 Sentry
- 支持从 Context 中创建、获取带有 **Trace ID** 的 logger
- 提供 gin 的日志中间件，支持使用 logging 的 logger 打印访问日志，支持 Trace ID，可以记录更加详细的请求和响应信息，支持通过配置自定义。
- 支持服务内部函数方式和外部 HTTP 方式**动态调整日志级别**，无需修改配置、重启服务
- 支持自定义 logger Encoder 配置
- 支持将日志保存到文件并自动 rotate
- 支持 Gorm 日志打印 Trace ID

logging 只提供 zap 使用时的常用方法汇总，不是对 zap 进行二次开发，拒绝过度封装。

## 安装

```
go get -u github.com/axiaoxin-com/logging
```

## 开箱即用

logging 提供的开箱即用方法都是使用自身默认 logger 克隆出的 CtxLogger 实际执行的。
在 logging 被 import 时，会生成内部使用的默认 logger 。
默认 logger 使用 JSON 格式打印日志内容到 stderr 。
默认不带 Sentry 上报功能，可以通过设置环境变量或者替换 logger 方法支持。
默认 logger 可通过代码内部动态修改日志级别， 默认不支持 HTTP 方式动态修改日志级别，需要指定端口创建新的 logger 来支持。
默认带有初始字段 pid 打印进程 ID 。

开箱即用的方法第一个参数为 context.Context, 可以传入 gin.Context ，会尝试从其中获取 Trace ID 进行日志打印，无需 Trace ID 可以直接传 nil

**示例 [example/logging.go](example/logging.go)**

全局开箱即用的方法默认不支持 sentry 自动上报 Error 级别的事件，有两种方式可以使其支持：

1. 通过设置系统环境变量 `SENTRY_DSN` 和 `SENTRY_DEBUG` 来实现自动上报。

2. 也可以通过替换默认 logger 来实现让全局方法支持 Error 以上级别自动上报。

**示例 [example/replace.go](example/replace.go)**

## 快速获取、创建你的 Logger

logging 提供多种方式快速获取一个 logger 来打印日志

**示例 [example/logger.go](example/logger.go)**

## 带 Trace ID 的 CtxLogger

每一次函数或者 gin 的 http 接口调用，在最顶层入口处都将一个带有唯一 trace id 的 logger 放入 context.Context 或 gin.Context ，
后续函数在内部打印日志时从 Context 中获取带有本次调用 trace id 的 logger 来打印日志几个进行调用链路跟踪。

**示例 1 普通函数中打印打印带 Trace ID 的日志 [example/context.go](example/context.go)**

**示例 2 gin 中打印带 Trace ID 的日志 [example/gin.go](example/gintraceid.go)**:

## 动态修改 logger 日志级别

logging 可以在代码中对 AtomicLevel 调用 SetLevel 动态修改日志级别，也可以通过请求 HTTP 接口修改。
创建 logger 时可自定义端口运行 HTTP 服务来接收请求修改日志级别。实际使用中日志级别通常写在配置文件中，
可以通过监听配置文件的修改来动态调用 SetLevel 方法。

**示例 [example/atomiclevel.go](example/atomiclevel.go)**

## 自定义 logger Encoder 配置

**示例 [example/encoder.go](example/encoder.go)**

## 日志保存到文件并自动 rotate

使用 lumberjack 将日志保存到文件并 rotate ，采用 zap 的 RegisterSink 方法和 Config.OutputPaths 字段添加自定义的日志输出的方式来使用 lumberjack 。

**示例 [example/lumberjack.go](example/lumberjack.go)**

## 支持 Gorm 日志打印 Trace ID

使用 gorm v2 支持 context logger 打印 trace id

**示例 [example/gorm.go](example/gorm.go)**

## gin middleware: GinLogger

GinLogger uses zap to log detailed access logs in JSON or text format with trace id, supports flexible and rich configuration,
and supports automatic reporting of log events above error level to sentry

相关文章： <https://github.com/axiaoxin/axiaoxin/issues/17>

示例： [example/ginlogger.go](./example/ginlogger.go)
