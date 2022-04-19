package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct{
	Ip string
	Port int
	UsersInfo map[string]*User
	Mute sync.RWMutex
	Message chan string
}

func CreateNewServer(ip string, port int) *Server{
	server := &Server{
		Ip:ip,
		Port: port,
		UsersInfo:make(map[string]*User),
		Message: make(chan string),
	}

	return server
}

// BoardCast 广播到user的C中，当消息发送过去后User对应的消息监听事件会自动讲消息发送
func (s *Server) BoardCast(){
	for  {
		msg := <- s.Message
		//因为要写入数据，所以要加锁
		s.Mute.Lock()
		for _, v := range s.UsersInfo{
			v.C <- msg
		}
		s.Mute.Unlock()
	}
}

// BoardMessage 用来发送广播信息到Message中
func (s *Server) BoardMessage(user *User, msg string){
	message := "[" + user.Addr + "] " + user.Name + " : " + msg
	s.Message <- message
}

//监听user发送的信息并广播
func (s *Server)ListenUserMessage(user *User, conn net.Conn){

	for  {
		msg := make([]byte,4096)
		n, err := conn.Read(msg)

		if n==0{
			s.BoardMessage(user,"下线")
			return
		}

		if err != nil && err != io.EOF {
			fmt.Println("conn read err:",err)
			return
		}

		s.BoardMessage(user, string(msg[:n-1]))
	}
}

func (s *Server)hanldeConnect(conn net.Conn){
	//创建当前用户
	user := CreateNewUser(conn, s)

	//将用户加入表中
	user.UserOnline()
	//广播给其他用户该用户已上线
	s.BoardMessage(user, "已上线")

	go s.ListenUserMessage(user, conn)
	//handle阻塞
	select {

	}
}

func (s *Server) Start(){
	ln, err := net.Listen("tcp",fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err", err)
	}
	// 关闭监听
	defer ln.Close()
	//广播用户上线
	go s.BoardCast()
	fmt.Println("TCP listen已建立")
	for{
		 conn, connErr := ln.Accept()
		 if connErr != nil {
			 fmt.Println("连接建立失败", connErr)
			 continue
		 }
		 //客户端连接建立成功，处理相关业务
		 go s.hanldeConnect(conn)
	}
}