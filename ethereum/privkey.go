package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// 非确定性钱包:地址独立,无关联性,但需要备份所有地址的私钥

// createPrivKeyAndAddress 生成加密货币私钥和地址。
// 该函数首先生成一个随机私钥，然后从私钥派生出公钥和地址。
func GeneratePrivKey() {
	// 生成随机私钥
	privKey, _ := crypto.GenerateKey()
	// 将私钥转换为字节序列
	privKeyBytes := privKey.D.Bytes()
	// 将私钥字节序列转换为十六进制字符串
	privKeyHex := hex.EncodeToString(privKeyBytes)
	// 打印私钥的十六进制表示
	fmt.Println(privKeyHex)

	// 提取私钥对应的公钥
	PublicKey := privKey.PublicKey
	// 将公钥转换为以太坊地址
	addr := crypto.PubkeyToAddress(PublicKey)
	fmt.Println(addr)
}

func importPrivKey(privKeyHex string) {
	// 检查私钥字符串长度
	if len(privKeyHex) != 64 {
		fmt.Println("Invalid private key length")
		return
	}
	// 将十六进制字符串解码为字节形式的私钥
	// privKeyBytes, _ := hex.DecodeString(privKeyHex)
	privKey, _ := crypto.HexToECDSA(privKeyHex)

	// 提取私钥对应的公钥
	publicKey := privKey.PublicKey
	fmt.Println("Public Key X:", publicKey.X)
	fmt.Println("Public Key Y:", publicKey.Y)

	// 将公钥转换为以太坊地址
	addr := crypto.PubkeyToAddress(publicKey)
	fmt.Println(addr)
}

// createPrivKey 生成一个椭圆曲线私钥，并据此计算公钥和以太坊地址。
func createPrivKey() {
	// 使用secp256k1曲线生成私钥
	curve := secp256k1.S256()
	// 创建一个字节切片，长度足够存储曲线参数N的比特长度
	b := make([]byte, curve.Params().N.BitLen()/8)
	// 从随机源读取数据填充私钥字节切片
	io.ReadFull(rand.Reader, b)
	// 将字节切片转换为大整数形式的私钥
	key := new(big.Int).SetBytes(b)
	// 输出私钥的字节长度和十六进制表示
	fmt.Println("key:", len(key.Bytes()))
	fmt.Println("key:", hex.EncodeToString(key.Bytes()))

	// 使用私钥计算公钥，X和Y是公钥的坐标
	X, Y := curve.ScalarBaseMult(key.Bytes())
	// 将公钥坐标转换为压缩形式的公钥切片
	pubKey := elliptic.Marshal(curve, X, Y)
	// 输出公钥
	fmt.Println("pubKey:", pubKey)

	// 对公钥进行Keccak256哈希，得到压缩公钥
	compressPubKey := crypto.Keccak256(pubKey[1:])
	// 从压缩公钥中提取以太坊地址
	addr := common.BytesToAddress(compressPubKey[12:])
	// 输出以太坊地址
	fmt.Println("addr:", addr.String())
}

func createPrivKey2() {
	curve := secp256k1.S256()
	b := make([]byte, curve.Params().N.BitLen()/8)
	io.ReadFull(rand.Reader, b)
	key := new(big.Int).SetBytes(b)
	fmt.Println("key:", len(key.Bytes()))
	fmt.Println("key:", hex.EncodeToString(key.Bytes()))
	// 使用私钥计算公钥，X和Y是公钥的坐标
	x, y := curve.ScalarBaseMult(key.Bytes())

	// 获取公钥的字节表示形式
	publicKey := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	pubKeyBytes := crypto.FromECDSAPub(&publicKey)
	// 输出公钥
	fmt.Println("pubKey:", hex.EncodeToString(pubKeyBytes))

	compressPubKey := crypto.Keccak256(pubKeyBytes[1:])
	addr := common.BytesToAddress(compressPubKey[12:])

	// 输出以太坊地址
	fmt.Println("addr:", addr.String())
}
