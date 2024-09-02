package bitcoin

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestCreateAddr(t *testing.T) {
	type args struct {
		network *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"testNewBtcAddr", args{MainNetParams}, false},
		/*  可自行导入钱包进行验证
		&{7c688e12f67b8bf33e07db1e9ee421d3f35d1f47c2eb67dd60bed24fe722094b
		L1PYYFSJ3WksMJZuCto9xu3zxMsNdyxaqCkMCtUDnUfixFJLMVf9
		03c7db7c49b0d57a9bd7c608995d4e85c3ba19b2d6d5ab2ca425f034c95fd08a7c
		1JMwkFL8LiVYJnqSdA5zUoTh44tV2Xv1FN
		3K7Au7SZKG1UDBxZrv2J6XpTgsxRhZUccd
		bc1qheczs27g3d9jm88txyq7rws6805dc07y7ea7kg
		bc1p8tfvst3ad265tje8fqflvm2sygz05vmhwcg67n56jx9kaxr5ejpsn0df6x}
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBtcAddr(tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)

		})
	}
}

func TestGeneartePrivKey(t *testing.T) {
	type args struct {
		network *chaincfg.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"testGeneartePrivKey", args{MainNetParams}, false},
		/*
			b80bf961c414e6373846cb670fe40dbd61d115deeeae8192467f68eff23fb64f
			L3PURysgS3sbByk2sGBwyn8nm5Rcped8fC667qc8Yzkuy3vhRPNt
			038623e98ea08e4713e68feeb405940cb02c67e6c7b848d34bd1389202c04416fc
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := GeneartePrivKey(tt.args.network)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneartePrivKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got, got1, got2)
		})
	}
}

func TestNewHDWallet(t *testing.T) {
	type args struct {
		network   *chaincfg.Params
		mnemonic  string
		addrIndex uint32
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"testNewHDWallet", args{MainNetParams, "", 0}, false},
		{"testNewHDWallet", args{MainNetParams, "share mirror defy grief flower mosquito speak noise since trim mix behave", 0}, false},
		{"testNewHDWallet", args{MainNetParams, "share mirror defy grief flower mosquito speak noise since trim mix behave", 1}, false},
		{"testNewHDWallet", args{MainNetParams, "share mirror defy grief flower mosquito speak noise since trim mix behave", 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHDWallet(tt.args.network, tt.args.mnemonic, tt.args.addrIndex)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHDWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
