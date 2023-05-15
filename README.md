# boltlog
![](https://github.com/harley9293/blotlog/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/harley9293/blotlog/branch/master/graph/badge.svg?token=gLxJw9swlO)](https://codecov.io/gh/harley9293/blotlog)

A simple log library written in go language, which outputs logs in the following format on the console by default

```text
2022-07-01 20:44:15.055|DEBUG|[log_test.go:15:hello]|gid:19|message
```

## Installation

Install blotlog using the go get command:

```shell
go get -u github.com/harley9293/blotlog
```

## Usage

First, import the blotlog library:

```go
import "github.com/harley9293/blotlog"
```

You can use this library directly without any configuration, by default all logs will be printed to the console

```go
log.Debug/Info/Warn/Error("hello world %d", 2023)
```

You can also modify the default configuration to use the logging library according to your needs

```go
// set log level
log.SetLevel(log.ErrorLevel)

// discard console output
log.ConsoleOff()

// set rotate config
log.AddRotateHook(&log.RotateConf{})
```