package model

import (
	"HugeChattingSystem/common/message"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// MyUserDao 我们在服务器启动时，就初始化一个userDao对象，把它做成全局便利那个，需要和redis操作时候，就直接使用
var (
	MyUserDao *UserDao
)

// UserDao 定义一个userDao结构体，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// NewUserDao 使用工厂模式创建一个UserDao的对象
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 在这里非常重要，在go里面如果我们用这种返回方式，说明已经定义了返回的变量
// 所以我们无需再定义新的返回变量
// 非常重要，不要忘记，所以在前面的很多method里面，直接return就行了，因为err已经定了
func (ud UserDao) getUserByUd(conn redis.Conn, id int) (user *message.User, err error) {

	// 通过给定的id去redis里面查询用户
	res, err := redis.String(conn.Do("HGet", "users", id))

	if err != nil {
		if err == redis.ErrNil {
			// 代表没查到，id不存在
			err = ERROR_USER_NOT_EXISTS
		}
		return
	}
	user = &message.User{} // or user = new(User)
	// 这里我们需要把res反序列化成为一个user对象
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err = ", err)
		return
	}
	return
}

// Login 完成对用户的验证
// 1。 完成对用户的验证
// 2。 如果用户的id和pwd都正确，则返回一个user对象
// 3。 如果用户的id或者pwd错误，则返回一个错误信息
func (ud *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	conn := ud.pool.Get()
	defer conn.Close()
	user, err = ud.getUserByUd(conn, userId)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (ud *UserDao) Register(user *message.User) (err error) {
	conn := ud.pool.Get()

	defer conn.Close()

	_, err = ud.getUserByUd(conn, user.UserId)
	if err == nil {
		// 如果没有错误，反而证明user已经存在
		err = ERROR_USER_EXISTS
		return
	}

	// 注册用户
	data, err := json.Marshal(user)

	if err != nil {
		return err
	}

	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("注册用户错误 err = ", err)
		return
	}
	return
}
