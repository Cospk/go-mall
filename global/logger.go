package global

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"runtime"
)

type logger struct {
	ctx     context.Context
	traceId string
	spanId  string
	pSpanId string
	_logger *zap.Logger
}

// traceId : 表示追踪整条请求的唯一标识
// spanId  : 表示当前请求的唯一标识
// pSpanId : 表示当前请求的父请求的唯一标识

// NewLogger 创建一个logger
func NewLogger(ctx context.Context) *logger {
	var traceId, spanId, pSanId string
	if ctx.Value("Trace-Id") != nil {
		traceId = ctx.Value("Trace-Id").(string)
	}
	if ctx.Value("Span-Id") != nil {
		spanId = ctx.Value("Span-Id").(string)
	}
	if ctx.Value("Parent-Span-Id") != nil {
		pSanId = ctx.Value("Parent-Span-Id").(string)
	}
	return &logger{
		ctx:     ctx,
		traceId: traceId,
		spanId:  spanId,
		pSpanId: pSanId,
		_logger: Logger,
	}
}

// Debug 打印debug日志
func (l *logger) Debug(msg string, kv ...interface{}) {
	l.log(zapcore.DebugLevel, msg, kv...)
}

// Info 打印info日志
func (l *logger) Info(msg string, kv ...interface{}) {
	l.log(zapcore.InfoLevel, msg, kv...)
}

// Warn 打印warn日志
func (l *logger) Warn(msg string, kv ...interface{}) {
	l.log(zapcore.WarnLevel, msg, kv...)
}

// Error 打印error日志
func (l *logger) Error(msg string, kv ...interface{}) {
	l.log(zapcore.ErrorLevel, msg, kv...)
}

// kv 应该是成对的数据, 类似: name,张三,age,10,...
func (l *logger) log(lvl zapcore.Level, msg string, kv ...interface{}) {
	// 保证要打印的日志信息成对出现
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 日志行信息中增加追踪参数
	kv = append(kv, "traceid", l.traceId, "spanid", l.spanId, "pspanid", l.pSpanId)
	// 增加日志调用者信息, 方便查日志时定位程序位置
	funcName, file, line := l.getLoggerCallerInfo()
	kv = append(kv, "func", funcName, "file", file, "line", line)

	fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		fields = append(fields, zap.Any(k, kv[i+1]))
	}
	ce := l._logger.Check(lvl, msg)
	ce.Write(fields...)
}

// getLoggerCallerInfo 日志调用者信息 -- 方法名, 文件名, 行号
func (l *logger) getLoggerCallerInfo() (funcName, file string, line int) {
	// 回溯拿调用日志方法的业务函数的信息
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	file = path.Base(file)
	funcName = runtime.FuncForPC(pc).Name()
	return
}
