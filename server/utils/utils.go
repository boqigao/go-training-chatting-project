package utils

import (
	"HugeChattingSystem/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type Transfer struct {
	// 分析他应该有哪些字段
	Conn   net.Conn
	Buffer [8094]byte //这是传输时候使用的缓冲
}

// ReadPkg readPkg整个封装出来，实现了两个功能，读取长度以后，读取所需长度的message
// return 读取出来的message，还有error
func (t *Transfer) ReadPkg() (mes message.Message, err error) {

	// 这个read是读取的消息长度
	// conn.Read 在conn没有被关闭的情况下才会阻塞
	// 如果客户端关闭了，就不会阻塞
	n, err := t.Conn.Read(t.Buffer[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Read err = ", err)
		return
	}
	fmt.Println("读到的buf=", t.Buffer[0:4])
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(t.Buffer[0:4])

	// 这个read是读取的消息实际内容
	n, err = t.Conn.Read(t.Buffer[:pkgLen])

	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read(buffer[:pkgLen]) fail, err = ", err)
		return
	}

	// 根据pkgLen的长度将把buf反序列化 - > message.Message
	err = json.Unmarshal(t.Buffer[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err = ", err)
		return
	}

	return

}

// WritePkg writePkg也整个封装出来，write的是一个切片，其实叫writePkg更有点像sendPkg
func (t *Transfer) WritePkg(data []byte) (err error) {

	// 先发送一个长度给对方，把长度转成了一个bytes
	// 发送长度的原因是，对面是一个buffer不知道长度的话会读出乱码，所以需要截取固定的长度
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	// 把数组转换成切片
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)

	// 发送长度，也是要一个切片
	n, err := t.Conn.Write(bytes[0:4])

	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail err = ", err)
		return err
	}

	// 发送数据本身
	n, err = t.Conn.Write(data)

	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail err = ", err)
		return err
	}
	return
}
