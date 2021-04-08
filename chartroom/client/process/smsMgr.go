package process

import (
	"encoding/json"
	"fmt"
	"gocode/chartroom/common/message"
)

func outputGroupMes(mes *message.Message) { // 这个地方mes类型一定SmsMes
	// 反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("outputGroupMes json.Unmarshal err=", err.Error())
		return
	}
	// 显示信息
	content := smsMes.Content
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s", smsMes.UserId, content)
	fmt.Println(info)
	fmt.Println()
}
