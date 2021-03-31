package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gocode/chartroom/common/message"
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
	fmt.Println("读取到buf消息长度", buf[:4])
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
