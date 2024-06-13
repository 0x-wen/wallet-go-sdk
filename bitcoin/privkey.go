package bitcoin

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

// 1.非确定性钱包地址生成
func createPrivkey() {
	// secp256k1 生成私钥
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		panic(err)
	}

	privKeyHex := hex.EncodeToString(privKey.Serialize())
	fmt.Println("Private Key:", privKeyHex)

	// 创建一个新的WIF编码的私钥
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true) // true 表示压缩公钥
	if err != nil {
		log.Fatalf("Error creating WIF: %v", err)
	}
	// 打印WIF格式的私钥
	fmt.Println("WIF Private Key:", wif.String())

	// 从私钥获取公钥
	pubKey := privKey.PubKey()
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	// 使用公钥和主网络参数创建一个P2PKH地址
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Unable to create address: %v", err)
	}
	fmt.Println("P2PKH Address:", addr.EncodeAddress())
}

func HDWallet() {
	// 创建一个新的HD钱包种子
	// seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	// if err != nil {
	// 	log.Fatalf("Unable to generate seed: %v", err)
	// }

	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("助记词:", mnemonic)
	// 由助记词生成种子(Seed), password为空可兼容其他钱包
	seed := bip39.NewSeed(mnemonic, "")

	// 使用种子创建一个新的主私钥
	masterPrivKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Unable to create master private key: %v", err)
	}

	// Derive according to BIP44 path m/44'/0'/0'/0/0
	purposeIndex := hdkeychain.HardenedKeyStart + 44
	coinTypeIndex := hdkeychain.HardenedKeyStart + 0 // Bitcoin
	accountIndex := hdkeychain.HardenedKeyStart + 0  // Account 0
	changeIndex := 0                                 // External address (receiving)
	addressIndex := 0                                // First address

	// Correct usage of Derive method
	purposeKey, err := masterPrivKey.Derive(uint32(purposeIndex))
	if err != nil {
		log.Fatalf("Unable to derive purpose key: %v", err)
	}
	coinTypeKey, err := purposeKey.Derive(uint32(coinTypeIndex))
	if err != nil {
		log.Fatalf("Unable to derive coin type key: %v", err)
	}
	accountKey, err := coinTypeKey.Derive(uint32(accountIndex))
	if err != nil {
		log.Fatalf("Unable to derive account key: %v", err)
	}
	externalKey, err := accountKey.Derive(uint32(changeIndex))
	if err != nil {
		log.Fatalf("Unable to derive external chain key: %v", err)
	}
	childKey, err := externalKey.Derive(uint32(addressIndex))
	if err != nil {
		log.Fatalf("Unable to derive receive address key: %v", err)
	}
	// 获取子私钥
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		log.Fatalf("Unable to get private key: %v", err)
	}
	privKeyHex := hex.EncodeToString(privKey.Serialize())
	fmt.Println("Private Key:", privKeyHex)

	// 创建一个新的WIF编码的私钥
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true) // true 表示压缩公钥
	if err != nil {
		log.Fatalf("Error creating WIF: %v", err)
	}
	// 打印WIF格式的私钥
	fmt.Println("WIF Private Key:", wif.String())

	// 从私钥获取公钥
	pubKey := privKey.PubKey()

	// 使用公钥和主网络参数创建一个P2PKH地址
	addr, err := btcutil.NewAddressPubKeyHash(btcutil.Hash160(pubKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatalf("Unable to create address: %v", err)
	}
	fmt.Println("P2PKH Address:", addr.EncodeAddress())

	// TODO:
	// pay-to-script-hash (P2SH) btcutil.NewAddressScriptHash
	// pay-to-witness-pubkey-hash (P2WPKH) btcutil.NewAddressWitnessScriptHash
	// pay-to-taproot (P2TR)   btcutil.NewAddressTaproot

}
