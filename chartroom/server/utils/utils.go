package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	xerrors "github.com/pkg/errors"
	"gocode/chartroom/common/message"
	"net"
)

// 这里的方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte // 数组 使用时当切片
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("等待读取...")
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		return mes, xerrors.Wrap(err, "ReadPkg: read pkg header error")
	}
	fmt.Println("读取到buf", this.Buf[:4])
	// 根据buf[:4] 转成uint32
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	// 根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return mes, xerrors.Wrap(err, "ReadPkg: read pkg body error")
	}
	// 把pkg反序列化 -> Message  注意传指针
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		return mes, xerrors.Wrap(err, "ReadPkg: json.Unmarshal fail")
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	// 1.先发送长度
	var pkgLen uint32 // 定义的包的长度
	pkgLen = uint32(len(data))
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen) // 将数字转为byte切片
	// 发送长度 n是发送成功的字节数
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		return xerrors.Wrap(err, "WritePkg: conn.Write(bytes) fail")
	}
	fmt.Printf("发送消息的长度=%d 内容=%s\n", len(data), string(data))
	// 2.发送消息本身
	n, err = this.Conn.Write(data)
	if err != nil || n != int(pkgLen) {
		return xerrors.Wrap(err, "WritePkg: conn.Write(bytes) fail")
	}
	return
}
