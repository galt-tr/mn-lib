package script

import (
	"encoding/hex"

	"github.com/libsv/go-bt/bscript"
	"github.com/libsv/go-bt/sighash"
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

func NewMetanetP2PKH(address, parentTxId, data string) (*bscript.Script, error) {
	s := &bscript.Script{}

	var err error

	s, err = AppendPushTx(s)
	if err != nil {
		return nil, err
	}

	//Push SHA1 Hash for Filtering
	if err = s.AppendPushDataHexString("cb030491157b26a570b6ee91e5b068d99c3b72f6"); err != nil {
		return nil, err
	}

	//append meta flag
	var str string
	str = "meta"
	if err = s.AppendPushDataString(str); err != nil {
		return nil, err
	}

	//append OP_SHA1
	s.AppendOpCode(bscript.OpSHA1)

	//append node address
	if err = s.AppendPushDataString(address); err != nil {
		return nil, err
	}
	// append parentTxId
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
	s.AppendOpCode(bscript.Op2DROP)
	s.AppendOpCode(bscript.OpEQUAL)
	if data != "" {
		s.AppendOpCode(bscript.OpRETURN)
		s.AppendPushDataString(data)
	}

	return s, nil
}

// AppendPushTx assumes preimage in the unlocking script
// Leaves a copy of the preimage on the stack

func AppendPushTx(s *bscript.Script) (*bscript.Script, error) {
	var err error
	//Add number of items back in the stack preimage is
	s.AppendOpCode(bscript.Op0)
	//Copy preimage to top of the stack
	s.AppendOpCode(bscript.OpPICK)
	//Double SHA256 hash preimage
	s.AppendOpCode(bscript.OpHASH256)
	//Split first byte of preimage hash, move to top of the stack, convert to num, and add 1
	s.AppendOpCode(bscript.Op1)
	s.AppendOpCode(bscript.OpSPLIT)
	s.AppendOpCode(bscript.OpSWAP)
	s.AppendOpCode(bscript.OpBIN2NUM)
	s.AppendOpCode(bscript.Op1ADD)

	//Concatenate new first byte with preimage hash
	s.AppendOpCode(bscript.OpSWAP)
	s.AppendOpCode(bscript.OpCAT)

	//Push Private Key derived from hash
	if err = s.AppendPushDataHexString("3044022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f817980220"); err != nil {
		return nil, err
	}

	s.AppendOpCode(bscript.OpSWAP)
	s.AppendOpCode(bscript.OpCAT)
	//Not working in mainnet, so going to calculate this manually
	//s.AppendOpCode(bscript.OpDATA65)
	s.AppendOpCode(bscript.Op16)
	s.AppendOpCode(bscript.Op4)
	s.AppendOpCode(bscript.OpMUL)
	s.AppendOpCode(bscript.Op1ADD)
	s.AppendOpCode(bscript.OpCAT)

	pubKeyHashBytes, err := hex.DecodeString("02b405d7f0322a89d0f9f3a98e6f938fdc1c969a8d1382a2bf66a71ae74a1e83b0")
	if err != nil {
		return nil, err
	}
	if err = s.AppendPushData(pubKeyHashBytes); err != nil {
		return nil, err
	}
	s.AppendOpCode(bscript.OpCHECKSIGVERIFY)

	// DROPPING FOR TESTING
	s.AppendOpCode(bscript.OpDROP)
	return s, nil
}

func NewMetanetUnlockingScript(pubKey, preimage, sig []byte, sigHashFlag sighash.Flag) (*bscript.Script, error) {
	sigBuf := []byte{}
	sigBuf = append(sigBuf, sig...)
	sigBuf = append(sigBuf, uint8(sigHashFlag))

	scriptBuf := [][]byte{sigBuf, pubKey}
	s := &bscript.Script{}
	err := s.AppendPushDataArray(scriptBuf)
	if err != nil {
		return nil, err
	}

	if preimage != nil {
		if err = s.AppendPushData(preimage); err != nil {
			return nil, err
		}
	}

	return s, nil

}
