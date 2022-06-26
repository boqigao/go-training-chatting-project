package main

import (
	"HugeChattingSystem/common/message"
	process2 "HugeChattingSystem/server/process"
	"HugeChattingSystem/server/utils"
	"fmt"
	"io"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// ServerProcessMes 编写一个serverProcessMes函数
// 根据客户端发送消息种类不同，决定调用哪个函数来处理
func (p *Processor) ServerProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		// 创建一个userProcessor
		up := process2.UserProcess{
			Conn: p.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//
	default:
		fmt.Println("消息类型不存在，无法处理。。。")
	}
	return
}

func (p Processor) Process3() {
	// 循环的读取客户端发送的消息
	transfer := utils.Transfer{
		Conn:   p.Conn,
		Buffer: [8094]byte{},
	}
	for {
		mes, err := transfer.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出")
				return
			} else {
				fmt.Println("readPkg(conn) fail, err = ", err)
				return
			}
		}
		fmt.Println("mes=", mes)
		p.ServerProcessMes(&mes)
	}
}
