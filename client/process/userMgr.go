package process

import (
	"HugeChattingSystem/client/model"
	"HugeChattingSystem/common/message"
	"fmt"
)

// 客户端也要维护一个在线的用户的map
var onlineUsers = make(map[int]*message.User, 10)
var CurUser model.CurUser // 我们在用户登录成功后，完成对curUser的初始化

// 显示当前在线的用户
func outputOnlineUser() {
	fmt.Println("当前在线用户列表：")
	for id, _ := range onlineUsers {
		fmt.Println("用户id:\t", id)
	}
	fmt.Println("")
}

// 编写一个方法，处理返回的NotifyUserStatusMes

func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	// 适当地优化一下
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}
