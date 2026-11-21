package perf

import "testing"

func TestPerf(t *testing.T) {
	kv := NewPerf()

	Simple("", "SET", "", kv)

	value, exists := kv.Get("SET")
	if !exists {
		t.Fatalf("[ERROR]: /%s/ could not be found", "SET")
	}

	if value != 1 {
		t.Errorf("[ERROR]: value came back as /%s/ not one", value)
	}
}

func TestPerf1000sTimes(t *testing.T) {
	kv := NewPerf()

	_ = "INCRBY"
	key_name := "SET_KEY"
	impliment_times := 1000

	for range impliment_times {
		resp, err := Simple("", key_name, "", kv)
		if err != nil {
			t.Error(err.Error())
		}

		if resp != nil {
			t.Log(resp)
		}
	}

	value, exists := kv.Get(key_name)
	if !exists {
		t.Fatalf("[ERROR]: /%s/ could not be found", key_name)
	}

	if value != impliment_times {
		t.Errorf("[ERROR]: value came back as /%d/ not %d", value, impliment_times)
	}

	t.Log("Buckets: ", len(kv.Buckets))
	t.Log("Entries: ", kv.Size)
}
