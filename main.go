package main

import (
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	metanet "github.com/galt-tr/mn-lib/metanet"
)

func main() {
	privKey, _ := bsv.WifToPrivateKey("L1AyB7th6Waq5AhPy4X7pdj5F1mM8Bk1e2iR4nNqLa4UiuasznfA")
	privKey2, _ := bsv.WifToPrivateKey("L19tF5jckU6RajrMzyCiLpL8Rxyp9qNkKiAgXqWNmuGp23gVbJbY")
	fmt.Println(bsv.PubKeyFromPrivateKey(privKey, true))
	pubKey := bsv.PubKeyFromPrivateKey(privKey, true)
	address, _ := bsv.GetAddressFromPubKeyString(pubKey, true)

	node := &metanet.MetanetNode{
		Prefix:        "meta",
		NodeAddress:   address.String(),
		NodePublicKey: pubKey,
		ParentTxId:    "eec332437097611afc4c09856c7ffe17da6b2d66ffa776a6a57afe9d98ebba56",
		Input: []*bsv.Utxo{&bsv.Utxo{
			Satoshis:     5000,
			ScriptPubKey: "14f987be172d68cb99e8e51a8dada6e0bc38a1d54751795a5a9554937f517f7c817f7b6d01157f5e5a9558937f01217f517f01207f557f01147f527f75777e777e7b7c7ea77b885179880079aa517f7c818b7c7e263044022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179802207c7e01417e2102b405d7f0322a89d0f9f3a98e6f938fdc1c969a8d1382a2bf66a71ae74a1e83b0ad046d65746121027271a7a975d4e9decf81b5430a54481ebc1e3226e392069dbb07e66726992a9a20ec514b43feb34359d6382f57e36cb845e8ac0a04f7fbf6f5347078155759e058756d76a914e41b788a1c96d6cf878feb1bd717e82015415faf88ac6a0774657374696e67",
			TxID:         "eec332437097611afc4c09856c7ffe17da6b2d66ffa776a6a57afe9d98ebba56",
			Vout:         0,
		}},
		InputPrivateKey: privKey2,
		ChangeAddress:   address.String(),
		Data:            "testing",
	}
	//hex, preimage, err := metanet.CreateSpendableNode(node)
	//if err != nil {
	//	fmt.Println(err)
	//}
	rawTx, _ := metanet.CreateSpendableNode(node)
	//fmt.Println(hex)
	//fmt.Println(preimage)
	//fmt.Println()
	//p, err := pushtx.ParsePreimageHex(preimage)
	//if err != nil {
	//	fmt.Println(err)
	//}
	fmt.Println(rawTx)
	//hex := "0100000001a4a06f99bcb3f5827824cd42296307fb4d6f69c51db67e39f85d03799aac867c010000006a47304402203af515b6308515001eba7980712be3015108f059937e258e34b4d3aa8b0761360220587133fb52ac0d9f6aa929677b4c9ee71145fbd5354c5e614bbfeb9d25228bfd412103559216cce466d1b8134f83c088f98467cdc462a16e65b3359af63bc9c81e1fb8ffffffff03e8030000000000001976a9148e1fdac4e443d616884f53162787cddd7808e9ff88acd4190000000000001976a9149a25d691778a9933a56c56eaa115064a0acdca2888ac0000000000000000aa006a046d657461223144785637587954616f3351543351374a77676f6d48787076704d677a444d6f67534036613532323636623664376335316639326537393938656562633435643439636337393263633264616132366637386439646639653762656261633563396262017c086272697465767565027631146c6f636174696f6e526576696577506172656e741b4368494a317779514c5178742d544952505a59413652314c70764d00000000"
	//boolean, _ := metanet.IsChildNodeValid(hex, 0, 2)
	//fmt.Println(boolean)
}
