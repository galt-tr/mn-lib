package main

import (
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	metanet "github.com/galt-tr/mn-lib/metanet"
)

func main() {
	privKey, _ := bsv.WifToPrivateKey("L1qKK6yz2ZEXwe9SptJDKs3G1BR2Np828uQiq2ojqsPLafYbze2r")
	node := &metanet.MetanetNode{
		Prefix:      "meta",
		NodeAddress: "166Uuce9aPFftGcFUgTj5p6NAzLwyEYHYT",
		ParentTxId:  "NULL",
		Input: []*bsv.Utxo{&bsv.Utxo{
			Satoshis:     4754815,
			ScriptPubKey: "76a914dfe61ca18783253b027cd2cd1f655ce68f9c470888ac",
			TxID:         "3b9cbdd54d382a9f23e6030623ccd4f49728675ca897f933de54c465a5ae7847",
			Vout:         1,
		}},
		InputPrivateKey: privKey,
		ChangeAddress:   "166Uuce9aPFftGcFUgTj5p6NAzLwyEYHYT",
		Data:            "testing",
	}
	hex, preimage, err := metanet.CreateSpendableNode(node)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(hex)
	fmt.Println(preimage)
	//hex := "0100000001a4a06f99bcb3f5827824cd42296307fb4d6f69c51db67e39f85d03799aac867c010000006a47304402203af515b6308515001eba7980712be3015108f059937e258e34b4d3aa8b0761360220587133fb52ac0d9f6aa929677b4c9ee71145fbd5354c5e614bbfeb9d25228bfd412103559216cce466d1b8134f83c088f98467cdc462a16e65b3359af63bc9c81e1fb8ffffffff03e8030000000000001976a9148e1fdac4e443d616884f53162787cddd7808e9ff88acd4190000000000001976a9149a25d691778a9933a56c56eaa115064a0acdca2888ac0000000000000000aa006a046d657461223144785637587954616f3351543351374a77676f6d48787076704d677a444d6f67534036613532323636623664376335316639326537393938656562633435643439636337393263633264616132366637386439646639653762656261633563396262017c086272697465767565027631146c6f636174696f6e526576696577506172656e741b4368494a317779514c5178742d544952505a59413652314c70764d00000000"
	//boolean, _ := metanet.IsChildNodeValid(hex, 0, 2)
	//fmt.Println(boolean)
}
