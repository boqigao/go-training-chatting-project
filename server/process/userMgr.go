package process

import "fmt"

// 因为UserMgr对象在服务器端有且只有一个
// 而且在很多地方都会使用到，因此我们将其定义为全局变量

var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUser的添加
func (um *UserMgr) AddOnlineUser(up *UserProcess) {
	um.onlineUsers[up.UserId] = up
}

// DelOnlineUser 删除特定用户
func (um *UserMgr) DelOnlineUser(userId int) {
	delete(um.onlineUsers, userId)
}

// GetAllOnlineUser 得到所有在线用户
func (um *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return um.onlineUsers
}

// GetOnlineUserById 获得某个特定用户的up，其实就是为了得到其中的conn
func (um *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	up, ok := um.onlineUsers[userId]

	if !ok {
		// 说明现在查找的用户当前不在线
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}

	return
}
