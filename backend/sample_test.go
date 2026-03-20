package main

import "testing"

func TestSample(t *testing.T) {
	if 1+1 != 2 {
		t.Errorf("1+1 should be 2")
	}
}