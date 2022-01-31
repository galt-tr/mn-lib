package main

import (
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	woc "github.com/galt-tr/mn-lib/internal/woc"
	metanet "github.com/galt-tr/mn-lib/metanet"
)

func main() {
	childPrivKey, _ := bsv.WifToPrivateKey("L54vKPJjxy7GueTjmbXZXDbhM3qJKxTge4Nny6txyV71maT5XJqk")
	parentPrivKey, _ := bsv.WifToPrivateKey("L4wdGhsyeGhGG7kmDZUav8qBRdavZh8ixb4oJ1YfULey9btQe7XU")
	pubKey := bsv.PubKeyFromPrivateKey(childPrivKey, true)
	address, _ := bsv.GetAddressFromPubKeyString(pubKey, true)

	var sats uint64
	var vOut uint32

	txId := "a67b0067999c734363e6f1f5743fc490a72e19c687db9529a2bf1b69bc0e1586"
	vOut = 0

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
	rawTx, err := metanet.CreateOpPushTx(node)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rawTx)
}
