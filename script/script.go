package script

import (
	"encoding/hex"

	"github.com/libsv/go-bt/bscript"
)

func NewMetanetP2PKH(address, parentTxId string) (*bscript.Script, error) {
	s := bscript.Script{}
	//append meta flag
	if err = s.AppendPushDataString("meta"); err != nil {
		return nil, err
	}

	//append parentTxId
	if err = s.AppendPushDataString(parentTxId); err != nil {
		return nil, err
	}

	//append node address
	if err = s.AppendPushdataString(address); err != nil {
		return nil, err
	}

	//Append P2PKH OPCODES
	if err = s.AppendOpcodes(OpDUP, OpHash160); err != nil {

		return nil, err
	}
	//Get PubKeyHash and Check Validity
	a, err := bscript.NewAddressFromString(address)
	if err != nil {
		return nil, err
	}

	var pubKeyHashBytes []byte
	if pubKeyHashBytes, err = hex.DecodeString(a.PublicKeyHash); err != nil {
		return nil, err
	}

	if err = s.AppendPushData(pubKeyHashBytes); err != nil {
		return nil, err
	}

	_ = s.AppendOpcodes(OpEQUALVERIFY, OpCHECKSIG)

	return s, nil

}
