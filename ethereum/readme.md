## ethereum钱包实践

### 目录结构讲解[./ethereum]
- privkey.go
  
  > 主要实现 非确定性钱包: 地址独立,无关联性,但需要备份所有地址的私钥
  
  生成公私钥对，公钥生成地址


- address.go

  > bip39、bip32、bip44 结合使用，HD wallet

  bip39 -> 生成一个BIP39兼容的助记词

  bip32 -> 从单个种子衍生出无数个公私钥对的密钥体系

  bip44 -> 解决了多币种、多账户钱包管理的问题, pubKeyToAddress 实现了一个符合 BIP-44 标准的路径以太坊钱包


- client.go

  > 组装交易和签名并发送

  transaction 实现了组装交易与签名和发送，并提供通过http发送的交易参数


- api.go

  > 主要通过http请求获取相关信息

  - 发送交易：params参数获取可参考transaction方法中 ` params := hex.EncodeToString(encodedTx)`

    ```shell
    curl --request POST \
         --url https://eth-sepolia.g.alchemy.com/v2/k7J02LbbJiACCe52gTgZ64sY-sj-AZux \
         --header 'accept: application/json' \
         --header 'content-type: application/json' \
         --data '
    {
      "id": 1,
      "jsonrpc": "2.0",
      "params": [
    "0xf86f04852e90edd0008252089475751bf3a86ea2f19660229c112df7dad84b8c0287038d7ea4c68000808401546d72a0d7618c5a377ae77c40c0e092c05545bbb24aacfff02569d74a1ab3bf27620d0ba02b5e3fe241f53d44a00691084d9e0d08475ce136b92da53b8b0b9b4fed285167"
      ],
      "method": "eth_sendRawTransaction"
    }'
    ```

  - 根据txhash查询交易信息

    ```shell
    curl --request POST \
         --url https://eth-sepolia.g.alchemy.com/v2/k7J02LbbJiACCe52gTgZ64sY-sj-AZux \
    --header 'Content-Type: application/json' \
    --header 'Cookie: _cfuvid=wq9cOWHZrQJC_tuTXZLmlWkdI8GPYot1N2.FsyNsrbw-1716723624624-0.0.1.1-604800000' \
    --data  '
    {
      "id": 1,
      "jsonrpc": "2.0",
      "method": "eth_getTransactionByHash",
      "params": [
        "0x18fb59de0ac1ab62401e08f6260c0fd838a35f113da467671c894c65d1b1ccc6"
      ]
    }'
    ```
