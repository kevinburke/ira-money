package main

import (
	"fmt"
	"testing"
)

func TestDeterminePrice(t *testing.T) {
	_, err := determinePrice(0, 0, 0, -1)
	if err == nil {
		t.Fatalf("expected determinePrice(-1) to return error, didn't")
	}
	if err.Error() != "invalid probability: -1" {
		t.Fatalf("expected determinePrice(-1) to return error, returned different one: %s", err)
	}

	_, err = determinePrice(0, 0, 0, 1.1)
	if err == nil {
		t.Fatalf("expected determinePrice(1.1) to return error, didn't")
	}
	if err.Error() != "invalid probability: 1.1" {
		t.Fatalf("expected determinePrice(1.1) to return error, returned different one: %s", err)
	}
}

func TestGetExponent(t *testing.T) {
	fmt.Println(getExponent(0.98))
	t.Fail()
}
