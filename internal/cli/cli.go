package cli

import (
	"fmt"

	"github.com/jonathon-chew/KVStore/pkg/resp"
	"github.com/jonathon-chew/KVStore/pkg/server"
)

type CommandBuilder struct {
	command string
	key     string
	value   string
	port    string
}

func CLI(arguments []string) {

	var cb CommandBuilder

	for idx := 0; idx <= len(arguments); idx++ {
		command := arguments[idx]
		switch command {
		case "--command":
			if len(arguments) > idx+1 {
				cb.command = arguments[idx+1]
			}

			idx++
		case "--key":
			if len(arguments) > idx+1 {
				cb.key = arguments[idx+1]
			}

			idx++
		case "--value":
			if len(arguments) > idx+1 {
				cb.value = arguments[idx+1]
			}

			idx++
		case "--port":
			if len(arguments) > idx+1 {
				cb.port = arguments[idx+1]
			}

			idx++
		default:
			fmt.Println("Unable to recognise command: ", command)
		}
	}

	if cb.command == "" || cb.key == "" {
		fmt.Printf("[ERROR]: Missing elements, minimum is a key and a command %v\n", cb)
		return
	}

	resp, err := resp.MakeResp(cb.command, cb.key, cb.value)
	if err != nil {
		fmt.Println("[ERROR]: unable to make resquest ", err.Error())
	}

	tcp_resp, err := server.SendRequest(resp, cb.port)
	if err != nil {
		fmt.Printf(err.Error(), "\n")
		return
	}

	fmt.Println(tcp_resp)
}
