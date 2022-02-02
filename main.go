package main

import (
	"fmt"

	"github.com/galt-tr/mn-lib/internal/woc"
	"github.com/galt-tr/mn-lib/metanet"
	"github.com/galt-tr/mn-lib/types"
	"github.com/libsv/go-bk/wif"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
)

func main() {
	childPrivKey, _ := wif.DecodeWIF("L4DtXj5wxuzTfhoGY5vNzSxuoxi7MApFpt4fM4g8RrLhW6GHypWV")
	parentPrivKey, _ := wif.DecodeWIF("L4aUoic8n7ofQxgMhtbB1vgkS87kCK3ECCuZVSuTgihsr8CPHGKr")
	pubKey := childPrivKey.PrivKey.PubKey()
	address, _ := bscript.NewAddressFromPublicKey(pubKey, true)
	var sats uint64
	var vOut uint32

	txId := "274c1f33b2b160b54d5b159e818e94204d27eb50e64e1f55165a5319d671d8e8"
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
		NodeAddress:     address.AddressString,
		NodePublicKey:   pubKey.SerialiseCompressed(),
		ParentTxId:      txId,
		Satoshis:        amount,
		Input:           input,
		InputPrivateKey: parentPrivKey.PrivKey,
		Data:            "child node",
		ChangeAddress:   address.AddressString,
	}

	rawTx, err := metanet.NewMetanetNode(mn)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rawTx)

}
