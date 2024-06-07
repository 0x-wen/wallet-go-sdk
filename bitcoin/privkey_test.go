package bitcoin

import "testing"

func Test_createPrivkey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"生成公私钥对与地址"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPrivkey()
		})
	}
}
