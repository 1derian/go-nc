package proto

// 该包自定义了协议,用来解决粘包问题

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"strings"
)

// Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	// message 接收到的数据
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头 , 内容是message的长度
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体 , 内容是数据内容,字节类型
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	// 返回封包好的数据
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 读取 4 个字节的数据，文件读偏移量不动，返回的字节切片是内部 buf 的引用
	lengthByte, _ := reader.Peek(4)  // 是一个 lengthByte []byte 类型

	// 初始化一个新的 buf , 参数是一个字节类型 , 大小为 需要接收的数据的 长度
	lengthBuff := bytes.NewBuffer(lengthByte)
	// 初始化一个 length 变量
	var length int32
	// 一句话就是字节转成int类型
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// Buffered返回缓冲区中现有的可读取的字节数。

	if length+4 > int32(reader.Buffered()) {
		// 当真正的要接受的数据(全部)大于缓冲区,你必须循环接收
		// 思路: 循环接收的读取,然后读取的内容转为字符串拼接,然后判断剩余内容
		// length+4  是数据的总大小
		var recvSize int32
		var builder strings.Builder
		s := int32(4096)  // 接收的大小
		for  recvSize < length+4 {
			if recvSize + int32(4096) > length+4{
				s = length + 4 - recvSize
			}
			pack := make([]byte,s) // 这个要动态的变化,默认都是4096,最后一次要是剩余的大小.这是和python唯一的区别.go会补齐,python不会
			_, err = reader.Read(pack)
			if err != nil {
				return "", err
			}
			builder.WriteString(strings.TrimSpace(string(pack[4:])))
			recvSize += int32(4096)
		}
		return builder.String(), nil
	}

	// 读取真正的消息数据 , reader.Buffered() == length+4
	// 当接收的数据没有缓冲区大,直接读取即可
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	// 返回真正的数据内容,是一个切片
	return string(pack[4:]), nil

	// 这个思路是这样的,先读取4个字节 , 不移动偏移量
	// 然后这四个字节的内容是字节类型的, 字节-->整形,借助
}