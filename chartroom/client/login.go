package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gocode/chartroom/common/message"
	"net"
)

func login(userId int, userPwd string) (err error) {
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
	// 7. 到这个时候 data就是我们要发送的消息
	// 7.1 先把 data的长度发送给服务器 write接收的是byte切片
	// 先获取到 data的长度->转成一个表示长度的byte切片
	var pkgLen uint32 // 定义的包的长度
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen) // 将数字转为byte切片
	// 发送长度 n是发送成功的字节数
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Printf("客户端，发送消息的长度=%d 内容=%s\n", len(data), string(data))
	// 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	// 这里还需要处理服务器端返回的消息
	mes, err = readPkg(conn) // mes 就是
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}
	// 将mes的Data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		fmt.Println("登录成功")
	} else if loginResMes.Code == 500 {
		fmt.Println(loginResMes.Error)
	}
	return
}
