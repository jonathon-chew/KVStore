package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/jonathon-chew/KVStore/pkg/commands"
	"github.com/jonathon-chew/KVStore/pkg/kvstore"
)

func handleConnection(conn net.Conn, kv kvstore.Store) {
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

func SendRequest(resp, port string) (string, error) {

	cleanPort, err := resolvePort(port)
	if err != nil {
		return "", err
	}

	conn, err := net.Dial("tcp", "localhost"+cleanPort)
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

func resolvePort(port string) (string, error) {
	var cleanPort string
	if port == "" || port == "*" {
		cleanPort = ":6379"
	}

	if !strings.Contains(port, ":") && cleanPort == "" {
		check_port := port
		_, err := strconv.Atoi(check_port)
		if err != nil {
			// fmt.Printf("[ERROR]: %s\n", err.Error())
			return "", fmt.Errorf("[ERROR]: %s\n", err.Error())
		}

		cleanPort = ":" + port
	}

	if strings.Contains(port, ":") && cleanPort == "" {
		if len(port) > 1 {
			check_port := port[1:]
			_, err := strconv.Atoi(check_port)
			if err != nil {
				// fmt.Printf("[ERROR]: %s\n", err.Error())
				return "", fmt.Errorf("[ERROR]: %s\n", err.Error())
			}
		}

		cleanPort = port
	}

	return cleanPort, nil
}

/*
This starts a KVStore TCP server on the port passed in
Port is expected to have ":" in it
Requires the port number to continue
if "" or * is passed in then the default port will be used ":6379"
*/
func Server(port string) int {

	fmt.Println("Up and running...")
	kv := kvstore.NewKVStore(10)

	cleanPort, err := resolvePort(port)
	if err != nil {
		fmt.Print(err.Error())
		return 1
	}

	listener, err := net.Listen("tcp", cleanPort)
	if err != nil {
		fmt.Printf("[ERROR]: %s\n", err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("[ERROR]: %s\n", err.Error())
			return 1
		}

		go handleConnection(conn, kv)
	}
}
