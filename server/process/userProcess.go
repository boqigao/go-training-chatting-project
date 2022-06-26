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
	//
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

	// 如果用户的id为100，密码是123456就是合法的，否则不合法
	if err == nil {
		// 合法
		loginResMes.Code = 200
		loginResMes.Error = ""
	} else {
		// 不合法
		loginResMes.Code = 500
		loginResMes.Error = "用户不存在，需要注册再使用"
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
