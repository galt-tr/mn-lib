package main

import (
	"fmt"

	"github.com/galt-tr/mn-lib/metanet"
	"github.com/galt-tr/mn-lib/types"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
)

func main() {
	childPrivKey, _ := wif.DecodeWIF("L448vgvBygTDNNRRxPLoxpVbkYQA6xcGDyPgUkahAP4tNAkVa6pD")
	parentPrivKey, _ := wif.DecodeWIF("KzisysrR57Gjd3z7iZBYhVMPU7vQbAw8RnRfM5e4u5WTyWFqmHXv")
	pubKey := childPrivKey.PrivKey.PubKey()
	pubKeyNode := hex.EncodeToString(pubKey.SerializeCompressed())
	fmt.Println(pubKey.SerialiseCompressed())
	address, _ := bscript.NewAddressFromPublicKey(pubKey, true)
	var sats uint64
	var vOut uint32

	txId := "de1452e582ed1ed4b10ad52d2412251d5d756d0c1152e726f4cb4d57369d2cca"
	vOut = 0
	amount := uint64(3000)

	o, _ := woc.GetTransactionOutput(txId, int(vOut))

	sats = uint64(o.Value * 100000000)
	scriptPubKey, err := bscript.NewFromHexString(o.ScriptPubKey.Hex)
	if err != nil {
		fmt.Println(err)
	}

	input := &bt.Input{
		PreviousTxSatoshis: sats,
		PreviousTxScript:   scriptPubKey,
		PreviousTxOutIndex: vOut,
	}
	mn := &types.MetanetNode{
		Prefix:          "meta",
		NodeAddress:     address,
		NodePublicKey:   pubKeyNode,
		ParentTxId:      txId,
		Satoshis:        amount,
		Input:           input,
		InputPrivateKey: parentPrivKey,
		Data:            "parent node",
		ChangeAddress:   NodeAddress,
	}

	rawTx, err := metanet.NewMetanetNode(mn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rawTx)

}
