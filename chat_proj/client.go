package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func (client *Client) Input(inputchan chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error readstring", err.Error())
			continue
		}
		if input == "/quit\n" {
			fmt.Println("Byebye...")
			client.conn.Close()
			os.Exit(0)
		}
		inputchan <- input
	}
}

func (client *Client) Send(inputchan chan string) {
	for {
		input := <-inputchan
		_, err := client.conn.Write([]byte(input))
		if err != nil {
			fmt.Println(err.Error())
			client.conn.Close()
			os.Exit(2)
		}
	}
}

func (client *Client) Recv(outputchan chan string) {
	buf := make([]byte, 1024)
	for {
		length, err := client.conn.Read(buf)
		if checkError(err, "connection") == false {
			fmt.Println("Server is dead... Byebye")
			os.Exit(0)
		}
		outputchan <- string(buf[0:length])
	}
}

func (client *Client) Output(outputchan chan string) {
	for {
		outstr := <-outputchan
		fmt.Println(outstr)
	}
}

func StartClient(tcpAddr *net.TCPAddr) (cli *Client) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")
	var client Client
	var user User
	user.addr = user.addr
	user.name = user.addr
	client.user = &user
	client.conn = conn
	inputchan := make(chan string)
	outputchan := make(chan string)
	fmt.Println("Input you name")
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	go client.Send(inputchan)
	go client.Recv(outputchan)
	go client.Output(outputchan)
	go client.Input(inputchan)
	if err == nil {
		user.name = name
		cmd := fmt.Sprintf("/username %s", name)
		fmt.Println(cmd)
		inputchan <- cmd
	}
	return &client
}
