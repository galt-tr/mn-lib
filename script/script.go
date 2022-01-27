package script

import (
	"encoding/hex"

	pushtx "github.com/murray-distributed-technologies/go-pushtx/script"

	bsv "github.com/bitcoinschema/go-bitcoin"

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

func NewMetanetP2PKH(pubKey, parentTxId, data string) (*bscript.Script, error) {
	s := &bscript.Script{}

	var err error

	//Push SHA1 Hash for Filtering
	// SHA1 Hash of template is 'f987be172d68cb99e8e51a8dada6e0bc38a1d547'
	// hex: '51795a5a9554937f517f7c817f7b6d01157f5e5a9558937f01217f517f01207f557f01147f527f75777e777e7b7c7ea77b885179880079aa517f7c818b7c7e263044022079be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179802207c7e01417e2102b405d7f0322a89d0f9f3a98e6f938fdc1c969a8d1382a2bf66a71ae74a1e83b0ad046d6574612120756d76a91488ac'
	if err = s.AppendPushDataHexString("f987be172d68cb99e8e51a8dada6e0bc38a1d547"); err != nil {
		return nil, err
	}

	// grab preimage off the stack
	s.AppendOpCode(bscript.Op1)
	s.AppendOpCode(bscript.OpPICK)

	// get locking script from preimage

	if s, err = pushtx.AppendGetLockingScriptFromPreimage(s); err != nil {
		return nil, err
	}

	// strip data from script template and check hash against sha1 hash of template
	// also checks that pubkey of parent is same as pubkey used to sign input

	if s, err = StripTemplateData(s); err != nil {
		return nil, err
	}

	// add PushTX Check to validate we are in current transaction
	s, err = pushtx.AppendPushTx(s)
	if err != nil {
		return nil, err
	}

	//append meta flag

	if err = s.AppendPushDataString("meta"); err != nil {
		return nil, err
	}

	//append pubkey
	if err = s.AppendPushDataHexString(pubKey); err != nil {
		return nil, err
	}

	//append txid
	if err = s.AppendPushDataHexString(parentTxId); err != nil {
		return nil, err
	}
	s.AppendOpCode(bscript.OpDROP)
	s.AppendOpCode(bscript.Op2DROP)

	address, err := bsv.GetAddressFromPubKeyString(pubKey, true)
	if err != nil {
		return nil, err
	}

	//Add P2PKH Script
	if s, err = AppendP2PKHLockingScript(s, address.String()); err != nil {
		return nil, err
	}

	if data != "" {
		s.AppendOpCode(bscript.OpRETURN)
		s.AppendPushDataString(data)
	}

	return s, nil
}

func AppendMetanetUnlockingScript(s *bscript.Script, preimage []byte) (*bscript.Script, error) {
	var err error
	if preimage != nil {
		if err = s.AppendPushData(preimage); err != nil {
			return nil, err
		}
	}

	return s, nil
}

//Currently the same as Metanet Unlocking Script, but separating for future

func NewOpPushTxUnlockingScript(pubKey, preimage, sig []byte, sigHashFlag sighash.Flag) (*bscript.Script, error) {
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

func AppendP2PKHLockingScript(s *bscript.Script, address string) (*bscript.Script, error) {
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

	return s, nil
}

// 25 bytes
// assumes locking script is on top of the stack
func StripTemplateData(s *bscript.Script) (*bscript.Script, error) {

	var err error

	// split 21 byte (0x15) hash from script
	if err = s.AppendPushDataHexString("15"); err != nil {
		return nil, err
	}
	s.AppendOpCode(bscript.OpSPLIT)

	//push 148  to stack to split template after 'meta'+ pushdata prefix for pubkey
	// i don't think I can push hex str need to instead get 139 manually i think
	//if err = s.AppendPushDataHexString("79"); err != nil {
	//	return nil, err
	//}
	s.AppendOpCode(bscript.Op14)
	s.AppendOpCode(bscript.Op10)
	s.AppendOpCode(bscript.OpMUL)
	s.AppendOpCode(bscript.Op8)
	s.AppendOpCode(bscript.OpADD)

	s.AppendOpCode(bscript.OpSPLIT)

	//TODO: Should be smart and grab data prefix to better split the template. This would allow for data of non-fixed lengths

	// push 33 (0x21) to stack to split template after public key of node
	if err = s.AppendPushDataHexString("21"); err != nil {
		return nil, err
	}

	s.AppendOpCode(bscript.OpSPLIT)
	// split first data prefix byte from txid
	s.AppendOpCode(bscript.Op1)
	s.AppendOpCode(bscript.OpSPLIT)

	// push 32 (0x20) to stack to split template after parent txid
	if err = s.AppendPushDataHexString("20"); err != nil {
		return nil, err
	}

	s.AppendOpCode(bscript.OpSPLIT)

	// push Op_5 to split off OP_DROP, OP_DROP, and first 3 bytes of P2PKH script
	s.AppendOpCode(bscript.Op5)
	s.AppendOpCode(bscript.OpSPLIT)

	// push 20 (0x14) to stack to split off pubkeyhash in P2PKH script
	if err = s.AppendPushDataHexString("14"); err != nil {
		return nil, err
	}
	s.AppendOpCode(bscript.OpSPLIT)

	s.AppendOpCode(bscript.Op2)
	s.AppendOpCode(bscript.OpSPLIT)

	// drop anything after p2pkh
	s.AppendOpCode(bscript.OpDROP)
	// nip pubkeyhash
	s.AppendOpCode(bscript.OpNIP)
	// concatenate script template
	s.AppendOpCode(bscript.OpCAT)
	s.AppendOpCode(bscript.OpNIP)
	s.AppendOpCode(bscript.OpCAT)
	s.AppendOpCode(bscript.OpROT)
	s.AppendOpCode(bscript.OpSWAP)
	s.AppendOpCode(bscript.OpCAT)

	// SHA1 hash the script template
	s.AppendOpCode(bscript.OpSHA1)

	// move pubkey to bottom of stack and check hashes
	s.AppendOpCode(bscript.OpROT)
	s.AppendOpCode(bscript.OpEQUALVERIFY)

	// grab public key and check they are equal for valid metanet node
	s.AppendOpCode(bscript.Op1)
	s.AppendOpCode(bscript.OpPICK)
	s.AppendOpCode(bscript.OpEQUALVERIFY)

	return s, nil

}
