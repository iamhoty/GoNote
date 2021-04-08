package process

import (
	"fmt"
	"gocode/chartroom/client/model"
	"gocode/chartroom/common/message"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser // 用户登录成功后 完成对CurUser的初始化

// 在客户端显示当前在线的用户
func outputOnlineUser()  {
	fmt.Println("当前在线用户列表:")
	// 遍历
	for id, _ := range onlineUsers{
		fmt.Println("用户id:\t", id)
	}
}

// 编写一个方法，处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes)  {
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 原来没有 才创建
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,

		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}