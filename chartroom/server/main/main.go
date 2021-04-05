package main

import (
	"fmt"
	"gocode/chartroom/server/model"
	"net"
	"time"
)

func init() {
	// 当服务器启动时 就去初始化redis连接池
	initPool("127.0.0.1:6379", 16, 0, 300*time.Second)
	initUserDao()
}

// 处理客户端的通讯
func process(conn net.Conn) {
	fmt.Println("客户端的地址:", conn.RemoteAddr().String())
	defer conn.Close()
	// 这里调用总控, 创建一个
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误=err", err)
		return
	}
}

// 这里我们编写一个函数，完成对UserDao的初始化任务
func initUserDao()  {
	// 这里的pool 本身就是一个全局的变量 在redis中 *redis.pool
	// 这里需要注意一个初始化顺序问题 先initPool初始化pool 再initUserDao
	model.MyUserDao = model.NewUserDao(pool)
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
