package cosmos

import (
	"testing"
)

func Test_createPrivKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"钱包地址生成"},
		// privKey: aca02a78ff1425fea15eb3209ac40aef1a397632b345bca41d9fd93f0f01cef9
		// address: cosmos1zfrt0kc5988rv86la59vlhgxrhcqf4xpedm9em
		// keplr-wallet: cosmos1zfrt0kc5988rv86la59vlhgxrhcqf4xpedm9em
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPrivKey()
		})
	}
}

func TestNewHDWallet(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"生成确定性分层钱包"},
		// 助记词： tree bunker flower kid permit kick dismiss actor truth uncover message increase
		// privKey: b53b6ecb38e33e1c66f85e9d42f45296e1f6a19fe5703695acc39d7a556d7af2
		// accAddr: cosmos1p4p6s9km58d5h3c8w5a5h5vcgg53e3knduyg9t
		// keplr-wallet: cosmos1p4p6s9km58d5h3c8w5a5h5vcgg53e3knduyg9t
		// cosmos1p4p6s9km58d5h3c8w5a5h5vcgg53e3knduyg9t
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewHDWallet()
		})
	}
}
