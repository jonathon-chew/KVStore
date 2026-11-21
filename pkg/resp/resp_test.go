package resp

import (
	"testing"
)

func TestResp(t *testing.T) {

	request, err := MakeResp("SET", "greeting", "Hello")
	if err != nil {
		t.Error(err)
	}

	var exp_string string = "*3\r\n$3\r\nSET\r\n$8\r\ngreeting\r\n$5\r\nHello\r\n"
	if request != exp_string {
		t.Errorf("%s\n This does not match the expected string\n%s\n", request, exp_string)
	}
}
