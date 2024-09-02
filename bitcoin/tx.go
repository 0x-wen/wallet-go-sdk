package bitcoin

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

var (
	privKey     = "cPi6St8HSQdE5Ayq8pPs8Wp61endKKBtS5bUokTyu3BNvNThmVpR"
	destination = "n3hpgkysvCBwpiAVQVEFQcS4B6i1RUyh7k"
	
)

type Data struct {
	// PreviousTxid string `json:"prev_tx_hash"`
	// Balance      int64  `json:"balance"`
	// PubKeyScript string `json:"pub_key_script"`
	My []BlockChairResp
}


type Prevout struct {
	Scriptpubkey string `json:"scriptpubkey"`
	Scriptpubkey_asm string `json:"scriptpubkey_asm"`
	Scriptpubkey_type string `json:"scriptpubkey_type"`
	Scriptpubkey_address string `json:"scriptpubkey_address"`
	Value int `json:"value"`
}

type Vin struct {
	Txid string `json:"txid"`
	Vout int `json:"vout"`
	Prevout Prevout `json:"prevout"`
	Scriptsig string `json:"scriptsig"`
	Scriptsig_asm string `json:"scriptsig_asm"`
	Witness []string `json:"witness"`
	Is_coinbase bool `json:"is_coinbase"`
	Sequence int `json:"sequence"`
}

type Vout struct {
	Scriptpubkey string `json:"scriptpubkey"`
	Scriptpubkey_asm string `json:"scriptpubkey_asm"`
	Scriptpubkey_type string `json:"scriptpubkey_type"`
	Scriptpubkey_address string `json:"scriptpubkey_address"`
	Value int `json:"value"`
}

type Status struct {
	Confirmed bool `json:"confirmed"`
	Block_height int `json:"block_height"`
	Block_hash string `json:"block_hash"`
	Block_time int `json:"block_time"`
}

type BlockChairResp struct {
	Txid string `json:"txid"`
	Version int `json:"version"`
	Locktime int `json:"locktime"`
	Vin []Vin `json:"vin"`
	Vout []Vout `json:"vout"`
	Size int `json:"size"`
	Weight int `json:"weight"`
	Sigops int `json:"sigops"`
	Fee int `json:"fee"`
	Status Status `json:"status"`
}


func Tx() {
	rawTx, err := CreateTx(privKey, destination, 30000)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("raw signed transaction is: ", rawTx)
}

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func GetUTXO(address string) (string, int64, string, error) {

	// Provide your url to get UTXOs, read the response
	// unmarshal it, and extract necessary data
	newURL := fmt.Sprintf("https://mempool.space/testnet/api/address/%s/txs", address)

	response, err := http.Get(newURL)
	if err != nil {
		fmt.Println("error in GetUTXO, http.Get")
		return "", 0, "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error in GetUTXO, io.ReadAll")
		return "", 0, "", err
	}

	// based on the response you get, should define a struct
	// so before unmarshaling check your JSON response model

	var blockChairResp = Data{}
	err = json.Unmarshal(body, &blockChairResp)
	if err != nil {
		fmt.Println("error in GetUTXO, json.Unmarshal")
		return "", 0, "", err
	}

	// var previousTxid string = "16688d2946c3e029ca91ce730109994c2bcafb859d580a6f7c820fb2bb5b6afc"
	// var balance int64 = 62000
	// var pubKeyScript string = "76a91455d5e92958a8b06b4ff15cd2dd3d254f375e98db88ac"
	return blockChairResp.My[0].Txid, int64(blockChairResp.My[0].Vout[1].Value), blockChairResp.My[0].Vout[1].Scriptpubkey, nil
}

// CreateTx creates a transaction to send a specified amount of bitcoins to a destination address.
// privKey: The private key used to sign the transaction.
// destination: The address of the recipient.
// amount: The amount of bitcoins to send.
// Returns the signed transaction in hexadecimal string format and an error if any occurs.
func CreateTx(privKey string, destination string, amount int64) (string, error) {
	// Decode the private key in WIF format.
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// Generate a public key address from the private key for retrieving UTXOs.
	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), TestNet3Params)
	if err != nil {
		return "", err
	}

	// Retrieve the UTXO (unspent transaction output) information of the address, including transaction ID, balance, and script.
	txid, balance, pkScript, err := GetUTXO(addrPubKey.EncodeAddress())
	if err != nil {
		return "", err
	}

	// Check if the balance is sufficient to cover the amount to be sent.
	/*
	 * 1 or unit-amount in Bitcoin is equal to 1 satoshi and 1 Bitcoin = 100000000 satoshi
	 */
	// checking for sufficiency of account
	if balance < amount {
		return "", fmt.Errorf("the balance of the account is not sufficient")
	}

	// Decode the destination address and prepare the payment script.
	// extracting destination address as []byte from function argument (destination string)
	destinationAddr, err := btcutil.DecodeAddress(destination, TestNet3Params)
	if err != nil {
		return "", err
	}
	destinationAddrByte, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return "", err
	}

	// Create a new transaction.
	// creating a new bitcoin transaction, different sections of the tx, including
	// input list (contain UTXOs) and outputlist (contain destination address and usually our address)
	// in next steps, sections will be field and pass to sign
	redeemTx, err := NewTx()
	if err != nil {
		return "", err
	}

	// Convert the transaction ID string to a hash object.
	utxoHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return "", err
	}

	// Create an outpoint to specify the UTXO to be spent.
	// the second argument is vout or Tx-index, which is the index
	// of spending UTXO in the transaction that Txid referred to
	// in this case is 1, but can vary different numbers
	outPoint := wire.NewOutPoint(utxoHash, 1)

	// Create a transaction input and add it to the transaction.
	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// Create a transaction output, specifying the recipient address and amount, and add it to the transaction.
	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	// Sign the transaction using the private key and the corresponding script.
	// now sign the transaction
	finalRawTx, err := SignTx(privKey, pkScript, redeemTx)
	if err != nil {
		return "", err
	}

	return finalRawTx, nil
}

func SignTx(privKey string, pkScript string, redeemTx *wire.MsgTx) (string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return "", nil
	}

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePKScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return "", nil
	}

	// since there is only one input, and want to add
	// signature to it use 0 as index
	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}
