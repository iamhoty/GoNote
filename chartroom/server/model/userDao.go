package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gocode/chartroom/common/message"
)

// 我们在服务器启动后，就初始化一个userDao实例
// 把它做成全局的变量，在需要和redis操作时，就直接使用即可
var (
	MyUserDao *UserDao
)

// 定义一个UserDao 结构体
// 完成对User 结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 1. 根据用户id 返回 一个User实例或没有返回err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *message.User, err error) {
	user = &message.User{}
	// 通过id 去redis获取用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在 users 哈希中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	// res {\"userId\":123,\"userPwd\":\"123456\",\"userName\":\"scott\"}
	// 这里我们需要把res 反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// 完成登录的校验 Login
// 1.Login 完成对用户的验证
// 2.如果用户的id和pwd都正确，则返回一个user实例
// 3.如果用户的id或pwd有错误，则返回对应的错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	//先从UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}
	// 证明用户是获取到的
	// 校验密码
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

// 注册
func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao 的连接池中取出一根连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err != nil {
		// 不存在
		if err != ERROR_USER_NOTEXISTS {
			return
		}
	} else {
		err = ERROR_USER_EXISTS
		return
	}
	// 这时，说明id在redis还没有，则可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	// 入库redis
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
