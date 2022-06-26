package message

// 确定一些消息类型
const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"

	RegisterMesType = "RegisterMes"
)

// Message 可以理解为一个最基础的message类
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

// LoginMes struct of login message of client
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

// LoginResMes struct of login reply message from server
type LoginResMes struct {
	Code  int    `json:"code"` // 500表示用户未注册，200表示登录成功
	Error string `json:"error"`
}

type RegisterMes struct {
}