package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
)

var (
	cil, _ = newClient()
)

// 主要练习常用接口api使用方式
// https://docs.alchemy.com/reference/eth-gettransactionbyhash
// 1.获取最新块高
// 2.根据块高获取块里面的信息
// 3.根据交易 Hash 获取交易详情
// 4.获取交易状态

func openAPI() {
	// 1.获取最新块高
	blockNumber, err := cil.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("blockNumber:", blockNumber)

	// 2.根据块高获取块里面的信息
	/*
	type Transaction struct {
		inner TxData    // 它包含了交易的共识内容:发送者地址、接收者地址、交易金额、交易费用、数据字段等。
		time  time.Time // 该交易首次在本地节点被观察到的时间

		// caches
		hash atomic.Pointer[common.Hash] // 存储的是交易的哈希值
		size atomic.Uint64        // 存储了交易的RLP（Recursive Length Prefix）编码后的大小
		from atomic.Pointer[sigCache]  // 缓存交易签名的解析结果，以便在需要验证交易签名时可以快速获取
	}
	*/
	block, err := cil.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Fatal(err)
	}
	txs := block.Transactions()
	fmt.Println(txs)
}
