package transaction

import (
	"github.com/galt-tr/mn-lib/script"
	"github.com/galt-tr/mn-lib/types"
	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript"
	pushtxpreimage "github.com/murray-distributed-technologies/go-pushtx/preimage"
)

func CreateNode(mn *types.MetanetNodes) (string, error) {
	var err error
	tx := bt.NewTx()

	err = tx.From(mn.ParentTxId,
		mn.Input.PreviousTxOutIndex,
		mn.Input.PreviousTxScript.String(),
		mn.Input.PreviousTxSatoshis)
	if err != nil {
		return "", err
	}
	// add metanet output
	if tx, err = AddMetanetOutput(tx, mn); err != nil {
		return "", err
	}

	// add change output
	// TODO: if input is nonstandard change doesn't work
	if input.PreviousTxScript.IsP2PKH() {
		fq := bt.NewFeeQuote()
		if err = tx.ChangeToAddress(mn.ChangeAddress, fq); err != nil {
			return "", err
		}
	}
	if !input.PreviousTxScript.IsP2PKH() {
		lockingScript, err := bscript.NewP2PKHFromAddress(mn.ChangeAddress)
		if err != nil {
			return "", err
		}
		amount := (mn.Input.PreviousTxSatoshis - satoshis - 500)
		changeOutput := bt.Output{
			Satoshis:      amount,
			LockingScript: lockingScript,
		}
		tx.AddOutput(&changeOutput)
	}

	unlocker := Getter{PrivateKey: mn.InputPrivateKey}

	//sign Input
	if err = tx.FillAllInputs(context.Background(), &unlocker); err != nil {
		return "", err
	}
	return tx.String(), nil

}

func AddMetanetOutput(tx *bt.Tx, mn *types.MetanetNode) (*bt.Tx, error) {
	lockingScript, err := script.NewMetanetLockingScript(mn)
	if err != nil {
		return nil, err
	}

	output := bt.Output{
		Satoshis:      mn.Satoshis,
		LockingScript: lockingScript,
	}
	tx.AddOutput(&output)
	return tx, nil
}

type Getter struct {
	PrivateKey *bec.PrivateKey
}

func (g *Getter) Unlocker(ctx context.Context, lockingScript *bscript.Script) (bt.Unlocker, error) {

	if lockingScript.IsP2PKH() {
		return &btunlocker.Simple{PrivateKey: g.PrivateKey}, nil
	}
	if script.IsMetanetTx(lockingScript) {
		return &UnlockMetanetTx{PrivateKey: g.PrivateKey}, nil
	}
	return nil, errors.New("locking script is not P2PKH or Metanet")
}

type UnlockMetanetTx struct {
	PrivateKey *bec.PrivateKey
}

func (u *UnlockMetanetTx) UnlockingScript(ctx context.Context, tx *bt.Tx, params bt.UnlockerParams) (*bscript.Script, error) {
	if params.SigHashFlags == 0 {
		params.SigHashFlags = sighash.AllForkID
	}
	preimage, err := tx.CalcInputPreimage(params.InputIdx, params.SigHashFlags)
	if err != nil {
		return nil, err
	}
	preimage, nLockTime, err := pushtxpreimage.CheckForLowS(preimage)
	if err != nil {
		return nil, err
	}
	tx.LockTime = nLockTime

	// defaultHex is used to fix a bug in the original client (see if statement in the CalcInputSignatureHash func)
	var defaultHex = []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var sh []byte
	sh = crypto.Sha256d(preimage)

	if bytes.Equal(defaultHex, preimage) {
		sh = preimage
	}

	sig, err := u.PrivateKey.Sign(sh)
	if err != nil {
		return nil, err
	}

	pubKey := u.PrivateKey.PubKey().SerialiseCompressed()
	signature := sig.Serialise()

	uscript, err := script.NewMetanetUnlockingScript(pubKey, preimage, signature, params.SigHashFlags)
	if err != nil {
		return nil, err
	}

	return uscript, nil

}
