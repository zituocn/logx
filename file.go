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
	StorageTypeMinutes StorageType = iota
	StorageTypeHour
	StorageTypeDay
	StorageTypeMonth
)

var formats = map[StorageType]string{
	StorageTypeMinutes: "2006-01-02-15-04",
	StorageTypeHour:    "2006-01-02-15",
	StorageTypeDay:     "2006-01-02",
	StorageTypeMonth:   "2006-01",
}

var defaultMaxDay = 7

func (s StorageType) getFileFormat() string {
	return formats[s]
}

type FileOptions struct {
	StorageType StorageType
	MaxDay      int
	Dir         string
	Prefix      string
}

type FileWriter struct {
	file   *os.File
	mu     sync.Mutex
	closed bool
	stopCh chan struct{} // 通知后台 goroutine 退出
	FileOptions
	date string
}

func NewFileWriter(opts ...FileOptions) *FileWriter {
	opt := prepareFileWriterOption(opts)
	w := &FileWriter{
		FileOptions: opt,
		stopCh:      make(chan struct{}),
	}
	go w.run() // 单一后台 goroutine 处理清理 + 轮转
	return w
}

// run 单一 goroutine，每小时检查一次
func (w *FileWriter) run() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.mu.Lock()
			w.clearExpiredFiles()
			w.tryRotate()
			w.mu.Unlock()
		case <-w.stopCh:
			return
		}
	}
}

// Close 安全关闭 FileWriter，flush 并停止后台 goroutine
func (w *FileWriter) Close() error {
	close(w.stopCh)
	w.mu.Lock()
	defer w.mu.Unlock()
	w.closed = true
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *FileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, os.ErrClosed
	}
	w.initFileLocked()
	return w.file.Write(p)
}

func (w *FileWriter) initFileLocked() {
	date := time.Now().Format(w.StorageType.getFileFormat())
	if w.date == date && w.file != nil {
		return // 文件正确，无需操作
	}
	// 日期变化，关闭旧文件
	if w.file != nil {
		_ = w.file.Close()
		w.file = nil
	}
	dir := filepath.Dir(w.Dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("%s.%s.log", w.Prefix, date)
	// 使用 filepath.Join 拼接路径
	fullPath := filepath.Join(w.Dir, fileName)
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic(err)
	}
	w.file = file
	w.date = date
}

// tryRotate 按日期轮转（由定时器触发）
func (w *FileWriter) tryRotate() {
	date := time.Now().Format(w.StorageType.getFileFormat())
	if w.date != date && w.file != nil {
		_ = w.file.Close()
		w.file = nil
		w.date = ""
	}
}

// clearExpiredFiles 清理过期文件（调用方需持锁）
func (w *FileWriter) clearExpiredFiles() {
	files := getDirFiles(w.Dir)
	now := time.Now()
	maxDur := time.Hour * 24 * time.Duration(w.MaxDay-1)
	for _, item := range files {
		if item.ModTime.Add(maxDur).Before(now) {
			// 使用 filepath.Join
			_ = os.Remove(filepath.Join(w.Dir, item.Name))
		}
	}
}

// FileInfo 文件信息
type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}

// getDirFiles 返回目录下所有非目录文件（使用 os.ReadDir 替代 ioutil.ReadDir）
func getDirFiles(path string) (files []*FileInfo) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, &FileInfo{
			Name:    entry.Name(),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		})
	}
	return
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
	// 确保 Dir 以路径分隔符结尾
	if len(opt.Dir) > 0 && opt.Dir[len(opt.Dir)-1] != '/' {
		opt.Dir = opt.Dir + "/"
	}
	return opt
}
