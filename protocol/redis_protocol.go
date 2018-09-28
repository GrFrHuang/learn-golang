package main

import (
	"fmt"
	"os"
	"net"
	"strconv"
	"github.com/lunny/log"
)

const (
	RedisServerAddress = "127.0.0.1:6380"
	RedisServerNetwork = "tcp"
)

type RedisError struct {
	msg string
}

func (this *RedisError) Error() string {
	return this.msg
}

// 连接到redis server
func conn() (net.Conn, error) {
	conn, err := net.Dial(RedisServerNetwork, RedisServerAddress)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return conn, err
}

// 将参数转化为redis请求协议
func getCmd(args []string) []byte {
	cmdString := "*" + strconv.Itoa(len(args)) + "\r\n"
	for _, v := range args {
		cmdString += "$" + strconv.Itoa(len(v)) + "\r\n" + v + "\r\n"
	}
	cmdByte := make([]byte, len(cmdString))
	copy(cmdByte[:], cmdString)
	log.Info(cmdString)
	return cmdByte
}

func dealReply(reply []byte) (interface{}, error) {
	responseType := reply[0]
	switch responseType {
	case '+':
		return dealStatusReply(reply)
	case '$':
		return dealBulkReply(reply)
	default:
		return nil, &RedisError{"proto wrong!"}
	}
}

// todo 错误响应\整数响应\批量主体响应
// 处理状态响应
func dealStatusReply(reply []byte) (interface{}, error) {
	statusByte := reply[1:]
	pos := 0
	for _, v := range statusByte {
		if v == '\r' {
			break
		}
		pos++
	}
	status := statusByte[:pos]
	return string(status), nil
}

// 处理主体响应
func dealBulkReply(reply []byte) (interface{}, error) {
	statusByte := reply[1:]
	// 获取响应文本第一行标示的响应字符串长度
	pos := 0
	for _, v := range statusByte {
		if v == '\r' {
			break
		}
		pos++
	}
	strlen, err := strconv.Atoi(string(statusByte[:pos]))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if strlen == -1 {
		return "nil", nil
	}
	nextLinePost := 1
	for _, v := range statusByte {
		if v == '\n' {
			break
		}
		nextLinePost++
	}
	result := string(statusByte[nextLinePost:nextLinePost+strlen])
	return result, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("usage: go run proto.go + redis command\nfor example:\ngo run proto.go PING")
		os.Exit(0)
	}
	// GO1.9增加了安全map sync.Map{}
	conn, _ := conn()
	cmd := getCmd(args)
	conn.Write(cmd)
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	res, _ := dealReply(buf[:n])
	fmt.Println("redis的返回结果是 ", res)
}
