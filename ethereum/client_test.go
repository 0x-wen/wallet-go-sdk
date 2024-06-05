package main

import (
	"testing"
)

func Test_client(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"client"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newClient()
		})
	}
}

func Test_transaction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"测试发送交易至eth-sepolia"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction()
		})
	}
}
