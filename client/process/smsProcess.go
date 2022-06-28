package process

import (
	"HugeChattingSystem/client/utils"
	"HugeChattingSystem/common/message"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

// 发送群聊的消息

func (sp *SmsProcess) SendGroupProcess(content string) (err error) {

	// 1 创建一个Msg
	var mes message.Message
	mes.Type = message.SmsMesType

	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	data, err := json.Marshal(smsMes)

	if err != nil {
		fmt.Println(err)
	}

	mes.Data = string(data)

	data, err = json.Marshal(mes)

	if err != nil {
		fmt.Println(err)
		return
	}
	trans := utils.Transfer{Conn: CurUser.Conn}
	trans.WritePkg(data)
	return
}
