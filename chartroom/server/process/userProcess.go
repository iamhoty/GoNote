package process2

import (
	"encoding/json"
	"fmt"
	xerrors "github.com/pkg/errors"
	"gocode/chartroom/common/message"
	"gocode/chartroom/server/model"
	"gocode/chartroom/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
}

// 专门处理注册请求
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1.从mes中 取出mes.Data 并直接反序列化成LoginMes
	var RegisterMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessRegister: json.Unmarshal fail")
	}
	// 2.声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	// 3.声明一个登录返回结果
	var RegisterResMes message.RegisterResMes

	// 我们需要到redis数据库去完成注册
	// 1.使用model.MyUserDao 到redis去注册
	err = model.MyUserDao.Register(&RegisterMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			RegisterResMes.Code = 505 // 表示用户存在
			RegisterResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			RegisterResMes.Code = 506
			RegisterResMes.Error = "注册发生未知错误..."
		}
	} else {
		RegisterResMes.Code = 200
		fmt.Println(&RegisterMes.User, "注册成功!")
	}
	// 序列化
	data, err := json.Marshal(RegisterResMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessRegister: json.Marshal fail")
	}
	// 4.将data赋值给resMes
	resMes.Data = string(data)
	// 序列化 data为byte切片
	data, err = json.Marshal(resMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessRegister: json.Marshal fail")
	}
	// 5.发送
	// 创建Transfer实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	return tf.WritePkg(data)
}

// 专门处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 内容={"type":"LoginMes","data":"{\"userId\":1,\"userPwd\":\"123\",\"userName\":\"\"}"}
	// 1.从mes中 取出mes.Data 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessLogin: json.Unmarshal fail")
	}
	// 2.声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	// 3.声明一个登录返回结果
	var loginResMes message.LoginResMes

	// 我们需要到redis数据库去完成验证.
	// 1.使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500 // 表示用户不存在
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}

	} else {
		loginResMes.Code = 200
		fmt.Println(user, "登录成功!")
	}

	// 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessLogin: json.Marshal fail")
	}
	// 4.将data赋值给resMes
	resMes.Data = string(data)
	// 序列化 data为byte切片
	data, err = json.Marshal(resMes)
	if err != nil {
		return xerrors.Wrap(err, "ServerProcessLogin: json.Marshal fail")
	}
	// 5.发送
	// 创建Transfer实例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	return tf.WritePkg(data)
}
