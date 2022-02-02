package metanet

import (
	"github.com/galt-tr/mn-lib/transaction"
	"github.com/galt-tr/mn-lib/types"
)

func NewMetanetNode(mn *types.MetanetNode) (string, error) {
	rawTx, err := transaction.CreateMetanetTransaction(mn)
	if err != nil {
		return "", err
	}
	return rawTx, nil

}
