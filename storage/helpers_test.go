package storage

import "testing"

func TestDataToBytesAndBackIsEqual(t *testing.T) {
	type Person struct {
		Name string
		age  int
	}

	da := Person{"Douglas Adams", 49}
	b, err := dataToBytes(da)
	if err != nil {
		t.Error(err)
	}
	da2 := Person{}
	err = bytesToData(*da2, b)
	if err != nil {
		t.Error(err)
	}

	if da != da2 {
		t.Fail()
	}
}
