package commands

import (
	"bufio"
	"fmt"
	"io"
	"testing"

	"github.com/jonathon-chew/KVStore/pkg/kvstore"
)

func Test_ReadCommand(t *testing.T) {

	var ioR io.Reader
	ioR.Read([]byte("*3"))
	buf := bufio.NewReader(ioR)

	for {
		command, err := ReadCommand(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Protocol error:", err)
			break
		}

		t.Log(command)
	}
}

func Test_ParseCommand(t *testing.T) {

	// SET
	kv_set := kvstore.NewKVStore(1)
	set_response, err := ParseCommand("SET", []string{"Set_Test"}, kv_set)
	if err != nil {
		t.Error(err)
	}

	if _, err := kv_set.Get("Set_Test"); !err {
		t.Errorf("[ERROR]: Unable to get the key that should have been set")
	}

	t.Log(set_response)

	// GET
	kv_get := kvstore.NewKVStore(1)
	kv_get.Add("Get_Test", "True")
	get_response, err := ParseCommand("GET", []string{"Get_Test"}, kv_get)
	if err != nil {
		t.Error(err)
	}

	t.Log(get_response)

	// INCR
	kv_inc := kvstore.NewKVStore(1)
	inc_response, err := ParseCommand("INCR", []string{"counter"}, kv_inc)
	if err != nil {
		t.Error(err)
	}

	_, ok := kv_inc.Get("counter")
	if !ok {
		t.Error("[ERROR]: counter should have been initialised but wasn't")
	}

	t.Log(inc_response)

	// INCRBY
	kv_incby := kvstore.NewKVStore(1)
	// Init to 1
	_, incby_err := ParseCommand("INCR", []string{"counter"}, kv_incby)
	if incby_err != nil {
		t.Error(incby_err)
	}

	// Incriment by more than 1
	_, incby_err_2 := ParseCommand("INCRBY", []string{"counter", "4"}, kv_incby)
	if incby_err_2 != nil {
		t.Error(incby_err_2)
	}

	value, ok := kv_inc.Get("counter")
	if !ok {
		t.Error("[ERROR]: counter should have been initialised but wasn't")
	}

	if value != 5 {
		t.Errorf("[ERROR]: Counter should be 5, initalised to 1 and incrimented by 4, it was actually %v", value)
	}

	t.Log(inc_response)

}
