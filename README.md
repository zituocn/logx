# logx

一个简单又好用的go日志库

## 特性

1. API简单，可快速上手使用
2. 可设置日志输出格式
3. 可定义日志输出等级
4. 可定义是否颜色输出
5. 可实现自己的 `io.writer`

## 安装

```sh
go get github.com/zituocn/logx
```


## 1、在控制台打印日志

### 1.1 快速上手

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

输出

![](https://p1.22v.net/topic/20221123/0bff62601c7a533f.png)

### 1.2 支持`printf`式的格式化

```go
package main

import (
	"time"

	"github.com/zituocn/logx"
)

func main() {
	logx.Infof("这是一个字串: %s", time.Now())
	logx.Debugf("这是一个对象: %#v", time.Now())
}
```

输出

```sh
2022/11/23 21:16:49.507 [INFO] main.go:10: 这是一个字串: 2022-11-23 21:16:49.507432 +0800 CST m=+0.000194033
2022/11/23 21:16:49.507 [DEBU] main.go:11: 这是一个对象: time.Date(2022, time.November, 23, 21, 16, 49, 507724000, time.Local)
```

### 1.3 自定义配置

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
	logx.Info("info str")
	logx.Debug("debug str")
}
```

输出

```sh
{"prefix":"demo","time":"2022/11/23 20:25:6.166","level":"","file":"main.go:17","msg":"info str"}
{"prefix":"demo","time":"2022/11/23 20:25:6.166","level":"","file":"main.go:18","msg":"debug str"}
```

#### 参数说明

* logx.SetWriter -> 设置输出到哪一个io.writer
* SetColor -> 设置是否输出颜色
* SetFormat(logx.LogFormatJSON) -> 日志输出为json格式
* SetPrefix -> 设置日志前缀

---
## 2. 同时写日志到文件
