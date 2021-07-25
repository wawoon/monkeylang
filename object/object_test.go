package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "hello"}
	hello2 := &String{Value: "hello"}
	diff1 := &String{Value: "a"}
	diff2 := &String{Value: "a"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("the hash value is not same %v and %v", hello1.HashKey(), hello2.HashKey())
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("the hash value is not same %v and %v", diff1.HashKey(), diff2.HashKey())
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("the hash value is same %v and %v", hello1.HashKey(), diff1.HashKey())
	}
}
