/*
TxCtx permet de gerer les transaction au travers des context

ceci peremt de gerer par scope sur plusieur base de donnée les transactions
*/
package txctx

import (
	"context"
	"errors"
	"sync"
)

type ContextTransactionKey string

const (
	TransactionKeyStore ContextTransactionKey = "txctx-transaction-store"
)

type Tx interface {
	Commit(_ context.Context) error
	Rollback(_ context.Context) error
}

type Txs interface {
	Tx
	Add(txs ...Tx)
	GetValue(key interface{}) (interface{}, bool)
	GetValueOrStore(key interface{}, value interface{}) (interface{}, bool)
	SetValue(key interface{}, value interface{})
}

type TxStore struct {
	parent  Txs
	context context.Context
	txs     []Tx

	values sync.Map
}

func NewTxStore(ctx context.Context) *TxStore {
	return &TxStore{txs: []Tx{}, context: ctx, values: sync.Map{}}
}

func (impl *TxStore) Commit(ctx context.Context) error {
	if impl.parent != nil {
		return nil
	}
	for _, tx := range impl.txs {
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	impl.txs = []Tx{}
	return nil
}

func (impl *TxStore) Rollback(ctx context.Context) error {
	if impl.parent != nil {
		return nil
	}
	errs := []error{}
	for _, tx := range impl.txs {
		if err := tx.Rollback(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	impl.txs = []Tx{}
	return nil
}

func (impl *TxStore) Add(txs ...Tx) {
	if impl.parent != nil {
		impl.parent.Add(txs...)
		return
	}
	impl.txs = append(impl.txs, txs...)
}

func (impl *TxStore) GetValue(key interface{}) (interface{}, bool) {
	if impl.parent != nil {
		return impl.parent.GetValue(key)
	}
	return impl.values.Load(key)
}

func (impl *TxStore) GetValueOrStore(key interface{}, value interface{}) (interface{}, bool) {
	if impl.parent != nil {
		return impl.parent.GetValueOrStore(key, value)
	}
	return impl.values.LoadOrStore(key, value)
}

func (impl *TxStore) SetValue(key interface{}, value interface{}) {
	if impl.parent != nil {
		impl.parent.SetValue(key, value)
		return
	}
	impl.values.Store(key, value)
}

/*
GetValue retrieve Key, value
*/
func GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		return v.GetValue(key)
	}
	return nil, false
}

/*
GetValueOrStore retrieve Key, value or store in tx store
*/
func GetValueOrStore(ctx context.Context, key interface{}, value interface{}) (interface{}, bool) {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		return v.GetValueOrStore(key, value)
	}
	return nil, false
}

/*
SetValue store Key, value in tx store
*/
func SetValue(ctx context.Context, key interface{}, value interface{}) {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		v.SetValue(key, value)
		return
	}
}

/*
IsTxContext check if context containe key for tx store
*/
func IsTxContext(ctx context.Context) bool {
	_, ok := ctx.Value(TransactionKeyStore).(*TxStore)
	return ok
}

/*
Start new context
*/
func Begin(ctx context.Context) context.Context {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		return context.WithValue(ctx, TransactionKeyStore, &TxStore{parent: v})
	} else {

		return context.WithValue(ctx, TransactionKeyStore, NewTxStore(ctx))
	}
}

/*
Add add tx in context
*/
func Add(ctx context.Context, txs ...Tx) {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		v.Add(txs...)
	}
}

/*
Commit transaction in context

Executed if is a first context
*/
func Commit(ctx context.Context) error {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		return v.Commit(ctx)
	}
	return nil
}

/*
Rollback transaction in context

Executed if is a first context
*/
func Rollback(ctx context.Context) error {
	if v, ok := ctx.Value(TransactionKeyStore).(*TxStore); ok {
		return v.Rollback(ctx)
	}
	return nil
}
