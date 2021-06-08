package main

import (
	"fmt"
	"testing"
)

func TestByte(t *testing.T) {

	b := []byte("whatever")
	b, err := PassThrough(b)
	if err != nil {

	}
	fmt.Printf("zie of array %d\n", len(b))
	fmt.Printf("zie of array %s\n", string(b))
}

func PassThrough(b2 []byte) ([]byte, error) {
	fmt.Printf("zie of array %d\n", len(b2))

	b2 = []byte("new")
	return b2, nil

}
