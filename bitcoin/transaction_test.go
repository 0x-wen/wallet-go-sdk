package bitcoin

import "testing"

func Test_transaction(t *testing.T) {
	tests := []struct {
		name string
	}{
		// txid: https://mempool.space/testnet/tx/e15c60c5206f19995c5ce1eb6bad4fb5aa7929b25aa2de017be2989861d7cb4e
		{"test-RawTx"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction()
		})
	}
}
