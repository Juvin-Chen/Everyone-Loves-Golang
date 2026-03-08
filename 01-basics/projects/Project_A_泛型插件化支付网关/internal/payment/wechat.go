package payment

import "fmt"

// 具体支付插件2 WeChat
type WeChat struct {
	BasePayment
}

func (w *WeChat) Pay(amount float64) (string, error) {
	// 校验扣减金额
	if amount <= 0 {
		return "", fmt.Errorf("支付金额必须为正数，当前金额：%.2f", amount)
	}

	// 扣减余额
	success := w.UpdateBalance(amount)
	if !success {
		return "", fmt.Errorf("余额扣减失败，无法完成支付")
	}

	// 执行微信支付
	w.Log(fmt.Sprintf("执行微信支付逻辑，扣减金额：%.2f元", amount))

	// 返回支付结果
	tradeNo := fmt.Sprintf("WX_%s_%f", w.PaymentID, amount*100) // 模拟交易号
	w.Log(fmt.Sprintf("微信支付成功，交易号：%s", tradeNo))
	return tradeNo, nil
}
