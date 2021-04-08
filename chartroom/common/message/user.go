package message

// 定义一个用户的结构体
type User struct {
	// 为了序列化和反序列化成功 定义json反射
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` //用户状态
	Sex        string `json:"sex"`        //性别
}
