package log

import (
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	lg      *slog.Logger
	lgLevel *slog.LevelVar
	opts    *slog.HandlerOptions
	once    sync.Once
)

// Debug 记录调试级别日志
func Debug(msg string, args ...any) {
	getLogger().Debug(msg, args...)
}

// Info 记录信息级别日志
func Info(msg string, args ...any) {
	getLogger().Info(msg, args...)
}

// Warn 记录警告级别日志
func Warn(msg string, args ...any) {
	getLogger().Warn(msg, args...)
}

// Error 记录错误级别日志
func Error(msg string, args ...any) {
	getLogger().Error(msg, args...)
}

// SetLevel 设置日志级别
func SetLevel(level slog.Level) {
	lgLevel.Set(level)
}

func NewOptions() *slog.HandlerOptions {
	// 设置 HandlerOptions，自定义时间属性
	lgLevel = &slog.LevelVar{}
	lgLevel.Set(slog.LevelDebug)
	return &slog.HandlerOptions{
		Level: lgLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 如果当前属性是时间戳
			if a.Key == slog.TimeKey && len(groups) == 0 {
				a.Key = "time" // 键名可以保持不变或修改
				// 将时间值转换为自定义格式
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}
}

// getLogger 获取单例日志实例
func getLogger() *slog.Logger {
	once.Do(func() {
		opts = NewOptions()
		lg = slog.New(slog.NewTextHandler(newLogWriter(), opts))
	})
	return lg
}

func newLogWriter() io.Writer {
	// TODO close the file
	f, err := os.OpenFile("lockunlock.run.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	return NewMultiWriter(os.Stdout, f)
}

type MultiWriter struct {
	writerList []io.Writer
}

func NewMultiWriter(wlist ...io.Writer) *MultiWriter {
	return &MultiWriter{writerList: wlist}
}

func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writerList {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return
}

func init() {
	getLogger()
}
