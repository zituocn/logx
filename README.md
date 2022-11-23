# logx 

一个好用又简单的go日志库

## 特性

1. API简单快速，可快速使用
2. 可设置日志输出格式
3. 可定义日志输出等级
4. 可定义是否颜色输出
5. 可实现自己的 `io.writer`

## 安装

```sh
go get github.com/zituocn/logx
```

## 使用


### 1、在控制台打印日志

快速使用，使用默认的配置

```go
package main

import (
	"github.com/zituocn/logx"
)

func main() {
	logx.Info("info")
	logx.Debug("debug")
	logx.Error("error")
	logx.Notice("notice")
	logx.Warn("warn")
	logx.Panic("panic")
	logx.Fatal("fatal")
}
```

### 自定义配置

可以配置输出格式、样式等

```go
package main

import (
	"os"
	"github.com/zituocn/logx"
)

func init() {
	logx.SetWriter(os.Stdout).
        SetColor(true).
        SetFormat(logx.LogFormatJSON).
        SetPrefix("demo")
}

func main() {
	logx.Info("info")
	logx.Debug("debug")
	logx.Error("error")
	logx.Notice("notice")
	logx.Warn("warn")
	logx.Panic("panic")
	logx.Fatal("fatal")
}

```
#### 参数说明

* logx.SetWriter -> 设置输出到哪一个io.writer
* SetColor -> 设置是否输出颜色
* SetFormat(logx.LogFormatJSON) -> 日志输出为json格式
* SetPrefix -> 设置日志前缀

