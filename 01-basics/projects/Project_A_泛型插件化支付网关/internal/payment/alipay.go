package payment

import "fmt"

// 具体支付插件1 Alipay
type AliPay struct {
	BasePayment
}

func (a *AliPay) Pay(amount float64) (string, error) {
	if amount <= 0 {
		return "", fmt.Errorf("支付金额必须为正数，当前金额：%.2f", amount)
	}

	// 扣减余额
	success := a.UpdateBalance(amount)
	if !success {
		return "", fmt.Errorf("余额扣减失败，无法完成支付")
	}

	// 执行支付宝支付
	a.Log(fmt.Sprintf("执行支付宝支付，扣减金额：%.2f元", amount))

	// 返回支付结果
	tradeNo := fmt.Sprintf("AliPay_%s_%f", a.PaymentID, amount*100) // 模拟交易号
	a.Log(fmt.Sprintf("支付宝支付成功，交易号：%s", tradeNo))
	return tradeNo, nil
}
