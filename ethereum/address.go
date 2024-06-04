package main

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func main() {
	s := bip39Mnemonic()
	bip32HD(s)
	pubKeyToAddress(s)
}

// bip39Mnemonic 生成一个BIP39兼容的助记词。
func bip39Mnemonic() (seed []byte) {
	// 生成256位的随机熵，用于创建BIP39助记词。即 128 个 16 进制字符
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	fmt.Println("助记词：", mnemonic)

	// 由助记词生成种子(Seed)
	seed = bip39.NewSeed(mnemonic, "salt")
	fmt.Println("New seed:", seed)
	return
}

// bip32HD 根据给定的种子字节生成BIP32 Hierarchical Deterministic钱包。
// 该函数实现了BIP32规范，用于创建一套可从单个种子衍生出无数个公私钥对的密钥体系。
// 参数:
//
//	seed []byte - 用于生成BIP32根密钥的随机种子字节。
func bip32HD(seed []byte) {
	// 由种子生成主账户扩展私钥(私钥和链码)
	masterKey, _ := bip32.NewMasterKey(seed) //从 0 开始），主账户密钥 masterKey 序号是 0，以此类推，这个就叫做索引号（32 位）
	fmt.Println("masterKey:", masterKey)

	childKey1, _ := masterKey.NewChildKey(1)

	// 用主账户公钥 派生 子账户公钥
	publicKey := masterKey.PublicKey()
	childPubKey1, _ := publicKey.NewChildKey(1)

	// 当前两种私钥派生得到的数据一致。
	// 第一种: 通过扩展公钥派生出子账户公钥,扩展私钥派生出子账户私钥,实现公私钥解耦
	// 第二种: 强化派生限制父公钥派生子公钥,严格通过 扩展私钥 -> 子私钥 -> 子公钥
	pubK1, _ := childPubKey1.Serialize()
	pubK2, _ := childKey1.PublicKey().Serialize()
	if !bytes.Equal(pubK1, pubK2) {
		fmt.Println("pubK1 != pubK2")
	}

	// 强化派生
	// 索引号在 0 和 2^31–1(0x0 to 0x7FFFFFFF)之间的只用于常规派生。
	// 索引号在 2^31 和 2^32– 1(0x80000000 to 0xFFFFFFFF)之间的只用于强化派生。
	childKeyPro, _ := masterKey.NewChildKey(bip32.FirstHardenedChild)
	pubKPro, _ := childKeyPro.PublicKey().Serialize()
	fmt.Println("pubKPro:", pubKPro)

	childPubKeyPro, _ := publicKey.NewChildKey(bip32.FirstHardenedChild)
	fmt.Println("childPubKeyPro:", childPubKeyPro) // childPubKeyPro: <nil>
}

func pubKeyToAddress(seed []byte) {
	// 创建子账户, 为了保护主私钥安全，所有主私钥派生的第一级账户，都采用强化派生。
	masterKey, _ := bip32.NewMasterKey(seed)
	// 以太坊的币种类型是60, 以路径（path: "m/44'/60'/0'/0/0"）为例
	key, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)  // 强化派生 对应 purpose'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(60)) // 强化派生 对应 coin_type'
	key, _ = key.NewChildKey(bip32.FirstHardenedChild + uint32(0))  // 强化派生 对应 account'
	key, _ = key.NewChildKey(uint32(0)) // 常规派生 对应 change
	key, _ = key.NewChildKey(uint32(0)) // 常规派生 对应 address_index

	fmt.Println("privKey:", key)
	fmt.Println("privKey2:", key.String())

	// 将key转换为十六进制字符串
	hexKey := fmt.Sprintf("%x", key.Key)
	fmt.Println("Hex Key:", hexKey)

	// 生成地址
	pubKey, _ := crypto.DecompressPubkey(key.PublicKey().Key)
	address := crypto.PubkeyToAddress(*pubKey).Hex()
	fmt.Println("address:", address)
}