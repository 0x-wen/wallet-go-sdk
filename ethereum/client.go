package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	formAddr   = "0xD0B7B6d2D38EA232C43355313E17B7cE81082cC6"
	privKeyHex = "3ac1f101223efc6c05002eec99a19654e19e2e0d2a698d6cc5d51922f477c9df"
	toAddr     = "0x75751bF3A86eA2F19660229C112dF7DAd84b8c02"
)

// client的代码去 https://dashboard.alchemy.com/ 点击APIKEY后复制
func newClient() (*ethclient.Client, error) {
	// client, err := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/k7J02LbbJiACCe52gTgZ64sY-sj-AZux")
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/k7J02LbbJiACCe52gTgZ64sY-sj-AZux")
	if err != nil {
		log.Fatal(err)
	}

	// Get the balance of an account
	account := common.HexToAddress(formAddr)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Account balance: %d\n", balance)

	// Get the latest known block
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Latest block: %d\n", block.Number().Uint64())
	return client, nil
}

func transaction() {
	privKey, _ := crypto.HexToECDSA(privKeyHex)
	formAddr := crypto.PubkeyToAddress(privKey.PublicKey)
	fmt.Println("formAddr:", formAddr.Hex())
	cli, _ := newClient()

	// 构造交易数据
	nonce, err := cli.NonceAt(context.Background(), formAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nonce:", nonce)
	toAddress := common.HexToAddress(toAddr)
	amount := big.NewInt(1000000000000000)
	gasLimit := uint64(21000)            // 标准交易的 gas 限制
	gasPrice := big.NewInt(200000000000) // gas 价格，你需要根据网络情况调整
	// gasPrice, err := cli.SuggestGasTipCap(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("gasPrice:", gasPrice)
	tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	chainID, err := cli.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// 对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	if err != nil {
		log.Fatal(err)
	}
	// 将签名后的数据通过clp编码为字节数组,前端eth_sendRawTransaction接口中params参数
	encodedTx, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatal(err)
	}
	params := hex.EncodeToString(encodedTx)
	fmt.Printf("eth_sendRawTransaction Params: %s\n", params)
	// 广播交易到测试节点
	if err = cli.SendTransaction(context.Background(), signedTx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction Hash: 0x%x\n", signedTx.Hash())
}

// tx1: 0x20369b6783aca2f5b138a34222be0c29da677f5a7aac67bb136eea0961dfdfdd
// 这笔交易由于gas比较低的原因,一直没被打包进区块中,然后使用第二笔交易成功之后，这个已经废弃
// tx2: 0x92018f4b8fe03774b3714f48507bc92739049f0911dd8691a2e5c460bf580837
// tx3: 0x1f191ed79008aeaa1722abb067f2775feba4c3bd4e541f8ca7d55d471aeb94b3
// - Pending  This txn hash was found in our secondary node and should be picked up by our indexer in a short while.
// - Indexing  This transaction has been included and will be reflected in a short while.
// - Success
