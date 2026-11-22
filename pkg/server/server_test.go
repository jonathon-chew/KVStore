package server

import (
	"testing"

	"github.com/jonathon-chew/KVStore/pkg/resp"
)

func TestResolvePort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty", "", ":6379"},
		{"Default", ":6379", ":6379"},
		{"DefaultNoColon", "6379", ":6379"},
		{"Custom", "4567", ":4567"},
		{"CustomColon", ":4568", ":4568"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolvePort(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestServerDefault(t *testing.T) {

	tests := map[string]string{
		"ASTERIX":              "*",
		"NON_DEFAULT_NO_COLON": "4567",
		"NON_DEFAULT_COLON":    ":4568",
	}

	for test_name, test_port := range tests {
		go Server(test_port)

		reques, err := resp.MakeResp("SET", "HELLO", "WORLD")
		if err != nil {
			t.Fatalf("[ERROR]: %s Unable to make the request string to send to the server", test_name)
		}

		_, err = SendRequest(reques, test_port)
		if err != nil {
			t.Fatalf("[ERROR]: %s Unable to send set comamnd, %s", test_name, err.Error())
		}

		check, err := resp.MakeResp("GET", "HELLO", "")
		if err != nil {
			t.Fatalf("[ERROR]: %s Unable to make the request string to send to the server", test_name)
		}

		confirm_respose, err := SendRequest(check, test_port)
		if err != nil {
			t.Fatalf("[ERROR]: %s Unable to send get comamnd, %s", test_name, err.Error())
		}

		if confirm_respose != "+WORLD\r\n" {
			t.Fatalf("[ERROR]: %s Expected /%s/ Got /%s/", "+WORLD\n", test_name, confirm_respose)
		}
	}
}
