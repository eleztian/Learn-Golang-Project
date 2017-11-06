package main

import (
	"fmt"
	"net"
)

type User struct {
	addr   string
	name   string
	passwd string
}

type Server struct {
	clients map[string]Client
	msgchan chan Msg
}

type Client struct {
	user *User
	conn net.Conn
}

type Msg struct {
	from *User
	to   *User
	msg  string
}

var server Server

func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + "  " + err.Error())
		return false
	}
	return true
}
