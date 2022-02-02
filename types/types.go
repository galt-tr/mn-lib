package types

import (
	"github.com/libsv/go-bk/bec"
	"github.com/libsv/go-bt/v2"
)

type MetanetNode struct {
	Prefix          string          //Prefix for Metanet Nodes - default to 'meta'
	NodeAddress     string          //Node Public Key Address
	NodePublicKey   []byte          //Public Key of Node
	ParentTxId      string          //Transaction ID of Parent Node
	Satoshis        uint64          //Number of Satoshis to lock in node
	Input           *bt.Input       //utxo of parent node
	InputPrivateKey *bec.PrivateKey //Private Key used to sign input
	Data            string          //Data to be added at end of OP_RETURN in Metanet Node
	ChangeAddress   string          //Use for Change
}
