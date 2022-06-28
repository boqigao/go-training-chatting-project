package process

import (
	"HugeChattingSystem/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) {
	var smsMes message.SmsMes

	json.Unmarshal([]byte(mes.Data), &smsMes)

	fmt.Println(smsMes.UserId, "对大家说：", smsMes.Content)
}
