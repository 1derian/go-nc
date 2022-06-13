package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"nc/proto"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("h","127.0.0.1","主机名")
	//            类型        短变量      默认值       使用帮助
	port := flag.Int("p",8889,"端口")
	flag.Parse()
	// 格式化字符串
	hostPort := fmt.Sprintf("%s:%d",*host,*port)
	// 1.建立链接
	conn, err := net.Dial("tcp", hostPort)
	// 2.如果建立链接失败 , 退出
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	fmt.Println("链接成功...")
	// 3.defer 关闭链接
	defer conn.Close()
	// 创建一个 inputReader 对象
	inputReader := bufio.NewReader(os.Stdin) // 从本地io中获取数据
	clientReader := bufio.NewReader(conn)
	for {
		// 读取到换行结束
		fmt.Printf("shell>")
		// 阻塞住 , 接收用户输入的数据,然后赋值给input
		// 实际上是 把终端接收到用户输入的内容写到一个文本,然后inputReader去读取这个文本的内容
		input, _ := inputReader.ReadString('\n') // 读取用户输入
		// 去除空白
		inputInfo := strings.Trim(input, "\r\n")
		if strings.ToUpper(inputInfo) == "QUITE" { // 如果输入q就退出
			os.Exit(0)
		}
		data, err := proto.Encode(inputInfo)
		// 如果封包失败,退出
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		// 发送封包好的数据
		conn.Write(data)

		// 接收服务端返回的命令执行的结果
		msg, err := proto.Decode(clientReader)
		// 7.如果报错,退出
		if err != nil {
			if err == io.EOF {
				os.Exit(0)
			}
			log.Print(err)
			os.Exit(1)
		}
		//8.输出服务端回复的数据
		fmt.Println(msg)
	}
	// go build -ldflags="-H windowsgui -w -s" nc_client.go
}
