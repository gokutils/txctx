package txctx

import (
	"context"
	"errors"
)

var longTx = map[interface{}]*TxStore{}

/*
CreateLongTranscation create and store reference of transcation with given context
*/
func CreateLongTranscation(ctx context.Context, key interface{}) context.Context {
	txctx := NewTxStore(ctx)
	longTx[key] = txctx
	return context.WithValue(txctx.context, TransactionKeyStore, txctx)
}

/*
GetLongTransaction retrieve long transction by key
*/
func GetLongTransaction(key interface{}) (context.Context, error) {
	if v, ok := longTx[key]; ok {
		return context.WithValue(v.context, TransactionKeyStore, &TxStore{parent: v}), nil
	} else {
		return nil, errors.New("transaction.not_found")
	}
}

/*
GetLongTransaction retrieve long transction by key
*/
func DeleteLongTransaction(key interface{}) {
	delete(longTx, key)
}
