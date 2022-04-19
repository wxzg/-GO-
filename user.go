package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C 	 chan string
	conn net.Conn
	s    *Server
}

// CreateNewUser 创建一个新用户
func CreateNewUser(conn net.Conn, server *Server) *User{

	//拿到用户的网络地址
	getAddr := conn.RemoteAddr().String()
	user := &User{
		Name: getAddr,
		Addr: getAddr,
		C	: make(chan string),
		conn: conn,
		s	: server,
	}

	//监听消息，一旦有消息就发送给用户
	go user.linstenInformation()

	return user
}

func (user *User)linstenInformation(){
	for{
		msg := <- user.C

		user.conn.Write([]byte(msg+"\n"))
	}
}

func (user *User)ListenUserMessage(){

}

// 用户上线
func (user *User)UserOnline(){

	user.s.Mute.Lock()
	user.s.UsersInfo[user.Name] = user
	user.s.Mute.Unlock()
}

//用户下线
func (user *User)Offline(){
	serve := user.s
	serve.Mute.Lock()
	delete(serve.UsersInfo, user.Name)
	close(user.C)
	serve.Mute.Unlock()
}


//发送信息给指定用户
func (user *User) SendMessage(msg string){
	message := "["+ user.Addr +"]"+user.Name+":"+msg
	user.C <- message
}