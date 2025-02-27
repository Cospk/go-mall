package utils

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 生成唯一的SpanId：
// 1、ip地址+时间戳+随机数
// 2、UUID完全唯一，但长度128位，性能较差，不适合追踪系统
// 3、雪花算法，分布式系统中生成唯一ID的一种算法，性能较高，但需要分配节点id

// 新增全局变量（带互斥锁的随机源）
var (
	randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
	randLock   sync.Mutex
)

// GenerateSpanId 生成spanId，思路：ip（空间唯一）与时间戳（时间唯一）异或运算+随机数
func GenerateSpanId(addr string) string {
	strAddr := strings.Split(addr, ":")
	ip := strAddr[0]
	// 获取IP数值
	ipVal, _ := ipv4ToUint32(ip)

	//获取时间戳
	times := uint64(time.Now().UnixNano())

	// 获取32位随机数
	randLock.Lock()
	random := uint64(randSource.Int31())
	randLock.Unlock()

	// 组合成spanId (将ip和时间戳进行异或运算，既保证唯一性又保证ip信息安全。然后左移保留低32位后与随机数进行或运算)
	spanId := ((times ^ uint64(ipVal)) << 32) | random

	return strconv.FormatUint(spanId, 16)
}
func ipToUint64(ipStr string) (uint64, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("非法ip地址")
	}
	// 统一转换为 16 字节
	ip = ip.To16()
	if ip == nil {
		return 0, fmt.Errorf("ip地址转换失败")
	}
	// 取后 8 字节进行异或运算
	lower64 := binary.BigEndian.Uint64(ip[8:])
	upper64 := binary.BigEndian.Uint64(ip[:8])
	return lower64 ^ upper64, nil
}

func ipv4ToUint32(ip string) (uint32, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(ipAddr.IP.To4()), nil
}

// GetTraceInfoFromCtx 从ctx中获取trace信息
func GetTraceInfoFromCtx(ctx context.Context) (traceId, spanId, pSpanId string) {
	if ctx.Value("Trace-Id") != nil {
		traceId = ctx.Value("Trace-Id").(string)
	}
	if ctx.Value("Span-Id") != nil {
		spanId = ctx.Value("Span-Id").(string)
	}
	if ctx.Value("Parent-Span-Id") != nil {
		pSpanId = ctx.Value("Parent-Span-Id").(string)
	}
	return
}
