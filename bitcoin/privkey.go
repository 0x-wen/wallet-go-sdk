package bitcoin

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

var (
	netParams = &chaincfg.MainNetParams
)

// 1.非确定性钱包地址生成
func createPrivkey() {
	// secp256k1 生成私钥
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建WIF
	wif, err := btcutil.NewWIF(privKey, netParams, true) // true 表示压缩私钥
	if err != nil {
		fmt.Println("Error creating WIF:", err)
		return
	}

	// 输出WIF
	fmt.Println("WIF:", wif.String())

	// 获取公钥的序列化压缩形式
	serializedPubKey := privKey.PubKey().SerializeCompressed()

	// 计算公钥的哈希（RIPEMD160(SHA256(pubKey)))
	pubKeyHash := btcutil.Hash160(serializedPubKey)

	// 使用公钥哈希生成比特币地址
	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Bitcoin address:", address.String())
}

// func main() {
// 	//createMasterKey()

// 	masterKey, _ := btcec.NewPrivateKey()
// 	validatorPrivPubKeyAddr(masterKey.Serialize())
// }

// func createMasterKey() {
// 	mnemonic := "bright radio barely throw ring used wage barrel another cream treat essence"
// 	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 从种子生成根扩展私钥 masterKey是整个扩展私钥的表示，包含了链码、深度、父指纹、索引等信息。
// 	masterKey, err := hdkeychain.NewMaster(seed, netParams)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 检查是否为私钥扩展
// 	if masterKey.IsPrivate() {
// 		// privKey是从masterKey导出的ECDSA私钥,并打印ECDSA私钥的序列化形式
// 		privKey, err := masterKey.ECPrivKey()
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Printf("Master Private Key: %x\n", privKey.Serialize())

// 		// 展示获取公钥的两种形式,并判断是否一致
// 		pubKey, err := masterKey.ECPubKey()
// 		if err != nil {
// 			panic(err)
// 		}
// 		//[3 104 88 188 30 161 245 67 5 123 15 134 171 214 214 208 207 73 42 87 9 36 94 219 209 98 130 122 8 71 22 69 163]
// 		//[3 104 88 188 30 161 245 67 5 123 15 134 171 214 214 208 207 73 42 87 9 36 94 219 209 98 130 122 8 71 22 69 163]
// 		pubKey2 := privKey.PubKey()
// 		if !pubKey2.IsEqual(pubKey) {
// 			log.Fatal("Error:", pubKey.SerializeCompressed(), pubKey2.SerializeCompressed())
// 		}
// 		fmt.Printf("Master Public Key Compressed: %x\n", pubKey.SerializeCompressed())
// 		fmt.Printf("Master Public Key Uncompressed: %x\n", pubKey.SerializeUncompressed())

// 		// 生成地址
// 		// 1. 这是直接使用公钥作为支付地址，不常见。源码中的 NewAddressPubKey 结构体和相关函数用于处理这种类型的地址
// 		//address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), netParams)
// 		// 2.Pay-to-Pubkey-Hash (P2PKH) 地址,由公钥的哈希（通过RIPEMD160算法）生成
// 		hash := btcutil.Hash160(pubKey.SerializeCompressed())
// 		address, err := btcutil.NewAddressPubKeyHash(hash, netParams)
// 		if err != nil {
// 			fmt.Println("Error generating address:", err)
// 			return
// 		}
// 		fmt.Printf("P2PKH Address: %s\n", address.String())

// 	} else {
// 		log.Fatal("The master key is not a private key")
// 	}
// }

// func isValidPrivateKey(privKey []byte) bool {
// 	// 检查私钥长度是否正确
// 	if len(privKey) != 32 {
// 		return false
// 	}

// 	// 检查私钥是否在曲线的范围内
// 	privKeyInt := new(big.Int).SetBytes(privKey)
// 	curveOrder := btcec.S256().Params().N
// 	return privKeyInt.Cmp(curveOrder) >= 0 && privKeyInt.Cmp(big.NewInt(1)) >= 0
// }

// func validatorPrivPubKeyAddr(privKeyBytes []byte) {
// 	// 验证私钥
// 	if isValidPrivateKey(privKeyBytes) {
// 		fmt.Println("The private key is valid.")
// 	} else {
// 		fmt.Println("The private key is invalid.")
// 	}

// 	// 从私钥派生公钥
// 	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
// 	pubKey := privKey.PubKey()

// 	// 验证公钥
// 	if pubKey.IsOnCurve() {
// 		fmt.Println("The public key is on the curve.")
// 	} else {
// 		fmt.Println("The public key is not on the curve.")
// 	}

// 	// 生成地址
// 	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
// 	address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
// 	if err != nil {
// 		fmt.Println("Failed to create address:", err)
// 		return
// 	}
// 	fmt.Println("P2PKH Address:", address.String())

// 	// 验证WIF
// 	wif, err := btcutil.NewWIF(privKey, netParams, true)
// 	if err != nil {
// 		fmt.Println("Failed to create WIF:", err)
// 		return
// 	}
// 	fmt.Println("WIF:", wif.String())
// }
