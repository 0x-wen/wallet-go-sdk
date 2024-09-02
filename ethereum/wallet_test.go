package main

import (
	"fmt"
	"testing"
)

func Test_bip39Mnemonic(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"testSeed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSeed := bip39Mnemonic()
			fmt.Println(len(gotSeed))
		})
	}
}
