package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type StorageType int

const (

	// StorageTypeMinutes 按分钟存储
	StorageTypeMinutes StorageType = iota

	// StorageTypeHour 按小时存储
	StorageTypeHour

	// StorageTypeDay 按天存储
	StorageTypeDay

	// StorageTypeMonth 按月存储
	StorageTypeMonth
)

var (
	formats = map[StorageType]string{
		StorageTypeMinutes: "2006-01-02-15-04",
		StorageTypeHour:    "2006-01-02-15",
		StorageTypeDay:     "2006-01-02",
		StorageTypeMonth:   "2006-01",
	}

	// defaultMaxDay 日志文件的默认最大保存天数
	// 7天之外的文件，会被自动清理
	defaultMaxDay = 7
)

func (s StorageType) getFileFormat() string {
	return formats[s]
}

// FileOptions 文件存储选项
type FileOptions struct {

	// StorageType 存储的时间类型
	StorageType StorageType

	// MaxDay 日志最大保存天数
	MaxDay int

	// Dir 日志保存目录
	Dir string

	// Prefix 文件名前缀
	Prefix string

	// date 日期
	data string
}

// FileWriter 文件存储实现
type FileWriter struct {
	file *os.File
	mu   *sync.Mutex

	FileOptions
}

func NewFileWriter(opts ...FileOptions) *FileWriter {
	opt := prepareFileWriterOption(opts)
	w := &FileWriter{
		FileOptions: opt,
		mu:          &sync.Mutex{},
	}
	return w
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.initFile()
	return w.file.Write(p)
}

func (w *FileWriter) initFile() {
	now := time.Now()
	date := now.Format(w.StorageType.getFileFormat())
	if w.data != date && w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	if w.file == nil {
		dir := filepath.Dir(w.Dir)
		err := os.MkdirAll(dir, 755)
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%s.%s.log", w.Prefix, date)
		file, errO := os.OpenFile(filepath.Join(w.Dir, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
		if errO != nil {
			panic(errO)
		}
		w.file = file
		w.data = date
	}
}

func prepareFileWriterOption(opts []FileOptions) FileOptions {
	var opt FileOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Dir == "" {
		opt.Dir = "./"
	}
	if opt.MaxDay <= 0 {
		opt.MaxDay = defaultMaxDay
	}
	if opt.Dir[len(opt.Dir)-1:] != "/" {
		opt.Dir = opt.Dir + "/"
	}
	return opt
}
