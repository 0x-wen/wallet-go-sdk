package cosmos

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/go-bip39"
)

func createPrivKey() {
	privKey := secp256k1.GenPrivKey()
	fmt.Println("privKey:", hex.EncodeToString(privKey.Bytes()))

	var privKeyBytes sdk.AccAddress = privKey.PubKey().Address().Bytes()
	baseAccount := authtypes.NewBaseAccount(privKeyBytes, privKey.PubKey(), 0, 0)

	fmt.Println("address:", baseAccount.Address)
}

func NewHDWallet() {
	// bip39  -> 助记词
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("助记词：", mnemonic)
	// 由助记词生成种子(Seed), password为空可兼容其他钱包
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, ch := hd.ComputeMastersFromSeed(seed)
	priv, _ := hd.DerivePrivateKeyForPath(masterKey, ch, "m/44'/118'/0'/0/0")

	privKey := &secp256k1.PrivKey{Key: priv}
	fmt.Println("privKey:", hex.EncodeToString(privKey.Bytes()))
	accAddr, err := bech32.ConvertAndEncode(sdk.Bech32MainPrefix, privKey.PubKey().Address())
	if err != nil {
		panic(err)
	}
	fmt.Println("accAddr:", accAddr)
}

func Import(privKeyHex string) (*secp256k1.PrivKey, error) {
	priBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil, err
	}
	return &secp256k1.PrivKey{Key: priBytes}, nil
}