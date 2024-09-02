package bitcoin

import (
	"encoding/hex"
	"errors"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/tyler-smith/go-bip39"
)

const (
	LEGACY        = "legacy"        // p2pkh   bip44
	SEGWIT_NESTED = "segwit_nested" // p2sh    bip49
	SEGWIT_NATIVE = "segwit_native" // p2wpkh  bip84
	TAPROOT       = "taproot"       // p2tr    bip86
)

var (
	TestNet3Params = &chaincfg.TestNet3Params
	MainNetParams  = &chaincfg.MainNetParams // 主网
)

type BtcAddr struct {
	PrivKey    string
	WIFPrivKey string
	PubKey     string
	P2PKH      string
	P2SH       string
	P2WPKH     string
	P2TR       string
}

func NewBtcAddr(network *chaincfg.Params) (*BtcAddr, error) {
	privKey, wifPk, pubKey, err := GeneartePrivKey(network)
	if err != nil {
		return nil, err
	}

	publicKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	addrMap, err := ConvertPubKeyToAddresses(publicKey, network)
	if err != nil {
		return nil, err
	}

	return &BtcAddr{
		PrivKey:    privKey,
		WIFPrivKey: wifPk,
		PubKey:     pubKey,
		P2PKH:      addrMap["P2PKH"],
		P2SH:       addrMap["P2SH"],
		P2WPKH:     addrMap["P2WPKH"],
		P2TR:       addrMap["P2TR"],
	}, nil
}

func ConvertPubKeyToAddresses(publicKey []byte, network *chaincfg.Params) (map[string]string, error) {
	// 定义地址类型和对应的标识
	addressTypes := map[string]string{
		"P2PKH":  LEGACY,
		"P2SH":   SEGWIT_NESTED,
		"P2WPKH": SEGWIT_NATIVE,
		"P2TR":   TAPROOT,
	}

	addresses := make(map[string]string)

	for addrType, addrTypeValue := range addressTypes {
		addr, err := PubKeyToAddr(publicKey, addrTypeValue, network)
		if err != nil {
			return nil, err
		}
		addresses[addrType] = addr
	}

	return addresses, nil
}

func PubKeyToAddr(publicKey []byte, addrType string, network *chaincfg.Params) (string, error) {
	if network == nil {
		network = &chaincfg.MainNetParams
	}
	if addrType == LEGACY {
		p2pkh, err := btcutil.NewAddressPubKey(publicKey, network)
		if err != nil {
			return "", err
		}

		return p2pkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NATIVE {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}

		return p2wpkh.EncodeAddress(), nil
	} else if addrType == SEGWIT_NESTED {
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey), network)
		if err != nil {
			return "", err
		}
		redeemScript, err := txscript.PayToAddrScript(p2wpkh)
		if err != nil {
			return "", err
		}
		p2sh, err := btcutil.NewAddressScriptHash(redeemScript, network)
		if err != nil {
			return "", err
		}

		return p2sh.EncodeAddress(), nil
	} else if addrType == TAPROOT {
		internalKey, err := btcec.ParsePubKey(publicKey)
		if err != nil {
			return "", err
		}
		p2tr, err := btcutil.NewAddressTaproot(txscript.ComputeTaprootKeyNoScript(internalKey).SerializeCompressed()[1:], network)
		if err != nil {
			return "", err
		}

		return p2tr.EncodeAddress(), nil
	} else {
		return "", errors.New("address type not supported")
	}
}

func GeneartePrivKey(network *chaincfg.Params) (string, string, string, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
		return "", "", "", err
	}
	privKeyHex := hex.EncodeToString(privKey.Serialize())

	wif, err := btcutil.NewWIF(privKey, network, true) // true 表示压缩公钥
	if err != nil {
		log.Fatalf("Error creating WIF: %v", err)
	}

	pubKey := privKey.PubKey()
	pubKeyHex := hex.EncodeToString(pubKey.SerializeCompressed())

	return privKeyHex, wif.String(), pubKeyHex, nil
}

type HDWallet struct {
	mnemonic  string
	masterKey string
	childAddr BtcAddr
}

func NewHDWallet(network *chaincfg.Params, mnemonic string, addrIndex uint32) (*HDWallet, error) {
	var seed []byte
	if mnemonic != "" {
		seed = bip39.NewSeed(mnemonic, "")
	} else {
		mnemonic, seed, _ = GenerateMnemonic(network)
	}

	masterPrivKey, err := hdkeychain.NewMaster(seed, network)
	if err != nil {
		log.Fatalf("Unable to create master private key: %v", err)
	}

	// Derive according to BIP44 path m/44'/0'/0'/0/0
	purposeIndex := hdkeychain.HardenedKeyStart + 44
	coinTypeIndex := hdkeychain.HardenedKeyStart + 0 // Bitcoin
	accountIndex := hdkeychain.HardenedKeyStart + 0  // Account 0
	changeIndex := 0                                 // External address (receiving)
	addressIndex := addrIndex                        // First address

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

	// 创建一个新的WIF编码的私钥
	wif, err := btcutil.NewWIF(privKey, network, true) // true 表示压缩公钥
	if err != nil {
		log.Fatalf("Error creating WIF: %v", err)
	}

	// 从私钥获取公钥
	pubKey := privKey.PubKey().SerializeCompressed()

	addrMap, err := ConvertPubKeyToAddresses(pubKey, network)
	if err != nil {
		return nil, err
	}

	addr1 := BtcAddr{
		PrivKey:    hex.EncodeToString(privKey.Serialize()),
		WIFPrivKey: wif.String(),
		PubKey:     hex.EncodeToString(pubKey),
		P2PKH:      addrMap["P2PKH"],
		P2SH:       addrMap["P2SH"],
		P2WPKH:     addrMap["P2WPKH"],
		P2TR:       addrMap["P2TR"],
	}

	return &HDWallet{
		mnemonic:  mnemonic,
		masterKey: masterPrivKey.String(),
		childAddr: addr1,
	}, nil
}

func GenerateMnemonic(network *chaincfg.Params) (string, []byte, error) {
	// 方式一: 创建一个新的HD钱包种子
	// seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	// if err != nil {
	// 	log.Fatalf("Unable to generate seed: %v", err)
	// }

	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	// 由助记词生成种子(Seed), password为空可兼容其他钱包
	seed := bip39.NewSeed(mnemonic, "")
	return mnemonic, seed, nil
}
