package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/jonathon-chew/KVStore/pkg/commands"
	"github.com/jonathon-chew/KVStore/pkg/kvstore"
)

func handleConnection(conn net.Conn, kv *kvstore.HashTable) {
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		command, err := commands.ReadCommand(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Protocol error:", err)
			break
		}
		fmt.Printf("Command: /%s/ Args: /%v/\n", command.Name, command.Args)
		rv, parseErr := commands.ParseCommand(command.Name, command.Args, kv)
		if parseErr != nil {
			conn.Write([]byte("-" + parseErr.Error() + "\r\n"))
			continue
		}

		if rv != nil {
			conn.Write(fmt.Appendf(nil, "+%v\r\n", rv))
		} else {
			conn.Write([]byte("+Ok\r\n"))
		}

	}
}

func SendRequest(resp string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		return "", fmt.Errorf("[ERROR]: issue making TCP Connection %s", err.Error())
	}
	defer conn.Close()

	_, err = conn.Write([]byte(resp))
	if err != nil {
		return "", fmt.Errorf("[ERROR]: issue writing TCP Connection %s", err.Error())
	}

	reader := bufio.NewReader(conn)
	tcp_resp, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("[ERROR]: reading from TCP response %s", err.Error())
	}

	return tcp_resp, nil
}

func Server() {

	fmt.Println("Up and running...")

	kv := kvstore.NewKVStore(10)

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[ERROR]: %s\n", err.Error())
		}

		go handleConnection(conn, kv)
	}
}
