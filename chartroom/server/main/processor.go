package main

import (
	"fmt"
	"gocode/chartroom/common/message"
	"gocode/chartroom/server/process"
	"gocode/chartroom/server/utils"
	"net"
)

type Processor struct {
	Conn net.Conn
}

// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func (this *Processor) ServerProcess(mes *message.Message) (err error) {
	// 看看是否能接收到客户端发送的群发的消息
	fmt.Println("mes=", mes)

	switch mes.Type {
	case message.LoginMesType:
		// 创建UserProcess实例
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 处理注册
		up := &process2.UserProcess{
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		// 创建一个SmsProcess实例完成转发群聊消息
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (this *Processor) process2() (err error) {
	for {
		// 创建Transfer 完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		// 读取客户端发送的消息
		mes, err := tf.ReadPkg()
		if err != nil {
			return err
		}
		fmt.Println("mes==", mes)
		err = this.ServerProcess(&mes)
		if err != nil {
			return err
		}
	}
	return
}
