package logpipe

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
	"time"
)

func TestPipe_Write(t *testing.T) {
	// 创建测试用的logger
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "写入普通数据",
			data:    []byte("test message\n"),
			wantErr: false,
		},
		{
			name:    "写入空数据",
			data:    []byte{},
			wantErr: false,
		},
		{
			name:    "写入多行数据",
			data:    []byte("line1\nline2\n"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(logger, "test-prefix")
			n, err := p.Write(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.data), n)
			}

			// 清理
			p.Close()
		})
	}
}

func TestPipe_Close(t *testing.T) {
	logger := zaptest.NewLogger(t)
	p := New(logger, "test-prefix")

	// 测试关闭管道
	err := p.Close()
	assert.NoError(t, err)

	// 验证写入已关闭的管道会返回错误
	_, err = p.Write([]byte("test"))
	assert.Error(t, err)
}

func TestNewWithOption(t *testing.T) {
	tests := []struct {
		name   string
		option Option
	}{
		{
			name: "使用自定义logger",
			option: Option{
				Logger: zap.NewExample(),
				Prefix: "custom-prefix",
			},
		},
		{
			name: "使用nil logger",
			option: Option{
				Logger: nil,
				Prefix: "nil-logger-prefix",
			},
		},
		{
			name:   "使用空Option",
			option: Option{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewWithOption(tt.option)
			assert.NotNil(t, p)
			assert.NotNil(t, p.w)

			// 写入数据并验证时间戳更新
			initialTime := p.LastLineTimestamp.Load().(time.Time)
			_, err := p.Write([]byte("test message\n"))
			assert.NoError(t, err)

			// 等待一小段时间确保后台处理完成
			time.Sleep(100 * time.Millisecond)

			updatedTime := p.LastLineTimestamp.Load().(time.Time)
			assert.True(t, updatedTime.After(initialTime))

			p.Close()
		})
	}
}

func TestNew(t *testing.T) {
	logger := zaptest.NewLogger(t)
	prefix := "test-prefix"

	p := New(logger, prefix)
	assert.NotNil(t, p)
	assert.NotNil(t, p.w)

	// 测试基本功能
	_, err := p.Write([]byte("test message\n"))
	assert.NoError(t, err)

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证时间戳已更新
	timestamp := p.LastLineTimestamp.Load().(time.Time)
	assert.False(t, timestamp.IsZero())

	p.Close()
}

func TestPipe_ConcurrentWrites(t *testing.T) {
	logger := zaptest.NewLogger(t)
	p := New(logger, "test-prefix")

	// 并发写入测试
	const numGoroutines = 10
	const messagesPerGoroutine = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < messagesPerGoroutine; j++ {
				_, err := p.Write([]byte("test message\n"))
				assert.NoError(t, err)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// 等待处理完成
	time.Sleep(100 * time.Millisecond)

	// 验证最终时间戳
	timestamp := p.LastLineTimestamp.Load().(time.Time)
	assert.False(t, timestamp.IsZero())

	p.Close()
}
