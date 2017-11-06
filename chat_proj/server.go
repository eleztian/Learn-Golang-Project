/**
  * email: eleztian@gmail.com
  * Blog:www.eleztian.xyz
***/

package main

import (
	"fmt"
	"net"
	"strings"
)

func init() {
	server.clients = make(map[string]Client, 10240)
	server.msgchan = make(chan Msg, 10240)
}

func (server *Server) Recv(client *Client) {
	defer client.conn.Close()
	buf := make([]byte, 1024)
	var msg string
	for {
		length, err := client.conn.Read(buf)
		if checkError(err, "Connection") == false {
			fmt.Println(client.user.name, client.user.addr, " leaved")
			break
		}
		if length > 0 {
			buf[length] = 0
		} else {
			fmt.Println(client.user.name, " recv nothing, continue")
			continue
		}

		if buf[0] == '/' { // cmd
			msg = string(buf[1 : length-1])
			climsg, _ := server.CmdParse(msg, client)
			server.msgchan <- climsg
		} else {
			msg = string(buf[0:length])
			var climsg = Msg{client.user, nil, msg}
			server.msgchan <- climsg
		}
	}
}

func (server *Server) CmdParse(cmdstr string, client *Client) (climsg Msg, err error) {
	cmd := strings.Split(cmdstr, " ")
	var msg string
	msg = "Not found the cmd.\n"
	climsg.from = client.user
	climsg.to = nil
	switch cmd[0] {
	case "username":
		msg = fmt.Sprintf("%s changed name to %s\n", client.user.name, cmd[1])
		delete(server.clients, client.user.name)
		client.user.name = cmd[1]
		server.clients[client.user.name] = *client
	case "to":
		if toClient, ok := server.clients[cmd[1]]; ok {
			msg = strings.Join(cmd[2:], " ") + "\n"
			climsg.to = toClient.user
		} else {
			climsg.to = client.user
			msg = fmt.Sprintf("zhe name %s not exit.\n", cmd[1])
		}
	}
	climsg.msg = msg
	return
}

func (server *Server) GetClient(user *User) (client Client, ok bool) {
	if user == nil {
		ok = false
		return
	}
	client, ok = server.clients[user.name]
	return
}

func (server *Server) Send(climsg *Msg) (err error) {
	client, ok := server.GetClient(climsg.to)
	if ok {
		conn := client.conn
		_, err = conn.Write([]byte(climsg.msg))
		if err != nil {
			fmt.Println(err.Error())
			msg := fmt.Sprintf("%s leaved.\n", climsg.to.name)
			var remsg Msg
			remsg.from = climsg.to
			remsg.to = climsg.from
			remsg.msg = msg
			server.msgchan <- remsg
			conn.Close()
			delete(server.clients, climsg.to.name)
		}
	}
	return
}

func (server *Server) Listen(addr string) (err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	fmt.Println("listening...")
	for {
		conn, err := l.Accept()
		checkError(err, "Accept")
		var client Client
		var user User
		client.user = &user
		client.user.addr = conn.RemoteAddr().String()
		client.conn = conn
		client.user.name = client.user.addr
		server.clients[client.user.name] = client
		fmt.Println("Accept", client.user.addr)
		go server.Recv(&client)
	}

	return
}

func (server *Server) ProcessMsg() {
	for {
		climsg := <-server.msgchan

		str := fmt.Sprintf("[%s]:%s", climsg.from.name, climsg.msg)
		fmt.Print(str)
		climsg.msg = str
		if climsg.to != nil {
			server.Send(&climsg)
		} else {
			for _, client := range server.clients {
				if client.user.name != climsg.from.name {
					climsg.to = client.user
					server.Send(&climsg)
				}
			}
		}
	}
}

func StartServer(port string) {
	fmt.Println("server start to run")
	addr := ":" + port
	go server.ProcessMsg()
	go server.Listen(addr)
}
