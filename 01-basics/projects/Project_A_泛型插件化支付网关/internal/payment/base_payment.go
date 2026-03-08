package payment

import (
	"encoding/json"
	"fmt"
	"os"
)

// 定义一个接口，规定所有支付插件必须具备的能力
type PaymentStrategy interface {
	Pay(amount float64) (string, error)
}

// 基础支付结构体
type BasePayment struct {
	Paytype   string
	PaymentID string  // 支付AppID号，如："WeChat1001 / AliPay1001"
	Balance   float64 // 账户余额
}

// 基类所有方法的实现
// 1. Log()
// 2. GetBalance()
// 3. UpdateBalance()

// 所有的支付方式都需要记录日志
func (b *BasePayment) Log(msg string) {
	fmt.Printf("[%s%s] 支付日志：%s\n", b.Paytype, b.PaymentID, msg)
}

// 获取账户余额
func (b *BasePayment) GetBalance() (float64, error) {
	filePath := "configs/payment_balance.json"

	// 1. 打开JSON文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("打开余额文件失败：%w", err) // %w包装原始错误，便于排查
	}
	defer file.Close()

	// 2. 定义与JSON结构匹配的结构体
	type PaymentBalance struct {
		PaymentID string  `json:"PaymentID"`
		Balance   float64 `json:"Balance"`
	}

	// 3. 解析JSON文件到切片
	var balanceList []PaymentBalance
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&balanceList); err != nil {
		return 0, fmt.Errorf("解析余额文件失败：%w", err)
	}

	// 4. 遍历切片，查找对应PaymentID的余额
	for _, item := range balanceList {
		if item.PaymentID == b.PaymentID {
			return item.Balance, nil // 找到，返回余额
		}
	}

	// 5. 未找到对应ID的情况
	return 0, fmt.Errorf("未找到支付ID【%s】的余额记录", b.PaymentID)
}

// 更新账户余额
func (b *BasePayment) UpdateBalance(amount float64) bool {
	filePath := "configs/payment_balance.json"

	// 1. 校验金额：必须为正数（代表要扣除的金额）
	if amount <= 0 {
		b.Log(fmt.Sprintf("扣减金额必须为正数，传入金额%.2f不符合要求", amount))
		return false
	}

	// 2. 读取JSON文件中的所有支付余额数据
	type PaymentBalance struct {
		PaymentID string  `json:"PaymentID"`
		Balance   float64 `json:"Balance"`
	}

	var balanceList []PaymentBalance

	// 打开文件（文件不存在直接失败，因为无余额可扣）
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			b.Log(fmt.Sprintf("余额文件不存在，无法扣减%.2f元", amount))
			return false
		} else {
			b.Log(fmt.Sprintf("打开余额文件失败：%v", err))
			return false
		}
	}
	defer file.Close()

	// 解析JSON数据到切片
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&balanceList); err != nil {
		b.Log(fmt.Sprintf("解析余额文件失败：%v", err))
		return false
	}

	// 3. 查找对应PaymentID的余额并执行扣减
	found := false
	for i, item := range balanceList {
		if item.PaymentID == b.PaymentID {
			// 检查余额是否足够扣减
			if item.Balance < amount {
				b.Log(fmt.Sprintf("余额不足！当前余额%.2f元，尝试扣减%.2f元", item.Balance, amount))
				return false
			}
			// 执行扣减操作
			balanceList[i].Balance -= amount
			found = true
			b.Log(fmt.Sprintf("余额扣减成功！扣减%.2f元，剩余余额：%.2f元", amount, balanceList[i].Balance))
			break
		}
	}

	// 未找到对应PaymentID的记录，扣减失败
	if !found {
		b.Log(fmt.Sprintf("未找到支付ID【%s】的记录，无法扣减%.2f元", b.PaymentID, amount))
		return false
	}

	// 4. 确保configs目录存在
	if err := os.MkdirAll("configs", 0755); err != nil {
		b.Log(fmt.Sprintf("创建configs目录失败：%v", err))
		return false
	}

	// 5. 将扣减后的数据写回JSON文件
	updatedData, err := json.MarshalIndent(balanceList, "", "    ")
	if err != nil {
		b.Log(fmt.Sprintf("序列化余额数据失败：%v", err))
		return false
	}

	// 写入文件（覆盖原有内容）
	if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
		b.Log(fmt.Sprintf("写入余额文件失败：%v", err))
		return false
	}

	// 同步更新结构体自身的Balance字段（保持内存数据和文件一致）
	b.Balance, _ = b.GetBalance()
	return true
}
