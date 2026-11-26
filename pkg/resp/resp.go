package resp

import (
	"fmt"
	"strings"
)

// ioR := strings.NewReader("*3\r\n$3\r\nSET\r\n$8\r\ngreeting\r\n$5\r\nHello\r\n")

// MakeResp builds a RESP command containing command, key, and optionally value.
func MakeResp(command, key, value string) (string, error) {

	var sb strings.Builder

	// Start message
	if value != "" {
		_, err := sb.Write([]byte("*3\r\n"))
		if err != nil {
			fmt.Println("[ERROR]: ", err.Error())
			return "", err
		}
	} else {
		_, err := sb.Write([]byte("*2\r\n"))
		if err != nil {
			fmt.Println("[ERROR]: ", err.Error())
			return "", err
		}
	}

	// Write the command get SET / GET / INCRBY
	_, err := sb.Write(fmt.Append(nil, fmt.Sprintf("$%d\r\n%s\r\n", len(command), command)))
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
		return "", err
	}

	// Write the key
	_, err = sb.Write(fmt.Append(nil, fmt.Sprintf("$%d\r\n%s\r\n", len(key), key)))
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
		return "", err
	}

	// Write the value IF relevent
	if value != "" {
		_, err = sb.Write(fmt.Append(nil, fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
		if err != nil {
			fmt.Println("[ERROR]: ", err.Error())
			return "", err
		}
	}

	return sb.String(), nil
}
