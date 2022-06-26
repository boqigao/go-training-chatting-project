package main

import (
	"HugeChattingSystem/client/process"
	"fmt"
	"os"
)

// 定义两个全局变量，一个表示用户的id，一个表示用户的密码
var userId int
var userPwd string

func main() {
	//接受用户的选择
	//判断是否还继续选择菜单

	var key int
	var loop = true

	for loop {
		fmt.Println("----------欢迎进入多人聊天系统----------")
		fmt.Println("\t\t\t 1. 登录聊天室")
		fmt.Println("\t\t\t 2. 注册用户")
		fmt.Println("\t\t\t 3. 退出系统")
		fmt.Println("\t\t\t 请选择 1 - 3")

		fmt.Scanf("%d\n", &key)

		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户的id号码")
			//这里也可以使用Scanln应该会自动识别
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			up := &process.UserProcess{}
			err := up.Login(userId, userPwd)
			if err != nil {

			}
			// loop = false
		case 2:
			fmt.Println("注册用户")
			// loop = false
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)

		default:
			fmt.Println("您的输入有误，请重新输入")
		}
	}

	//根据用户的输入选择新的提示信息
	if key == 1 {
		// 说明用户要登陆聊天室

	} else {
		fmt.Println("进行用户注册的逻辑")
	}
}
