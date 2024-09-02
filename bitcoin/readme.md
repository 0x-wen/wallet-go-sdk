## Bitcoin钱包开发实践

参考相关文档后，大致分解为三部分:

文档类: 官方文档查阅，API文档整理，测试网水龙头获取代币

离线地址生成: 
- 非确定性钱包: 每个地址独立,无关联性,但需要备份所有地址的私钥
- 确定性分层钱包: 遵循BIP32、BIP39、BIP43和BIP44等规范, 支持多币种和多账户

签名与广播: 使用私钥签名交易并广播上链


## 操作笔记

本练习在测试网完成，需要获取相关测试币。

### 环境准备

- 插件钱包：[unisat-wallet](https://chromewebstore.google.com/detail/unisat-wallet/ppbibelpcjmhbdihakflkdcoccbgbkpo?hl=zh-CN&utm_source=ext_sidebar)
- [测试网领水](https://bitcoinfaucet.uo1.net/send.php)



