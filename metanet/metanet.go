package metanet

import (
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"strings"

	bsv "github.com/bitcoinschema/go-bitcoin"
	"github.com/bitcoinsv/bsvd/bsvec"
	"github.com/galt-tr/mn-lib/transaction"
	"github.com/libsv/go-bt"
	"github.com/libsv/go-bt/bscript"
)

//Use Default "Meta" Prefix
//const prefix := "meta"

type MetanetNode struct {
	Prefix          string            //Prefix for Metanet Nodes - default to 'meta'
	NodeAddress     string            //Node Public Key Address
	ParentTxId      string            //Transaction ID of Parent Node
	Input           []*bsv.Utxo       //utxo of parent node
	InputPrivateKey *bsvec.PrivateKey //Private Key used to sign input
	Data            string            //Data to be added at end of OP_RETU in Metanet Node
	ChangeAddress   string            //Use for Change
}

func CreateSpendableNode(mn *MetanetNode) (string, string, error) {
	rawTx, preimage, err := createSpendableTransaction(mn)
	if err != nil {
		return "", "", err
	}
	return rawTx.ToString(), hex.EncodeToString(preimage), nil

}

func createSpendableTransaction(mn *MetanetNode) (*bt.Tx, []byte, error) {
	payTo := &transaction.PayToMetanetAddress{
		Address:       mn.NodeAddress,
		Satoshis:      900,
		ParentTxId:    mn.ParentTxId,
		ChangeAddress: mn.ChangeAddress,
	}

	rawTx, preimage, err := transaction.CreateSpendableMetanetTxWithChange(mn.Input, payTo, mn.Data, mn.ChangeAddress, nil, nil, mn.InputPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	return rawTx, preimage, nil

}

//Creates Metanet node and returns Raw Transaction Hex
func CreateNode(mn *MetanetNode) (string, error) {
	//Build OP_RETURN
	opReturn, err := buildOpReturn(mn)
	if err != nil {
		return "", err
	}
	rawTx, err := createTransaction(mn, opReturn)
	if err != nil {
		return "", err
	}
	return rawTx.ToString(), nil
}

func createTransaction(mn *MetanetNode, opReturn [][]byte) (*bt.Tx, error) {
	//TODO: Nodes should be spendable
	//payTo := &bsv.PayToAddress{
	//	Address:  mn.NodeAddress,
	//	Satoshis: 0,
	//}
	//Create Transaction
	rawTx, err := bsv.CreateTxWithChange(mn.Input, nil, []bsv.OpReturnData{opReturn}, mn.ChangeAddress, nil, nil, mn.InputPrivateKey)
	if err != nil {
		return nil, err
	}
	return rawTx, nil
}

func buildOpReturn(mn *MetanetNode) ([][]byte, error) {
	//Build OP_RETURN using Type from Go-Bitcoin
	opReturn := bsv.OpReturnData{[]byte(mn.Prefix), []byte(mn.NodeAddress), []byte(mn.ParentTxId), []byte(mn.Data)}
	return opReturn, nil
}

// Checks that Input Signature provided is valid for the Metanet Data Output
func IsChildNodeValid(rawTxHex string, vin, vout int) (bool, error) {
	transaction, err := bsv.TxFromHex(rawTxHex)
	if err != nil {
		return false, err
	}

	unlockingScript := transaction.Inputs[vin].UnlockingScript

	lockingScript := transaction.Outputs[vout].LockingScript

	if lockingScript.IsData() == false {
		return false, errors.New("output script is not data output")
	}
	lockingScriptArr, err := getScriptArray(lockingScript)
	if err != nil {
		return false, err
	}

	_, parentTxId, err := parseMetanetDataOutput(lockingScriptArr)
	if err != nil {
		return false, err
	}
	//fmt.Println(parentTxId)

	_, sigPubKey, err := getSignatureFromUnlockingScript(unlockingScript)
	sigAddress, err := bsv.GetAddressFromPubKeyString(sigPubKey, true)
	if err != nil {
		return false, err
	}

	parentTx, err := getParentTransaction(parentTxId)
	if err != nil {
		return false, err
	}

	parentLockingScripts, err := getLockingScriptsFromTransaction(parentTx)
	if len(parentLockingScripts) > 1 {
		//TODO: Handle Multiple Metanet Data Outputs. Currently assumes only One per Tx
		return false, errors.New("more than one metanet data output exists")
	}

	parentAddress, _, err := parseMetanetDataOutput(parentLockingScripts[0])
	if err != nil {
		return false, err
	}
	if parentAddress != sigAddress.String() {
		return false, errors.New("input signature does not match parent pubKey")
	}
	return true, nil
}

func getLockingScriptsFromTransaction(transaction *bt.Tx) ([][]string, error) {
	var lockingScripts [][]string
	for _, o := range transaction.Outputs {
		if o.LockingScript.IsData() {
			scriptArr, err := getScriptArray(o.LockingScript)
			if err != nil {
				return nil, err
			}
			check := validateMetanetDataOutput(scriptArr)
			if check == true {
				lockingScripts = append(lockingScripts, scriptArr)
			}
		}
	}
	return lockingScripts, nil
}
func getScriptArray(script *bscript.Script) ([]string, error) {
	scriptASM, err := script.ToASM()
	if err != nil {
		return nil, err
	}
	scriptArr := strings.Split(scriptASM, " ")
	if err != nil {
		return nil, err
	}
	return scriptArr, nil
}

func parseMetanetDataOutput(scriptArr []string) (addr, parTxId string, e error) {
	check := validateMetanetDataOutput(scriptArr)
	if check == false {
		return "", "", errors.New("data output does not begin with 'meta'")
	}
	address, err := hex.DecodeString(scriptArr[3])
	if err != nil {
		return "", "", err
	}
	parentTxId, err := hex.DecodeString(scriptArr[4])
	if err != nil {
		return "", "", err
	}
	return string(address), string(parentTxId), nil

}

func validateMetanetDataOutput(data []string) bool {
	if data[2] == "6d657461" {
		return true
	} else {
		return false
	}
}

func getSignatureFromUnlockingScript(script *bscript.Script) (signature, pubkey string, e error) {
	scriptASM, err := script.ToASM()
	if err != nil {
		return "", "", errors.New("unlocking script could not parse to ASM")
	}

	scriptSigArr := strings.Split(scriptASM, " ")
	if len(scriptSigArr) != 2 {
		return "", "", errors.New("scriptSig contains more than two objects")
	}
	return scriptSigArr[0], scriptSigArr[1], nil
}

func getParentTransaction(parentTxId string) (*bt.Tx, error) {
	resp, err := http.Get("https://api.whatsonchain.com/v1/bsv/main/tx/" + parentTxId + "/hex")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	parentTx, err := bsv.TxFromHex(string(body))
	if err != nil {
		return nil, err
	}
	return parentTx, nil
}
