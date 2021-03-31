package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息的类型
}

type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"`   // 返回状态码 500 表示该用户未注册 200表示登录成功
	UsersId []int  `json:"userId"` // 增加字段，保存用户id的切片
	Error   string `json:"error"`  // 返回错误信息
}
