package process

import (
	"HugeChattingSystem/client/utils"
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

	fmt.Scanln(&key)
	switch key {
	case 1:
		fmt.Println("显示在线用户列表")
	case 2:
		fmt.Println("发送消息")
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
	}
}
