package sign

import (
	"github.com/bitcoinsv/bsvd/bsvec"
	mnScript "github.com/galt-tr/mn-lib/script"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
	"github.com/libsv/go-bt/sighash"
)

type Signer struct {
	PrivateKey  *bsvec.PrivateKey
	SigHashFlag sighash.Flag
}

func (s *Signer) SignMetanetTransaction(index uint32, preimage []byte, unsignedTx *bt.Tx) (signedTx *bt.Tx, err error) {
	if s.SigHashFlag == 0 {
		s.SigHashFlag = sighash.AllForkID
	}

	var sh []byte
	if sh, err = unsignedTx.GetInputSignatureHash(index, s.SigHashFlag); err != nil {
		return nil, err
	}

	var sig *bsvec.Signature
	if sig, err = s.PrivateKey.Sign(bt.ReverseBytes(sh)); err != nil {
		return nil, err
	}

	var script *bscript.Script
	if script, err = mnScript.NewMetanetUnlockingScript(s.PrivateKey.PubKey().SerializeCompressed(), preimage, sig.Serialize(), s.SigHashFlag); err != nil {
		return nil, err
	}

	if err = unsignedTx.ApplyUnlockingScript(index, script); err != nil {
		return nil, err
	}

	signedTx = unsignedTx

	return signedTx, err
}

func (s *Signer) SignOpPushTxTransaction(index uint32, preimage []byte, unsignedTx *bt.Tx) (signedTx *bt.Tx, err error) {
	if s.SigHashFlag == 0 {
		s.SigHashFlag = sighash.AllForkID
	}
	var sh []byte
	if sh, err = unsignedTx.GetInputSignatureHash(index, s.SigHashFlag); err != nil {
		return nil, err
	}

	var sig *bsvec.Signature
	if sig, err = s.PrivateKey.Sign(bt.ReverseBytes(sh)); err != nil {
		return nil, err
	}

	var script *bscript.Script

	if script, err = mnScript.NewOpPushTxUnlockingScript(s.PrivateKey.PubKey().SerializeCompressed(), preimage, sig.Serialize(), s.SigHashFlag); err != nil {
		return nil, err
	}

	if err = unsignedTx.ApplyUnlockingScript(index, script); err != nil {
		return nil, err
	}

	signedTx = unsignedTx

	return signedTx, err

}
