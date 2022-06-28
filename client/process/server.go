package process

import (
	"HugeChattingSystem/client/utils"
	"HugeChattingSystem/common/message"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

// ShowMenu 显示登录成功之后的界面
func ShowMenu() {

	fmt.Println("----------恭喜登录成功----------")
	fmt.Println("----------1。显示在线用户列表----------")
	fmt.Println("----------2。发送消息----------")
	fmt.Println("----------3。信息列表----------")
	fmt.Println("----------4。退出系统----------")
	fmt.Println("请选择1-4")
	var key int
	var content string

	smsProcess := &SmsProcess{}

	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("请输入你想对大家说点什么？")
		fmt.Scanln(&content)
		smsProcess.SendGroupProcess(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("输入错误，请重新选择")
	}
}

// 和服务器端保持通讯
func serverProcessMes(conn net.Conn) {
	// 创建一个transfer，让他不停的读取服务器的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() fail, err = ", err)
		}
		fmt.Println("mes =", mes)

		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 有人上线了
			// 1. 取出notifyMessage
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 2. 把这个用户的信息，状态，保存在客户端的map中（？）
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			// 返回了一个暂时不能识别的消息
			fmt.Println("服务器返回了未知的消息类型")

		}
	}
}
