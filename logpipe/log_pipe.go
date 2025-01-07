package logpipe

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"sync/atomic"
	"time"
)

// Pipe 用于将数据写入管道并记录最后一行的时间戳
type Pipe struct {
	w                 *io.PipeWriter // 管道的写入端
	LastLineTimestamp atomic.Value   // 最后一行的时间戳，使用 atomic.Value 以支持任意类型存储
}

// Write 将数据写入管道，并返回写入的字节数和可能的错误
func (p *Pipe) Write(data []byte) (int, error) {
	w := p.w
	n, err := w.Write(data)
	return n, err
}

// Close 关闭管道的写入端
func (p *Pipe) Close() error {
	return p.w.Close()
}

// Option 定义用于创建 Pipe 的选项
type Option struct {
	Logger *zap.Logger // 日志记录器，如果为 nil 则使用默认无操作记录器
	Prefix string      // 日志前缀
}

// NewWithOption 使用指定的选项创建一个新的 Pipe 实例
func NewWithOption(option Option) *Pipe {
	// 如果未提供日志记录器，则使用无操作记录器
	logger := option.Logger
	if logger == nil {
		logger = zap.NewNop()
	}
	// 禁用调用者信息
	logger = logger.WithOptions(zap.WithCaller(false))

	// 创建管道的读写端
	r, w := io.Pipe()
	// 使用 bufio.Scanner 按行扫描管道中的数据
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	// 初始化 Pipe 结构体
	pipe := &Pipe{
		w: w,
	}
	// 初始化 LastLineTimestamp
	var initialTime time.Time
	pipe.LastLineTimestamp.Store(initialTime)

	// 启动后台协程处理管道中的数据
	go func() {
		defer r.Close() // 确保在协程结束时关闭读端

		for scanner.Scan() {
			// 获取当前时间
			now := time.Now()
			// 更新最后一行的时间戳
			pipe.LastLineTimestamp.Store(now)
			// 使用结构化日志记录日志信息
			logger.Debug("Log line",
				zap.String("prefix", option.Prefix),
				zap.String("message", scanner.Text()),
				zap.Time("timestamp", now),
			)
		}

		// 检查扫描过程中是否发生错误
		if err := scanner.Err(); err != nil {
			logger.Error("Error scanning pipe", zap.Error(err))
		}
	}()

	return pipe
}

// New 使用默认选项创建一个新的 Pipe 实例
func New(logger *zap.Logger, prefix string) *Pipe {
	return NewWithOption(Option{Prefix: prefix, Logger: logger})
}
