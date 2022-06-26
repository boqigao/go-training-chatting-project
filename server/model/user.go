package model

// User 为了序列化或者反序列化成功，我们必须保证用户信息的字符串对应，一定要tag
type User struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}
