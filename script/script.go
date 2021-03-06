package script

import (
	"encoding/hex"

	pushtx "github.com/murray-distributed-technologies/go-pushtx/script"

	"github.com/galt-tr/mn-lib/types"
	"github.com/libsv/go-bt/v2/bscript"
	"github.com/libsv/go-bt/v2/sighash"
)

func IsMetanetTx(s *bscript.Script) bool {
	return true
}

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

	s.AppendOpcodes(bscript.OpSHA1)
	s.AppendOpcodes(bscript.OpEQUAL)
	return s, nil
}
func AppendFilterHash(s *bscript.Script) (*bscript.Script, error) {
	var err error
	//Push SHA1 Hash for Filtering
	// SHA1 Hash of template is '318672c5be62f835a01e67488ffd76c59c3e686c'
	if err = s.AppendPushDataHexString("318672c5be62f835a01e67488ffd76c59c3e686c"); err != nil {
		return nil, err
	}
	return s, nil
}

func NewMetanetLockingScript(mn *types.MetanetNode) (*bscript.Script, error) {
	var err error
	s := &bscript.Script{}

	// add hash for filtering
	if s, err = AppendFilterHash(s); err != nil {
		return nil, err
	}

	// grab preimage to top of stack
	s.AppendOpcodes(bscript.Op1, bscript.OpPICK)
	s.AppendOpcodes(bscript.OpDUP)
	if s, err = GetParentTxIDFromPreimage(s); err != nil {
		return nil, err
	}
	s.AppendOpcodes(bscript.Op1, bscript.OpPICK)

	// add get input locking script from preimage
	if s, err = pushtx.AppendGetLockingScriptFromPreimage(s); err != nil {
		return nil, err
	}

	// strip data from locking script to get script template
	// checks that SHA1 hash of script template verifies against filter hash
	// checks that public key in node matches with signature pubkey
	// pushes parentTxID of Metanet Node to ALTSTACK
	if s, err = StripTemplateData(s); err != nil {
		return nil, err
	}

	// check that parentTxID in preimage is equivalent to what is defined in metanet node
	if s, err = GetParentTxIDFromPreimage(s); err != nil {
		return nil, err
	}
	// grab parentTxID from altstack and check that it is equal to preimage txid
	s.AppendOpcodes(bscript.OpFROMALTSTACK)
	s.AppendOpcodes(bscript.OpEQUALVERIFY)
	// add OP_PUSH_TX
	if s, err = pushtx.AppendPushTx(s); err != nil {
		return nil, err
	}

	// append meta flag
	if err = s.AppendPushDataString("meta"); err != nil {
		return nil, err
	}
	// append pubkey
	if err = s.AppendPushData(mn.NodePublicKey); err != nil {
		return nil, err
	}
	// append parent Tx Id
	if err = s.AppendPushDataHexString(mn.ParentTxId); err != nil {
		return nil, err
	}

	//drop metanet data from stack
	s.AppendOpcodes(bscript.OpDROP, bscript.Op2DROP)

	// append P2PKH
	if s, err = pushtx.AppendP2PKH(s, mn.NodeAddress); err != nil {
		return nil, err
	}

	if mn.Data != "" {
		s.AppendOpcodes(bscript.OpRETURN)
		s.AppendPushDataString(mn.Data)
	}
	return s, nil
}

//Currently the same as Metanet Unlocking Script, but separating for future

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

func AppendP2PKHLockingScript(s *bscript.Script, address string) (*bscript.Script, error) {
	//Append P2PKH OPCODES
	s.AppendOpcodes(bscript.OpDUP)
	s.AppendOpcodes(bscript.OpHASH160)

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

	s.AppendOpcodes(bscript.OpEQUALVERIFY)

	s.AppendOpcodes(bscript.OpCHECKSIG)

	return s, nil
}

func GetParentTxIDFromPreimage(s *bscript.Script) (*bscript.Script, error) {
	var err error

	//assumes preimage is on top of stack
	// 68 bytes
	if err = s.AppendPushDataHexString("44"); err != nil {
		return nil, err
	}
	//split
	s.AppendOpcodes(bscript.OpSPLIT)
	s.AppendPushDataHexString("20")
	s.AppendOpcodes(bscript.OpSPLIT)
	s.AppendOpcodes(bscript.OpDROP)
	s.AppendOpcodes(bscript.OpNIP)

	return s, nil
}

// 25 bytes
// assumes locking script is on top of the stack
func StripTemplateData(s *bscript.Script) (*bscript.Script, error) {

	var err error
	// take first byte to check length of filter hash
	// should always be 21 bytes
	s.AppendOpcodes(bscript.Op1)
	s.AppendOpcodes(bscript.OpSPLIT)
	// split 21 byte hash from script
	s.AppendOpcodes(bscript.OpSWAP)
	s.AppendOpcodes(bscript.OpSPLIT)
	// drop filter from stack
	s.AppendOpcodes(bscript.OpNIP)

	//push 178 to stack to split template after 'meta'
	//if err = s.AppendPushDataHexString("79"); err != nil {
	//	return nil, err
	//}

	s.AppendOpcodes(bscript.Op16)
	s.AppendOpcodes(bscript.Op12)
	s.AppendOpcodes(bscript.OpMUL)
	s.AppendOpcodes(bscript.Op14)
	s.AppendOpcodes(bscript.OpADD)

	s.AppendOpcodes(bscript.OpSPLIT)

	// grab length of pubkey and drop from stack
	s.AppendOpcodes(bscript.Op1)
	s.AppendOpcodes(bscript.OpSPLIT)
	s.AppendOpcodes(bscript.OpOVER)

	s.AppendOpcodes(bscript.OpSPLIT)
	// split first data prefix byte from txid
	s.AppendOpcodes(bscript.Op1)
	s.AppendOpcodes(bscript.OpSPLIT)

	// push 32 (0x20) to stack to split template after parent txid
	//if err = s.AppendPushDataHexString("20"); err != nil {
	//	return nil, err
	//}
	// accomplished with op_OVER
	s.AppendOpcodes(bscript.OpOVER)

	s.AppendOpcodes(bscript.OpSPLIT)

	// have txid here

	// push Op_5 to split off OP_DROP, OP_DROP, and first 3 bytes of P2PKH script
	s.AppendOpcodes(bscript.Op5)
	s.AppendOpcodes(bscript.OpSPLIT)

	// push 20 (0x14) to stack to split off pubkeyhash in P2PKH script
	if err = s.AppendPushDataHexString("14"); err != nil {
		return nil, err
	}
	s.AppendOpcodes(bscript.OpSPLIT)

	s.AppendOpcodes(bscript.Op2)
	s.AppendOpcodes(bscript.OpSPLIT)

	// drop anything after p2pkh
	s.AppendOpcodes(bscript.OpDROP)
	// nip pubkeyhash
	s.AppendOpcodes(bscript.OpNIP)
	// concatenate script template
	s.AppendOpcodes(bscript.OpCAT)
	s.AppendOpcodes(bscript.OpSWAP)

	// send parentTxId to altstack
	s.AppendOpcodes(bscript.OpTOALTSTACK)
	s.AppendOpcodes(bscript.OpCAT)
	s.AppendOpcodes(bscript.OpROT)
	s.AppendOpcodes(bscript.OpSWAP)
	s.AppendOpcodes(bscript.OpCAT)
	s.AppendOpcodes(bscript.OpROT)
	s.AppendOpcodes(bscript.OpSWAP)
	s.AppendOpcodes(bscript.OpCAT)
	// SHA1 hash the script template
	s.AppendOpcodes(bscript.OpSHA1)

	s.AppendOpcodes(bscript.Op3)
	s.AppendOpcodes(bscript.OpPICK)

	// move pubkey to bottom of stack and check hashes
	s.AppendOpcodes(bscript.OpEQUALVERIFY)

	// grab public key and check they are equal for valid metanet node
	s.AppendOpcodes(bscript.Op3)
	s.AppendOpcodes(bscript.OpPICK)
	s.AppendOpcodes(bscript.OpEQUALVERIFY)
	s.AppendOpcodes(bscript.OpNIP)
	s.AppendOpcodes(bscript.OpSWAP)

	return s, nil

}

func convertTxIdLittleEndian(s *bscript.Script) *bscript.Script {
	s.AppendOpcodes(bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP,
		bscript.OpSIZE,
		bscript.Op1SUB,
		bscript.OpSPLIT,
		bscript.OpSWAP)
	return s
}

func concatenateTxIdBigEndian(s *bscript.Script) *bscript.Script {
	s.AppendOpcodes(bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT,
		bscript.OpCAT)
	return s

}
