package message

// 确定一些消息类型
const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"

	RegisterMesType    = "RegisterMes"
	RegisterResMesType = "RegisterResMes"

	NotifyUserStatusMesType = "NotifyUserStatusMes"

	SmsMesType = "SmsMes"
)

const (
	UserOnline = iota
	UserOffLine
	UserBusyStatus
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
	Code    int    `json:"code"`    // 500表示用户未注册，200表示登录成功
	UserIds []int  `json:"userIds"` // 增加字段，保存用户id的切片
	Error   string `json:"error"`
}

type RegisterMes struct {
	User User `json:"user"` // 类型就是user结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"` // 400 已经占有，200表示成功
	Error string `json:"error"`
}

// NotifyUserStatusMes 为了配合服务器端推送用户状态变化的消息，定义一个类型
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

type SmsMes struct {
	Content string `json:"content"`
	User           // 匿名结构体
}