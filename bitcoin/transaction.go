package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func GetUnspent() (*wire.OutPoint, *txscript.MultiPrevOutFetcher) {
	// 交易的哈希值,并指定使用UTXO的索引
	txHash, _ := chainhash.NewHashFromStr("2179bf15a81f1d3f4d9fbe135e2e5f1559eba9f080291339d36bcb35cdb89ffc")
	point := wire.NewOutPoint(txHash, uint32(0))

	/* 交易的锁定脚本，对应的是 ScriptPubKey 字段, NewTxOut 对应未花费的金额
	"vout":
		{
			"scriptpubkey": "51203f2ea7fb55982af54b8762b123b4eb775988b3d8d76050baa4184d4b8ee0665f",
			"scriptpubkey_asm": "OP_PUSHNUM_1 OP_PUSHBYTES_32 3f2ea7fb55982af54b8762b123b4eb775988b3d8d76050baa4184d4b8ee0665f",
			"scriptpubkey_type": "v1_p2tr",
			"scriptpubkey_address": "tb1p8uh20764nq402ju8v2cj8d8twavc3v7c6as9pw4yrpx5hrhqve0seupm29",
			"value": 20451
		}
	*/
	script, _ := hex.DecodeString("51203f2ea7fb55982af54b8762b123b4eb775988b3d8d76050baa4184d4b8ee0665f")
	output := wire.NewTxOut(int64(20451), script)
	fetcher := txscript.NewMultiPrevOutFetcher(nil)
	fetcher.AddPrevOut(*point, output)

	return point, fetcher
}

func DecodeTaprootAddress(strAddr string, cfg *chaincfg.Params) ([]byte, error) {
	taprootAddr, err := btcutil.DecodeAddress(strAddr, cfg)
	if err != nil {
		return nil, err
	}

	byteAddr, err := txscript.PayToAddrScript(taprootAddr)
	if err != nil {
		return nil, err
	}
	return byteAddr, nil
}

func transaction() {
	network := &chaincfg.TestNet3Params
	wif, _ := btcutil.DecodeWIF("cPi6St8HSQdE5Ayq8pPs8Wp61endKKBtS5bUokTyu3BNvNThmVpR")

	taprootAddr, _ := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wif.PrivKey.PubKey())), network)
	fmt.Printf("Taproot testnet address: %s\n", taprootAddr.String())

	point, fetcher := GetUnspent()

	toAddr := "tb1p4eqwtz7ny4cqkmwl4x67jq4pxxzlwnkjn52r4hy032g7g2c5e98s9cqa4y"
	toAddrByte, _ := DecodeTaprootAddress(toAddr, network)

	tx := wire.NewMsgTx(2)
	in := wire.NewTxIn(point, nil, nil)
	tx.AddTxIn(in)
	out := wire.NewTxOut(int64(10000), toAddrByte)
	tx.AddTxOut(out)

	// 打印交易输入和输出的详细信息
	for i, txIn := range tx.TxIn {
		fmt.Printf("Input %d: PreviousOutPoint: %s:%d\n", i, txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index)
	}

	for i, txOut := range tx.TxOut {
		fmt.Printf("Output %d: Value: %d, Script: %x\n", i, txOut.Value, txOut.PkScript)
	}

	// 打印输入和输出的总值
	var totalInputValue int64
	for _, txIn := range tx.TxIn {
		prevOutput := fetcher.FetchPrevOutput(txIn.PreviousOutPoint)
		totalInputValue += prevOutput.Value
	}

	var totalOutputValue int64
	for _, txOut := range tx.TxOut {
		totalOutputValue += txOut.Value
	}

	fmt.Printf("Total Input Value: %d\n", totalInputValue)
	fmt.Printf("Total Output Value: %d\n", totalOutputValue)

	// 获取前一笔交易输出
	prevOutput := fetcher.FetchPrevOutput(in.PreviousOutPoint)
	witness, _ := txscript.TaprootWitnessSignature(tx,
		txscript.NewTxSigHashes(tx, fetcher), 0, prevOutput.Value,
		prevOutput.PkScript, txscript.SigHashDefault, wif.PrivKey)
	tx.TxIn[0].Witness = witness

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)
	finalRawTx := hex.EncodeToString(signedTx.Bytes())

	fmt.Printf("Signed Transaction:\n%s", finalRawTx)
}
