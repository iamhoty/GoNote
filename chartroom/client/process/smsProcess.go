package process

import (
	"encoding/json"
	"fmt"
	"gocode/chartroom/client/utils"
	"gocode/chartroom/common/message"
)

type SmsProcess struct {
}

// 发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	// 1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType
	// 2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content // 发送内容
	smsMes.UserId = CurUser.UserId // 从CurUser全局变量获取userid
	smsMes.UserStatus = CurUser.UserStatus
	// 3.序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal smsMes fail =", err.Error())
		return
	}
	mes.Data = string(data)
	// 4.对mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal mes fail =", err.Error())
		return
	}
	// 5.群聊消息发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	// 6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes send mes err=", err.Error())
		return
	}
	return
}
