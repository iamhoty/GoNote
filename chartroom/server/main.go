package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gocode/chartroom/common/message"
	"io"
	"net"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("等待读取...")
	_, err = conn.Read(buf[:4])
	if err != nil {
		// 自定义错误
		//err = errors.New("read pkg header error")
		return
	}
	fmt.Println("读取到buf", buf[:4])
	// 根据buf[:4] 转成uint32
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[:4])
	// 根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	// 把pkg反序列化 -> Message  注意传指针
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		err = errors.New("json.Unmarshal error")
		return
	}
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	// 1.先发送长度
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
	fmt.Printf("发送消息的长度=%d 内容=%s\n", len(data), string(data))
	// 2.发送消息本身
	n, err = conn.Write(data)
	if err != nil || n != int(pkgLen) {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	return
}

// 专门处理登录请求
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	// 内容={"type":"LoginMes","data":"{\"userId\":1,\"userPwd\":\"123\",\"userName\":\"\"}"}
	// 1.从mes中 取出mes.Data 并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 2.声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType
	// 3.声明一个登录返回结果
	var loginResMes message.LoginResMes
	// 如果用户id= 100， 密码=123456, 认为合法，否则不合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		loginResMes.Code = 200
	} else {
		// 不合法
		loginResMes.Code = 500 // 表示用户不存在
		loginResMes.Error = "该用户不存在, 请注册再使用..."
	}
	// 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	// 4.将data赋值给resMes
	resMes.Data = string(data)
	// 序列化 data为byte切片
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
		return
	}
	// 5.发送
	err = writePkg(conn, data)
	return
}

// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func serverProcess(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		err = serverProcessLogin(conn, mes)
	case message.RegisterMesType:
		// 处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// 处理客户端的通讯
func process(conn net.Conn) {
	fmt.Println("客户端的地址:", conn.RemoteAddr().String())
	defer conn.Close()
	// 读取客户端发送的消息
	for {
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务器也退出...")
				return
			} else {
				fmt.Println("readPkg fail err=", err)
				return
			}
		}
		fmt.Println("mes==", mes)
		err = serverProcess(conn, &mes)
		if err != nil {
			return
		}
	}
}

func main() {

	fmt.Println("服务器在8889端口监听")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("监听失败, err", err)
		return
	}
	// 监听成功 等待客户端连接
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("客户端连接失败, err", err)
			continue
		}
		// 一旦连接成功 则启动一个协程和客户端保持通讯
		go process(conn)

	}

}
