package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"nc/proto"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/*
10.存在的问题,如果对方文本是utf-8编码的,操作系统是windows,这边识别的是windows就全部按照gbk解码了
这不又乱码了吗 , 鸡肋啊
*/

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}

func execCommand(command ,execName,arg string)(result string){
	out, err := exec.Command(execName,arg,command).Output()
	var builder strings.Builder
	if err == nil{
		// 拼接字符串
		res := ConvertByte2String(out,GB18030)
		builder.WriteString(res)
	}else {
		res := fmt.Sprintf("'%s' not found\n",command)
		builder.WriteString(res)
	}
	result = builder.String()
	return result
}



func handle(conn net.Conn){
	// 不要忘记关闭conn
	defer conn.Close()
	// 开始处理conn
	// 1.创建一个网络io对象
	ncReader := bufio.NewReader(conn)
	var commandRes string
	// 2.持续接收命令
	for  {
		// 3.接收消息(客户端发来的命令)
		// 返回字符串和nil
		command, err := proto.Decode(ncReader)
		// 4.判断是否出错
		if err != nil{
			if err == io.EOF {
				return
			}
			log.Println(err)
		}
		// 5.拿到command,执行输出的结果为out
		// 返回一个切片,如果命令运行失败,会返回
		// 判断当前平台
		if runtime.GOOS == "windows" {
			commandRes = execCommand(command,"cmd","/c")
		} else {
			commandRes = execCommand(command,"/bin/sh","-c")
		}

		// 6.把命令执行的结果发送给客户端
		data, err := proto.Encode(commandRes)
		// 如果封包失败,退出
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		// 发送封包好的数据
		// 命令执行是没有问题,即data是ok的
		conn.Write(data)
	}

}


func main()  {
	port := flag.Int("p",8889,"端口")

	flag.Parse()
	hostPort := fmt.Sprintf(":%d",*port)

	// 开启监听
	server,err := net.Listen("tcp",hostPort)
	// 判断监听是否成功
	if err != nil{
		log.Print(err)
		os.Exit(1)  // 1 有错误退出
	}

	defer server.Close() // 退出别忘记关闭链接
	// 死循环 开始接收链接
	for  {
		conn,err := server.Accept()
		if err != nil{
			os.Exit(1)  // 1 有错误退出
		}
		// 开启 goroutine  处理conn
		go handle(conn)

	}



}
