package object

import "testing"

func TestStrHashKey(t *testing.T) {
	hello1 := &Str{Value: "Hello World"}
	hello2 := &Str{Value: "Hello World"}
	
	diff1 := &Str{Value: "My name is tyson"}
	diff2 := &Str{Value: "My name is tyson"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}