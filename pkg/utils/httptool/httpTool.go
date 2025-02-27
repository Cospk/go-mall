package httptool

import (
	"bytes"
	"context"
	"github.com/Cospk/go-mall/pkg/errcode"
	"github.com/Cospk/go-mall/pkg/logger"
	"github.com/Cospk/go-mall/pkg/utils"
	"io/ioutil"
	"net/http"
	"time"
)

// Request 发起HTTP请求
func Request(method string, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 1、初始化配置
	start := time.Now()
	reqOpts := defaultRequestOptions() // 获取默认请求选项
	for _, opt := range options {      // 应用传入的配置选项
		if err = opt.apply(reqOpts); err != nil {
			return
		}
	}
	log := logger.NewLogger(reqOpts.ctx) //创建日志记录器
	defer func() {
		// 结束部分
		if err != nil {
			log.Error("HTTP请求失败日志", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err)
		}
	}()

	// 创建请求对象
	req, err := http.NewRequestWithContext(reqOpts.ctx, method, url, bytes.NewReader(reqOpts.data))
	if err != nil {
		return
	}
	defer req.Body.Close()

	// 在header添加追踪信息，把内部服务的日志关联起来
	setTraceHeaders(req, reqOpts.ctx)
	for k, v := range reqOpts.headers {
		req.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{Timeout: reqOpts.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	dur := time.Since(start).Milliseconds()
	if dur >= 3000 { // 超过3秒返回一条warn记录
		log.Warn("HTTP请求慢日志", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err, "cost/ms", dur)
	} else {
		log.Info("HTTP请求debug日志", "method", method, "url", url, "body", reqOpts.data, "reply", respBody, "err", err, "cost/ms", dur)
	}
	httpStatusCode = resp.StatusCode
	if httpStatusCode != http.StatusOK {
		// 状态码非200，当做error处理
		err = errcode.Wrap("HTTP请求失败", err)
		return
	}

	// 读取响应体
	respBody, _ = ioutil.ReadAll(resp.Body)

	return
}

func setTraceHeaders(req *http.Request, ctx context.Context) {
	traceId, spanId, _ := utils.GetTraceInfoFromCtx(ctx)
	req.Header.Set("Trace-Id", traceId)
	req.Header.Set("Span-Id", spanId)
}

// Get 发起GET请求
func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request("GET", url, options...)
}

// Post 发起POST请求
func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, respBody []byte, err error) {
	// 默认自带Header Content-Type: application/json 可通过 传递 WithHeaders 增加或者覆盖Header信息
	defaultHeader := map[string]string{"Content-Type": "application/json"}
	var newOptions []Option
	newOptions = append(newOptions, WithHeaders(defaultHeader), WithData(data), WithContext(ctx))
	newOptions = append(newOptions, options...)

	httpStatusCode, respBody, err = Request("POST", url, newOptions...)
	return
}

// 针对可选的HTTP请求配置项，模仿gRPC使用的Options设计模式实现
type requestOption struct {
	ctx     context.Context
	timeout time.Duration
	data    []byte
	headers map[string]string
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(option *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

// 默认请求选项
func defaultRequestOptions() *requestOption {
	return &requestOption{ // 默认请求选项
		ctx:     context.Background(),
		timeout: 5 * time.Second,
		data:    nil,
		headers: map[string]string{},
	}
}

// WithContext 设置请求上下文
func WithContext(ctx context.Context) Option {
	return optionFunc(func(opts *requestOption) error {
		opts.ctx = ctx
		return nil
	})
}

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

// WithHeaders 设置请求头
func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return
	})
}

// WithData 设置请求体数据
func WithData(data []byte) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}
