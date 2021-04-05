package process

import (
	"encoding/json"
	"fmt"
	"gocode/chartroom/client/utils"
	"gocode/chartroom/common/message"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要字段..
}

func (this *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	// 1.连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("连接服务器失败")
		return err
	}
	defer conn.Close()
	// 2.准备发消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3.创建一个登录消息
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	// 4.将registerMes序列化 返回byte切片
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json序列化失败, err", err)
		return err
	}
	// 5.把data赋给 mes.Data字段
	mes.Data = string(data)
	// 6.将mes进行序列化化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json序列化mes失败, err", err)
		return err
	}
	// 创建Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 7. 到这个时候 data就是我们要发送的消息
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
		return
	}

	mes, err = tf.ReadPkg() // mes 就是 RegisterResMes
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	// 将mes的Data部分反序列化成 RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功, 你重新登录一把")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return

}

func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	// 1.连接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("连接服务器失败")
		return err
	}
	defer conn.Close()
	// 2.准备发消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 3.创建一个登录消息
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 4.将loginMes序列化 返回byte切片
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json序列化失败, err", err)
		return err
	}
	// 5.把data赋给 mes.Data字段
	mes.Data = string(data)
	// 6.将mes进行序列化化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json序列化mes失败, err", err)
		return err
	}
	// 创建Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 7. 到这个时候 data就是我们要发送的消息
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误 err=", err)
		return
	}

	mes, err = tf.ReadPkg() // mes 就是 LoginResMes
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	// 将mes的Data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//这里我们还需要在客户端启动一个协程
		//该协程保持和服务器端的通讯.如果服务器有数据推送给客户端
		//则接收并显示在客户端的终端
		go serverProcessMes(conn)

		// 显示登录成功的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
