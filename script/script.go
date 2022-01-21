package script

import (
	"encoding/hex"

	"github.com/libsv/go-bt/bscript"
)

func NewMetanetP2PKH(address, parentTxId string) (*bscript.Script, error) {
	s := &bscript.Script{}

	var err error

	//Push SHA1 Hash for Filtering
	var hashBytes []byte
	if hashBytes, err = hex.DecodeString("8e9c49fd4e791448110a80548eb01783723d4deb"); err != nil {
		return nil, err
	}
	if err = s.AppendPushData(hashBytes); err != nil {
		return nil, err
	}
	//if err = s.AppendPushDataString("CB030491157B26A570B6EE91E5B068D99C3B72F6"); err != nil {
	//	return nil, err
	//}

	//append meta flag
	if err = s.AppendPushDataString("meta"); err != nil {
		return nil, err
	}

	//append OP_SHA1
	s.AppendOpCode(bscript.OpSHA1)

	//append node address
	if err = s.AppendPushDataString(address); err != nil {
		return nil, err
	}
	//append parentTxId
	if err = s.AppendPushDataString(parentTxId); err != nil {
		return nil, err
	}

	//Rotate signature and pubkey to top of stack
	s.AppendOpCode(bscript.Op2ROT)

	//Append P2PKH OPCODES
	s.AppendOpCode(bscript.OpDUP)
	s.AppendOpCode(bscript.OpHASH160)

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

	s.AppendOpCode(bscript.OpEQUALVERIFY)

	s.AppendOpCode(bscript.OpCHECKSIGVERIFY)
	//s.AppendOpCode(bscript.OpROT)
	s.AppendOpCode(bscript.Op2DROP)
	s.AppendOpCode(bscript.OpEQUAL)

	return s, nil

}
