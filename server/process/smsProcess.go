package process

import (
	"HugeChattingSystem/client/utils"
	"HugeChattingSystem/common/message"
	"encoding/json"
	"net"
)

type SmsProcess struct {
}

func (sp *SmsProcess) SendGroupMes(mes *message.Message) {

	var smsMes message.SmsMes
	json.Unmarshal([]byte(mes.Data), &smsMes)
	destId := smsMes.UserId

	data, err := json.Marshal(mes)
	if err != nil {
	}

	// 遍历服务器端的在线用户
	for id, up := range userMgr.onlineUsers {
		if id != destId {
			sp.SendToEachOnlineUse(data, up.Conn)
		}
	}
}

func (sp *SmsProcess) SendToEachOnlineUse(content []byte, conn net.Conn) {
	transfer := utils.Transfer{Conn: conn}
	transfer.WritePkg(content)
}
