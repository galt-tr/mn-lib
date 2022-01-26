package main

import (
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	metanet "github.com/galt-tr/mn-lib/metanet"
	pushtx "github.com/murray-distributed-technologies/go-pushtx"
)

func main() {
	privKey, _ := bsv.WifToPrivateKey("L4PF5nCubsEURgeZW6qCZs92CVesLEHtqvmcGXEwy4ppxdQNYCg9")
	node := &metanet.MetanetNode{
		Prefix:      "meta",
		NodeAddress: "1PHqZ8QcpiewjntKhXrztkXi9E4Lc5aVmx",
		ParentTxId:  "58fc9319d7b9d18a92724e3b7e95ceba99ccd780b7995362d44cbee8bc46b545",
		Input: []*bsv.Utxo{&bsv.Utxo{
			Satoshis:     2000,
			ScriptPubKey: "0079aa517f7c818b7c7e263044022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179802207c7e01417e2102b405d7f0322a89d0f9f3a98e6f938fdc1c969a8d1382a2bf66a71ae74a1e83b0ad7514cb030491157b26a570b6ee91e5b068d99c3b72f6046d657461a72231346e64483972374e327072396d71793666536f635955483263534a707644394a71044e554c4c7176a9142989611fd22fb65e8d6bb2c2b4e3a2b10dc604dd88ad6d876a0774657374696e67",
			TxID:         "58fc9319d7b9d18a92724e3b7e95ceba99ccd780b7995362d44cbee8bc46b545",
			Vout:         0,
		}},
		InputPrivateKey: privKey,
		ChangeAddress:   "1PHqZ8QcpiewjntKhXrztkXi9E4Lc5aVmx",
		Data:            "testing",
	}
	//hex, preimage, err := metanet.CreateSpendableNode(node)
	//if err != nil {
	//	fmt.Println(err)
	//}
	hex, preimage, err := metanet.CreateSpendableNode(node)
	fmt.Println(hex)
	fmt.Println(preimage)
	p, err := pushtx.ParsePreimageHex(preimage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)
	//hex := "0100000001a4a06f99bcb3f5827824cd42296307fb4d6f69c51db67e39f85d03799aac867c010000006a47304402203af515b6308515001eba7980712be3015108f059937e258e34b4d3aa8b0761360220587133fb52ac0d9f6aa929677b4c9ee71145fbd5354c5e614bbfeb9d25228bfd412103559216cce466d1b8134f83c088f98467cdc462a16e65b3359af63bc9c81e1fb8ffffffff03e8030000000000001976a9148e1fdac4e443d616884f53162787cddd7808e9ff88acd4190000000000001976a9149a25d691778a9933a56c56eaa115064a0acdca2888ac0000000000000000aa006a046d657461223144785637587954616f3351543351374a77676f6d48787076704d677a444d6f67534036613532323636623664376335316639326537393938656562633435643439636337393263633264616132366637386439646639653762656261633563396262017c086272697465767565027631146c6f636174696f6e526576696577506172656e741b4368494a317779514c5178742d544952505a59413652314c70764d00000000"
	//boolean, _ := metanet.IsChildNodeValid(hex, 0, 2)
	//fmt.Println(boolean)
}
