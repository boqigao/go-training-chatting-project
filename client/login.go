package main

import (
	"HugeChattingSystem/client/utils"
	"HugeChattingSystem/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

// 写一个函数完成登陆校验

func login(userId int, userPwd string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return err
	}
	defer conn.Close()

	// 创建一个通用的message结构体
	var mes message.Message
	mes.Type = message.LoginMesType

	// 创建一个送信专用的message结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	// 先发送一个长度给对方，把长度转成了一个bytes
	// 发送长度的原因是，对面是一个buffer不知道长度的话会读出乱码，所以需要截取固定的长度
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	// 把数组转换成切片
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)

	// 发送长度，也是要一个切片
	n, err := conn.Write(bytes[0:4])

	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail err = ", err)
		return err
	}

	fmt.Println("客户端发送消息的长度成功, len = ", len(data))
	fmt.Println("客户端发送的内容为：", string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail, err = ", err)
		return err
	}

	// 这里还要处理服务器返回的数据
	transfer := utils.Transfer{
		Conn:   conn,
		Buffer: [8094]byte{}, //这是传输时候使用的缓冲
	}
	mes, err = transfer.ReadPkg()
	if err != nil {
		fmt.Println("客户端读取服务端消息失败, err = ", err)
	}

	// 将其反序列化为LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("login success!")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
