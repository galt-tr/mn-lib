package transaction

import (
	"errors"
	"fmt"

	bsv "github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/galt-tr/mn-lib/script"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
	"github.com/libsv/go-bt/sighash"
)

type PayToMetanetAddress struct {
	Address        string `json:"address"`
	PublicKey      string `json:"publicKey"`
	Satoshis       uint64 `json:"satoshis"`
	ParentTxId     string `json:"parenttxid"`
	HasChange      bool   `json:"hasChange"`
	ChangeAddress  string `json:"changeAddress"`
	ChangeSatoshis uint64 `json"changeSatoshis"`
}

func CreateSpendableMetanetTxWithChange(utxos []*bsv.Utxo, mnAddress *PayToMetanetAddress, data, changeAddress string, standardRate, dataRate *bt.Fee, privateKey *bsvec.PrivateKey) (*bt.Tx, error) {

	if len(utxos) == 0 {
		return nil, errors.New("utxos(s) are required to create a tx")
	} else if len(changeAddress) == 0 {
		return nil, errors.New("change address is required")
	}

	// Accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64
	var totalPayToSatoshis uint64
	var remainder uint64

	// Loop utxos and get total usable satoshis
	for _, utxo := range utxos {
		totalSatoshis += utxo.Satoshis
	}

	// Loop all payout address amounts
	totalPayToSatoshis += mnAddress.Satoshis

	// Sanity check - already not enough satoshis?
	if totalPayToSatoshis > totalSatoshis {
		return nil, fmt.Errorf(
			"not enough in utxo(s) to cover: %d + (fee), total found: %d",
			totalPayToSatoshis,
			totalSatoshis,
		)
	}
	// THIS IS ALL FUCKED UP PLS CHECK https://github.com/BitcoinSchema/go-bitcoin/blob/master/transaction.go#L45
	// JUST GOING TO USE THIS FOR NOW

	// Add the change address as the difference (all change except 1 sat for Draft tx)
	// Only if the tx is NOT for the full amount

	if totalPayToSatoshis != totalSatoshis {
		mnAddress.HasChange = true
		mnAddress.ChangeSatoshis = (totalSatoshis - (totalPayToSatoshis + 1))
	}

	// Create the "Draft tx"
	fee, err := draftTx(utxos, mnAddress, data, privateKey, standardRate, dataRate)
	if err != nil {
		return nil, err
	}

	// Check that we have enough to cover the fee
	if (totalPayToSatoshis + fee) > totalSatoshis {
		mnAddress.HasChange = false

		//Re-run draft tx with no change address
		if fee, err = draftTx(utxos, mnAddress, data, privateKey, standardRate, dataRate); err != nil {
			return nil, err
		}
		mnAddress.HasChange = true

		// Get the remainder missing (handle negative overflow safer)
		totalToPay := totalPayToSatoshis + fee
		if totalToPay >= totalSatoshis {
			remainder = totalToPay - totalSatoshis
		} else {
			remainder = totalSatoshis - totalToPay
		}

		// Remove remainder from last used payToAddress (or continue until found)
		mnAddress.ChangeSatoshis = mnAddress.ChangeSatoshis - remainder

	} else {

		// Remove the change address (old version with original satoshis)
		// Add the change address as the difference (now with adjusted fee)
		mnAddress.ChangeSatoshis = totalSatoshis - (totalPayToSatoshis + fee)
	}
	//Create the final tx or error
	return CreateSpendableMetanetTx(utxos, mnAddress, data, privateKey)

}

// draftTx is a helper method to create a draft tx and associated fees
func draftTx(utxos []*bsv.Utxo, payTo *PayToMetanetAddress, data string, privateKey *bsvec.PrivateKey, standardRate, dataRate *bt.Fee) (uint64, error) {

	// Create the "Draft tx"
	tx, err := CreateSpendableMetanetTx(utxos, payTo, data, privateKey)
	if err != nil {
		return 0, err
	}

	// Calculate the fees for the "Draft tx"
	// todo: hack to add 1 extra sat - ensuring that fee is over the minimum with rounding issues in WOC and other systems
	fee := bsv.CalculateFeeForTx(tx, standardRate, dataRate) + 1
	return fee, nil
}

// Creates Transaction without Change
func CreateSpendableMetanetTx(utxos []*bsv.Utxo, mnAddress *PayToMetanetAddress, data string, privateKey *bsvec.PrivateKey) (*bt.Tx, error) {

	//start creating a new transaction
	tx := bt.NewTx()

	//accumulate the total satoshis from all utxo(s)
	var totalSatoshis uint64

	// loop all utxos and add to the transaction
	var err error
	for _, utxo := range utxos {
		if err = tx.From(utxo.TxID, utxo.Vout, utxo.ScriptPubKey, utxo.Satoshis); err != nil {
			return nil, err
		}
		totalSatoshis += utxo.Satoshis
	}

	var s *bscript.Script

	if s, err = script.NewMetanetP2PKH(mnAddress.PublicKey, mnAddress.ParentTxId, data); err != nil {
		return nil, err
	}

	tx.AddOutput(&bt.Output{
		Satoshis:      mnAddress.Satoshis,
		LockingScript: s,
	})

	// Add Change
	if mnAddress.HasChange == true {
		changeScript, err := bscript.NewP2PKHFromAddress(mnAddress.ChangeAddress)
		if err != nil {
			return nil, err
		}

		tx.AddOutput(&bt.Output{
			Satoshis:      mnAddress.ChangeSatoshis,
			LockingScript: changeScript,
		})
	}

	//if inputs are supplied make sure they are sufficient for this transaction
	if len(tx.GetInputs()) > 0 {
		totalOutputSatoshis := tx.GetTotalOutputSatoshis() // does not work properly
		if totalOutputSatoshis > totalSatoshis {
			return nil, errors.New("not enough in utxo(s) to cover")
		}
	}

	//Set Sighash Flag
	var sigHashFlag sighash.Flag
	sigHashFlag = 0x41

	var preimage []byte
	// Calculate Preimage
	// TODO: Currently defaults to signing for output 0 - should do smart checking
	if preimage, err = tx.GetInputPreimage(0, sigHashFlag); err != nil {
		return nil, err
	}

	//preimage = nil

	//TODO: Only signing output 0
	var index uint32
	index = 0

	// sign the transaction - currently only signing input 0
	// instead of implementing our own signer will continue to use go-bt v1 signer and append to unlocking script
	if privateKey != nil {
		signer := bt.InternalSigner{PrivateKey: privateKey, SigHashFlag: sigHashFlag}
		if tx, err = signer.Sign(index, tx); err != nil {
			return nil, err
		}
	}

	// Append to P2PKH Unlocking Script
	s = tx.Inputs[0].UnlockingScript
	if s, err = script.AppendMetanetUnlockingScript(s, preimage); err != nil {
		return nil, err
	}

	// return the transaction as a raw string
	return tx, nil

}
