package process

import (
	"HugeChattingSystem/common/message"
	"HugeChattingSystem/server/model"
	"HugeChattingSystem/server/utils"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	// 分析需要哪些字段
	Conn net.Conn
	// 增加一个字段，表示该Conn是哪个用户的
	UserId int
}

// NotifyOthersOnlineUser 编写通知所有在线用户的方法
func (up *UserProcess) NotifyOthersOnlineUser(userId int) (err error) {
	// 遍历所有在线用户,然后一个一个发送NotifyMes

	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id != userId {
			// 开始通知【单独写一个方法】
			// 这个up是所有的online的user的socket，并不是最开始在login的哪个用户的socket
			up.NotifyMeToOthers(userId)
		}
	}

	return err
}

func (up *UserProcess) NotifyMeToOthers(userId int) {
	// 组装我们的消息
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	data, err := json.Marshal(notifyUserStatusMes)

	if err != nil {
		fmt.Println(err)
		return
	}

	mes.Data = string(data)
	data, err = json.Marshal(mes)

	transfer := utils.Transfer{Conn: up.Conn}

	transfer.WritePkg(data)
}

func (up *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &registerMes) fail, err = ", err)
		return err
	}

	// 先声明一个resMessage
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	err = model.MyUserDao.Register(&registerMes.User)

	if err != nil {
		return
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		return
	}

	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail, err = ", err)
	}

	transfer := utils.Transfer{
		Conn: up.Conn,
	}
	err = transfer.WritePkg(data)
	return err

}

// ServerProcessLogin 编写一个函数serverProcessLogin，专门处理登录请求
func (up *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 先从message中取出mes.Data, 并且直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal([]byte(mes.Data), &loginMes) fail, err = ", err)
		return err
	}

	// 先声明一个resMessage
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	var loginResMes message.LoginResMes

	// 我们需要去用redis数据库完成验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err == nil {
		// 合法
		loginResMes.Code = 200
		loginResMes.Error = ""
		up.UserId = loginMes.UserId
		userMgr.AddOnlineUser(up)

		// 告诉在线的用户当前用户上线了
		up.NotifyOthersOnlineUser(loginMes.UserId)

		// 告诉当前登录的用户谁在线
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserIds = append(loginResMes.UserIds, id)
		}
	} else {
		// 不合法
		if err == model.ERROR_USER_NOT_EXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "未知错误"
		}
	}

	// 将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal(loginResMes) fail, err = ", err)
		fmt.Println(user, "登录成功")
	}

	resMes.Data = string(data)

	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal(resMes) fail, err = ", err)
	}

	// 发送的函数, 我们将其封装到一个writePkg函数中
	// 因为使用了分层模式，我们先创建一个Transfer
	transfer := utils.Transfer{
		Conn: up.Conn,
	}
	err = transfer.WritePkg(data)
	return err
}
