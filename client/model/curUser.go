package model

import (
	"HugeChattingSystem/common/message"
	"net"
)

// 因为在客户端在很多地方会使用到curUser,我们将其作为一个全局变量

type CurUser struct {
	Conn net.Conn
	message.User
}
