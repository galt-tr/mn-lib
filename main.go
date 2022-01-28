package main

import (
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	woc "github.com/galt-tr/mn-lib/internal/woc"
	metanet "github.com/galt-tr/mn-lib/metanet"
)

func main() {
	parentPrivKey, _ := bsv.WifToPrivateKey("L4PF5nCubsEURgeZW6qCZs92CVesLEHtqvmcGXEwy4ppxdQNYCg9")
	childPrivKey, _ := bsv.WifToPrivateKey("KyT6vXqr4LMENqFPBxkqtn3AkdErtYWgssNQz1rPTrmUTuYmHNtn")
	pubKey := bsv.PubKeyFromPrivateKey(childPrivKey, true)
	address, _ := bsv.GetAddressFromPubKeyString(pubKey, true)

	var sats uint64
	var vOut uint32

	txId := "58fc9319d7b9d18a92724e3b7e95ceba99ccd780b7995362d44cbee8bc46b545"
	vOut = 1

	o, _ := woc.GetTransactionOutput(txId, int(vOut))

	sats = uint64(o.Value * 100000000)
	scriptPubKey := o.ScriptPubKey.Hex

	node := &metanet.MetanetNode{
		Prefix:        "meta",
		NodeAddress:   address.String(),
		NodePublicKey: pubKey,
		ParentTxId:    txId,
		Input: []*bsv.Utxo{&bsv.Utxo{
			Satoshis:     sats,
			ScriptPubKey: scriptPubKey,
			TxID:         txId,
			Vout:         vOut,
		}},
		InputPrivateKey: parentPrivKey,
		ChangeAddress:   address.String(),
		Data:            "testing",
	}
	rawTx, err := metanet.CreateSpendableNode(node)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rawTx)
}
