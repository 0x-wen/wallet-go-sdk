package main

import (
	"testing"
)

func Test_privKeyPubKeyValidator(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"公私钥验证-1"},
		// 196a3f6dfbb2b057d2b204777d3228e26f5c57fda648c371675a494b29001136
		// 0x0709c9A4070E09d4e6467AB56841Ab818A9CdB3A
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPrivKeyAndAddress()
		})
	}
}

func Test_importPrivKey(t *testing.T) {
	type args struct {
		privKeyHex string
	}
	tests := []struct {
		name string
		args args
	}{
		{"根据十六进制私钥生成地址", args{"196a3f6dfbb2b057d2b204777d3228e26f5c57fda648c371675a494b29001136"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			importPrivKey(tt.args.privKeyHex)
		})
	}
}

func Test_createPrivKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"手动生成私钥"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPrivKey()
		})
	}
}

func Test_createPrivKey2(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"解决生成私钥,elliptic.Marshal弃用问题"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPrivKey2()
		})
	}
}
