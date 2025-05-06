package payment

import (
	"errors"
	"sync"
)

// 支付方式类型
const (
	PaymentTypeWechat = "wechat"
	PaymentTypeAlipay = "alipay"
)

var (
	paymentMap = make(map[string]Payment)
	mu         sync.RWMutex

	ErrPaymentNotFound = errors.New("payment method not found")
)

// Register 注册支付方式
func Register(name string, payment Payment) {
	mu.Lock()
	defer mu.Unlock()
	paymentMap[name] = payment
}

// Get 获取支付方式
func Get(name string) (Payment, error) {
	mu.RLock()
	defer mu.RUnlock()

	payment, ok := paymentMap[name]
	if !ok {
		return nil, ErrPaymentNotFound
	}

	return payment, nil
}

// GetAll 获取所有支付方式
func GetAll() map[string]Payment {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string]Payment, len(paymentMap))
	for k, v := range paymentMap {
		result[k] = v
	}

	return result
}

// Unregister 注销支付方式
func Unregister(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(paymentMap, name)
}
