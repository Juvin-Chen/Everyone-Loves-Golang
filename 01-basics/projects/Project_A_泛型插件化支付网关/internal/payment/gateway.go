package payment

import (
	"fmt"
)

// Result 泛型支付结果：T可适配不同支付方式的返回数据（如支付宝交易号、微信订单号）
type Result[T any] struct {
	Success bool    // 支付是否成功
	PayType string  // 支付方式（alipay/wechat）
	Amount  float64 // 支付金额
	Data    T       // 泛型字段：存储不同支付方式的专属返回数据（如交易号） string / int / 结构体类型
	Message string  // 支付结果描述
}

// 支付网关 (Gateway)
type Gateway[T any] struct{}

// ProcessPayment 网关核心方法：泛型处理支付请求,T：泛型参数，对应不同支付方式的专属返回数据（如string类型的交易号）
func (g *Gateway[T]) ProcessPayment(
	method PaymentStrategy,
	amount float64,
	payType string,
) Result[T] {
	// 1. 网关前置校验：金额合法性
	if amount <= 0 {
		return Result[T]{
			Success: false,
			PayType: payType,
			Amount:  amount,
			Message: "网关校验失败：支付金额必须大于0",
		}
	}

	// 2. 调用具体支付方式的Pay方法（插件化逻辑）
	tradeNo, err := method.Pay(amount)
	if err != nil {
		// 支付失败：返回失败结果
		return Result[T]{
			Success: false,
			PayType: payType,
			Amount:  amount,
			Message: fmt.Sprintf("支付失败：%v", err),
		}
	}

	// 3. 支付成功：构造泛型结果（将交易号转为泛型T）,这里通过类型断言适配泛型T（假设T是string类型，对应交易号）
	var data T
	if v, ok := any(tradeNo).(T); ok {
		data = v
	}

	return Result[T]{
		Success: true,
		PayType: payType,
		Amount:  amount,
		Data:    data,
		Message: "支付成功",
	}
}
