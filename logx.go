package logx

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (

	// LogDate 日期
	LogDate = 1 << iota

	// LogTime 时间
	LogTime

	// LogMicroSeconds ms
	LogMicroSeconds

	// LogLongFile 完整文件名
	LogLongFile

	// LogShortFile 较短的文件名
	LogShortFile

	// LogModule 包名
	LogModule

	// LogLevel 日志等级
	LogLevel
)

const (
	LevelTest = iota
	LevelDebug
	LevelInfo
	LevelNotice
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

const (
	endColor = "\033[0m"
)

var (
	// StdFlags 标准输出的格式
	StdFlags = LogDate | LogMicroSeconds | LogShortFile | LogLevel

	logColor = []string{
		LevelTest:   "\033[1;37m", //白色
		LevelInfo:   "\033[1;37m", //白色
		LevelDebug:  "\033[1;34m", //蓝色
		LevelNotice: "\033[1;33m", //黄色
		LevelWarn:   "\033[1;32m", //绿色
		LevelError:  "\033[1;31m", //红色
		LevelPanic:  "\033[1;35m", //紫红色
		LevelFatal:  "\033[1;36m", //青蓝色
	}

	smallLevels = []string{
		"T",
		"D",
		"I",
		"N",
		"W",
		"E",
		"P",
		"F",
	}

	levels = []string{
		"TEST",
		"DEBU",
		"INFO",
		"NOTI",
		"WARN",
		"ERRO",
		"PANI",
		"FATA",
	}
)

// LogFormat logformat enum
type LogFormat int

const (

	// LogFormatText text or default format
	LogFormatText LogFormat = iota

	//LogFormatJSON json format
	LogFormatJSON
)

// Logger logger struct
//	call logx.New() returns *Logger
type Logger struct {
	w         io.Writer
	ccPool    *sync.Pool
	bufPool   *sync.Pool
	level     int
	flag      int
	callDepth int
	prefix    string
	color     bool
	logFormat LogFormat
}

// New logx.New(...)
//	returns *Logger
func New(writer ...io.Writer) *Logger {
	var w io.Writer
	if len(writer) > 0 {
		w = writer[0]
	}
	return &Logger{
		flag:      StdFlags,
		level:     0,
		w:         w,
		callDepth: 2,
		ccPool: &sync.Pool{
			New: func() interface{} {
				return new(LogContent)
			},
		},
		bufPool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(nil)
			},
		},
	}
}

// SetWriter set io.writer
//	returns *Logger
func (log *Logger) SetWriter(w io.Writer) *Logger {
	log.w = w
	return log
}

// SetPrefix set *Logger's prefix
//	returns *Logger
func (log *Logger) SetPrefix(prefix string) *Logger {
	log.prefix = prefix
	return log
}

//	SetFlag set *Logger's flag
//	returns *Logger
func (log *Logger) SetFlag(flag int) *Logger {
	log.flag = flag
	return log
}

func (log *Logger) SetLevel(level int) *Logger {
	log.level = level
	return log
}

func (log *Logger) SetColor(color bool) *Logger {
	log.color = color
	return log
}

// SetFormat set log format
//	text or json or other
func (log *Logger) SetFormat(logFormat LogFormat) *Logger {
	log.logFormat = logFormat
	return log
}

func (log *Logger) SetCallDepth(depth int) *Logger {
	log.callDepth = depth
	return log
}

func (log *Logger) SetCallDepthPlus() *Logger {
	log.callDepth = log.callDepth + 1
	return log
}

func (log *Logger) GetCallDepth() int {
	return log.callDepth
}

func (log *Logger) Info(format string, v ...interface{}) {
	if LevelInfo < log.level {
		return
	}
	log.output(LevelInfo, fmt.Sprintf(format, v...))
}

func (log *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug < log.level {
		return
	}
	log.output(LevelDebug, fmt.Sprintf(format, v...))
}

func (log *Logger) Notice(format string, v ...interface{}) {
	if LevelNotice < log.level {
		return
	}
	log.output(LevelNotice, fmt.Sprintf(format, v...))
}

func (log *Logger) Error(format string, v ...interface{}) {
	if LevelError < log.level {
		return
	}
	log.output(LevelError, fmt.Sprintf(format, v...))
}

func (log *Logger) Warn(format string, v ...interface{}) {
	if LevelWarn < log.level {
		return
	}
	log.output(LevelWarn, fmt.Sprintf(format, v...))
}

func (log *Logger) Panic(format string, v ...interface{}) {
	if LevelPanic < log.level {
		return
	}
	s := fmt.Sprintf(format, v...)
	log.output(LevelPanic, s)
	panic(s)
}

func (log *Logger) Fatal(format string, v ...interface{}) {
	if LevelFatal < log.level {
		return
	}
	log.output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(-1)
}

func (log *Logger) output(level int, msg string) {
	var (
		now  = time.Now()
		name string
		line int
	)
	if log.flag&(LogShortFile|LogLongFile) != 0 {
		ok := false
		_, name, line, ok = runtime.Caller(log.callDepth)
		if !ok {
			name = "???"
			line = 0
		}
	}

	cc := log.ccPool.Get().(*LogContent)
	s := log.bufPool.Get().(*bytes.Buffer)
	cc.Color = log.color

	if log.prefix != "" {
		cc.Prefix = log.prefix
	}
	if log.flag&(LogDate|LogTime|LogMicroSeconds) != 0 {
		s.Reset()
		if log.flag&LogDate != 0 {
			year, month, day := now.Date()
			s.WriteString(strconv.FormatInt(int64(year), 10))
			s.WriteString("/")
			s.WriteString(strconv.FormatInt(int64(month), 10))
			s.WriteString("/")
			s.WriteString(strconv.FormatInt(int64(day), 10))
			s.WriteString(" ")
		}
		if log.flag&(LogTime|LogMicroSeconds) != 0 {
			hour, min, sec := now.Clock()
			s.WriteString(strconv.FormatInt(int64(hour), 10))
			s.WriteString(":")
			s.WriteString(strconv.FormatInt(int64(min), 10))
			s.WriteString(":")
			s.WriteString(strconv.FormatInt(int64(sec), 10))
			if log.flag&LogMicroSeconds != 0 {
				s.WriteString(".")
				s.WriteString(strconv.FormatInt(int64(now.Nanosecond()/1e6), 10))
			}
		}
		cc.Time = s.String()
	}
	if log.flag&LogLevel != 0 {
		cc.LevelInt = level
	}

	// filename and line
	if log.flag&(LogShortFile|LogLongFile) != 0 {
		s.Reset()
		if log.flag&LogShortFile != 0 {
			i := strings.LastIndex(name, "/")
			name = name[i+1:]
		}
		s.WriteString(name)
		s.WriteString(":")
		s.WriteString(strconv.FormatInt(int64(line), 10))
		cc.File = s.String()
	}
	cc.Msg = msg
	s.Reset()
	switch log.logFormat {
	case LogFormatText:
		s.Write(cc.Text())
	case LogFormatJSON:
		s.Write(cc.Json())
	default:
		s.Write(cc.Text())
	}
	s.WriteByte('\n')
	_, _ = log.w.Write(s.Bytes())

	log.ccPool.Put(cc)
	log.bufPool.Put(s)
}
