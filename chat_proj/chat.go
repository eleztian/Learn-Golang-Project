package main

import (
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strconv"
)

func main() {
	info := `
	chatit server [port] 
	    eg: chatit server 9090
	chatit client [Server_Addr]:[Server_Port] 
		eg: chatit client 192.168.0.74:9090 
	chatit client [Server_Addr]:[Server_Port] [count]      
		eg: chatit client 192.168.0.74:9090 500
	`
	if len(os.Args) < 2 {
		fmt.Println("Wrong parameter, usage:")
		fmt.Print(info)
		os.Exit(0)
	}

	if os.Args[1] == "server" && len(os.Args) == 3 {
		StartServer(os.Args[2])
	}

	if os.Args[1] == "client" && len(os.Args) == 3 {
		fmt.Println("start client to ", os.Args[2])
		tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
		}
		client := StartClient(tcpAddr)
		defer client.conn.Close()
	}

	if os.Args[1] == "client" && len(os.Args) == 4 {

		m, _ := strconv.Atoi(os.Args[3])
		tcpAddr, err := net.ResolveTCPAddr("tcp", os.Args[2])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(2)
		}
		debug.SetMaxThreads(30050)
		for i := 0; i < m; i++ {
			fmt.Println("start client", i, os.Args[2])
			StartClient(tcpAddr)
		}
	}
	select {}
	fmt.Println("end.")

}
