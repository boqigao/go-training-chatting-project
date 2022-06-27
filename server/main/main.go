package main

import (
	"HugeChattingSystem/server/model"
	"fmt"
	"net"
	"time"
)

// 处理和客户端的通讯
func process(conn net.Conn) {
	// 读客户端发送的信息
	defer conn.Close()
	p := Processor{
		Conn: conn,
	}
	p.Process3()
}

func init() {
	// 当服务器启动时候我们就去初始化连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	// Dao只需要启动一次，有点类似于单例
	initUserDao()
}

func main() {
	fmt.Println("服务器[新的结构]在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err = ", err)
	}
	defer listen.Close()

	for {
		fmt.Println("等待客户端连接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err)
		}

		//一旦连接成功，则启动一个写成和客户端保持通讯
		go process(conn)
	}
}

// 这里我们编写一个函数，完成对userDao的初始化
func initUserDao() {
	//这里的pool本身就是一个全局的变量，这里需要一个初始化顺序问题
	model.MyUserDao = model.NewUserDao(pool)
}
