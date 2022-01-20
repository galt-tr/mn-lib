package transaction

import (
	"errors"

	bsv "github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
)

type PayToMetanetAddress struct {
	Address    string `json:"address"`
	Satoshis   uint64 `json:"satoshis"`
	ParentTxId string `json:"parenttxid"`
}

func CreateSpendableMetanetTx(utxos []*bsv.Utxo, address *PayToMetanetAddress, opReturns []bsv.OpReturnData, privateKey *bsvec.PrivateKey) (*bt.Tx, error) {

	//Start creating a new transaction
	tx := bt.NewTx()

	//Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64

	// Loop all utxos and add to the transaction
	var err error
	for _, utxo := range utxos {
		if err = tx.From(utxo.TxID, utxo.Vout, utxo.ScriptPubKey, utxo.Satoshis); err != nil {
			return nil, err
		}
		totalSatoshis += utxo.Satoshis
	}

	var script *bscript.Script

	// Loop any pay addresses
	for _, address := range addresses {
		if script, err = NewMetanetP2PKH(address.Address, address.ParentTxId); err != nil {
			return nil, err
		}
		if err = tx.AddOutput(&bt.Output{
			Satoshis:      address.Satoshis,
			LockingScript: &script,
		}); err != nil {
			return nil, err
		}
	}

	// Loop any op returns
	var outPut *bt.Output
	for _, op := range opReturns {
		if outPut, err = bt.NewOpReturnPartsOutput(op); err != nil {
			return nil, err
		}
		tx.AddOutput(outPut)
	}

	//If inputs are supplied make sure they are sufficient for this transaction
	if len(tx.GetInputs()) > 0 {
		totalOutputSatoshis := tx.GetTotalOutputSatoshis() // Does not work properly
		if totalOutputSatoshis > totalSatoshis {
			return nil, errors.New("not enough in utxo(s) to cover")
		}
	}

	// Sign the transaction
	if privateKey != nil {
		signer := bt.InternalSigner{PrivateKey: privateKey, SigHashFlag: 0}
		if _, err = tx.SignAuto(&signer); err != nil {
			return nil, err
		}
	}

	// Return the transaction as a raw string
	return tx, nil

}
