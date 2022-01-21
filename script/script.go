package script

import (
	"encoding/hex"

	"github.com/libsv/go-bt/bscript"
)

func NewSha1HashPuzzle(str, hash string) (*bscript.Script, error) {
	s := &bscript.Script{}
	var err error

	//Push SHA1 Hash for Filtering
	var hashBytes []byte
	if hashBytes, err = hex.DecodeString(hash); err != nil {
		return nil, err
	}
	if err = s.AppendPushData(hashBytes); err != nil {
		return nil, err
	}

	if err = s.AppendPushDataString(str); err != nil {
		return nil, err
	}

	s.AppendOpCode(bscript.OpSHA1)
	s.AppendOpCode(bscript.OpEQUAL)
	return s, nil
}

func NewMetanetP2PKH(address, parentTxId string) (*bscript.Script, error) {
	s := &bscript.Script{}

	var err error

	//Push SHA1 Hash for Filtering
	var hashBytes []byte
	if hashBytes, err = hex.DecodeString("e58a58b0d3f0744e22339f8068db085ada2e2e82"); err != nil {
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
	s.AppendOpCode(bscript.OpEQUALVERIFY)

	//append node address
	//if err = s.AppendPushDataString(address); err != nil {
	//	return nil, err
	//}
	//append parentTxId
	//if err = s.AppendPushDataString(parentTxId); err != nil {
	//	return nil, err
	//}
	//s.AppendOpCode(bscript.Op2DROP)
	//s.AppendOpCode(bscript.OpDROP)

	//Rotate signature and pubkey to top of stack
	//s.AppendOpCode(bscript.Op2ROT)

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

	s.AppendOpCode(bscript.OpCHECKSIG)
	//s.AppendOpCode(bscript.OpROT)
	//s.AppendOpCode(bscript.Op2DROP)
	//s.AppendOpCode(bscript.OpEQUAL)

	return s, nil

}
