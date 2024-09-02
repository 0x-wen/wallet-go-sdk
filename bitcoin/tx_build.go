package bitcoin

import (
	"errors"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// Input 定义了交易输入的相关信息。
// 它包含了交易的唯一标识、输出序号、私钥、赎回脚本、地址和金额。
type Input struct {
	txId          string // 交易ID，用于唯一标识一个交易。
	vOut          uint32 // 交易输出序号，指定具体的交易输出。
	privateKeyHex string // 私钥的十六进制表示，用于签名交易。
	redeemScript  string // 赎回脚本，用于解锁特定类型的交易输出。
	address       string // 地址，交易输出的接收方地址。
	amount        int64  // 交易金额，指定输入的金额。
}

// Output 结构体代表一个交易输出。
// 它包含接收比特币的地址、脚本和对应的金额。
type Output struct {
	address string
	script  string
	amount  int64
}

type TransactionBuilder struct {
	inputs    []Input
	outputs   []Output
	netParams *chaincfg.Params
	tx        *wire.MsgTx
}

func NewTxBuild(version int32, netParams *chaincfg.Params) *TransactionBuilder {
	if netParams == nil {
		netParams = &chaincfg.MainNetParams
	}
	builder := &TransactionBuilder{
		inputs:    nil,
		outputs:   nil,
		netParams: netParams,
		tx:        &wire.MsgTx{Version: version, LockTime: 0},
	}
	return builder
}
func (t *TransactionBuilder) TotalInputAmount() int64 {
	total := int64(0)
	for _, v := range t.inputs {
		total += v.amount
	}
	return total
}

func (t *TransactionBuilder) TotalOutputAmount() int64 {
	total := int64(0)
	for _, v := range t.outputs {
		total += v.amount
	}
	return total
}

func (t *TransactionBuilder) AddInput(txId string, vOut uint32, privateKeyHex string,
	redeemScript string, address string, amount int64) {
	input := Input{txId: txId, vOut: vOut, privateKeyHex: privateKeyHex,
		redeemScript: redeemScript, address: address, amount: amount}
	t.inputs = append(t.inputs, input)
}

func (t *TransactionBuilder) AddOutput(address string, script string, amount int64) {
	output := Output{address: address, script: script, amount: amount}
	t.outputs = append(t.outputs, output)
}


func (t *TransactionBuilder) Build() (*wire.MsgTx, error) {
	if len(t.inputs) == 0 || len(t.outputs) == 0 {
		return nil, errors.New("invalid inputs or outputs")
	}

	tx := t.tx
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	var privateKeys []*btcec.PrivateKey
	for i := 0; i < len(t.inputs); i++ {
		input := t.inputs[i]
		txHash, err := chainhash.NewHashFromStr(input.txId)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(txHash, input.vOut)
		pkScript, err := AddrToPkScript(input.address, t.netParams)
		if err != nil {
			return nil, err
		}
		txOut := wire.NewTxOut(input.amount, pkScript)
		prevOutFetcher.AddPrevOut(*outPoint, txOut)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.TxIn = append(tx.TxIn, txIn)

		wif, err := btcutil.DecodeWIF(input.privateKeyHex)
		if err != nil {
			return nil, err
		}
		privateKeys = append(privateKeys, wif.PrivKey)
	}

	for i := 0; i < len(build.outputs); i++ {
		output := build.outputs[i]

		var pkScript []byte
		var err error
		if len(output.script) != 0 && len(output.address) == 0 {
			pkScript, err = hex.DecodeString(output.script)
			if err != nil {
				return nil, err
			}
		} else {
			pkScript, err = AddrToPkScript(output.address, build.netParams)
			if err != nil {
				return nil, err
			}
		}
		txOut := wire.NewTxOut(output.amount, pkScript)
		tx.TxOut = append(tx.TxOut, txOut)
	}
	if err := Sign(tx, privateKeys, prevOutFetcher); err != nil {
		return nil, err
	}
	return tx, nil
}
