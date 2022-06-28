package process

import (
	"HugeChattingSystem/client/utils"
	"HugeChattingSystem/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要任何字段
}

func (up *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return err
	}
	defer conn.Close()

	var mes message.Message
	mes.Type = message.RegisterMesType
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserName = userName
	registerMes.User.UserPwd = userPwd

	// 序列化具体的mes
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	// 序列化整个mes
	mes.Data = string(data)

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return err
	}

	transfer := utils.Transfer{
		Conn:   conn,
		Buffer: [8094]byte{},
	}

	err = transfer.WritePkg(data)

	if err != nil {
		fmt.Println("注册信息发送错误，err = ", err)
	}

	// 处理服务器返回的数据
	mes, err = transfer.ReadPkg()

	if err != nil {
		fmt.Println("服务器返回消息（注册信息）错误， err = ", err)
	}

	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)

	if registerResMes.Code == 200 {
		fmt.Printf("恭喜您！注册新用户%s成功，请您登录。\n", userName)
		os.Exit(0)
	} else {
		if registerResMes.Code == 400 {
			fmt.Printf("用户名%s已经占有，注册失败", userName)
			fmt.Printf("错误代码%v, 错误内容%s", registerResMes.Code, registerResMes.Error)
		}
		os.Exit(0)
	}

	return
}

// Login 写一个函数完成登陆校验
func (up *UserProcess) Login(userId int, userPwd string) (err error) {
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

		// 可以显示当前在线的用户列表
		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UserIds {
			fmt.Println("用户id:\t", v)
			// 把它扔到客户端维护的onlineUsers里面
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user

		}
		fmt.Println("")

		// 这里我们还需要在客户端启动一个协程
		// 该写成保持和服务器端的通讯，如果服务器端有东西推给客户端
		// 这个协程负责
		go serverProcessMes(conn)

		// 显示我们登录成功后的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Printf("登录失败，错误代码为%d, 错误详细信息为%s\n", loginResMes.Code, loginResMes.Error)
	}
	return
}
