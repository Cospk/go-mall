package middleware

import (
	"bytes"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// TraceMiddleware 添加traceId信息跟踪链路，方便后期做记录（注意：放第一个）
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.Request.Header.Get("Trace-Id")
		pSpanId := c.Request.Header.Get("Span-Id")
		spanId := utils.GenerateSpanId(c.Request.RemoteAddr)
		if traceId == "" {
			traceId = spanId
		}
		c.Set("Trace-Id", traceId)
		c.Set("Span-Id", spanId)
		c.Set("Parent-Span-Id", pSpanId)
		c.Next()
	}
}

// 自定义ResponseWriter
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写Write方法
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware 记录请求信息和响应信息的日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求信息日志（前置）
		var requestBody []byte
		// 若是文件上传请求，不在日志里记录body
		if !strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}
		start := time.Now()
		writeAccessLog(c, "access_start", time.Since(start), requestBody, nil)

		// 初始化响应记录器
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next() //跳出中间件，执行其他中间件以及路由处理函数，最后执行下面的代码

		// 记录响应信息
		var responseLogging string
		if c.Writer.Size() > 10*1024 { // 响应大于10KB 不记录
			responseLogging = "Response data size is too Large to log"
		} else {
			responseLogging = blw.body.String()
		}
		writeAccessLog(c, "access_end", time.Since(start), requestBody, responseLogging)
	}
}

func writeAccessLog(c *gin.Context, accessType string, cost time.Duration, body []byte, response interface{}) {
	req := c.Request
	logger.NewLogger(c).Info("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		//"token", req.Header.Get("token"),
		"method", req.Method,
		"path", req.URL.Path,
		"query", req.URL.RawQuery,
		"body", string(body),
		"response", response,
		"time(ms)", int64(cost/time.Millisecond),
	)
}

// RecoveryMiddleware gin虽然有默认的异常处理，但是无法记录完整的堆栈信息，这里自定义异常处理
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是因为连接中断导致的错误（非应用逻辑错误）
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					// 若是连接中断导致的错误，只要记录路径、错误信息以及请求信息即可，不需要记录堆栈信息直接终止处理
					logger.NewLogger(c).Error("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					c.Error(err.(error))
					c.Abort()
					return
				}

				// 普通panic记录完整堆栈
				logger.NewLogger(c).Error("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))
				// 返回 500 错误
				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}

		}()
		c.Next()
	}
}
