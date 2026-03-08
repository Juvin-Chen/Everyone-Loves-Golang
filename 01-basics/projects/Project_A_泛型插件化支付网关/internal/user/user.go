package user

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	Name     string `json:"Name"`
	AliPayID string `json:"AliPayID"`
	WeChatID string `json:"WeChatID"`
}

// 登录逻辑：仅校验姓名，匹配后加载对应的支付宝/微信ID
func (u *User) Login() bool {
	fmt.Println("请输入姓名：")
	_, err := fmt.Scan(&u.Name)
	if err != nil {
		fmt.Println("输入错误：", err)
		return false
	}

	// 1. 打开用户JSON配置文件
	file, err := os.Open("configs/users.json")
	if err != nil {
		fmt.Printf("读取用户数据失败：%v\n", err)
		return false
	}
	defer file.Close()

	// 2. 解析JSON到用户切片（仅包含3个字段）
	var userList []User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&userList)
	if err != nil {
		fmt.Printf("解析用户数据失败：%v\n", err)
		return false
	}

	// 3. 匹配姓名，加载支付ID
	loginSuccess := false
	for _, user := range userList {
		if user.Name == u.Name {
			u.AliPayID = user.AliPayID
			u.WeChatID = user.WeChatID
			loginSuccess = true
			break
		}
	}

	// 4. 输出结果
	if loginSuccess {
		fmt.Printf("登录成功！\n姓名：%s\n支付宝ID：%s\n微信ID：%s\n",
			u.Name, u.AliPayID, u.WeChatID)
		return true
	} else {
		fmt.Println("登录失败：姓名不存在")
		return false
	}
}
