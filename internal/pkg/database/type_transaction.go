package database

import (
	"context"
)

type TypeTransactionHandler[T any] struct {
	handler *TransactionHandlerDB
}

func NewTypedTransaction[T any](handler *TransactionHandlerDB) *TypeTransactionHandler[T] {
	return &TypeTransactionHandler[T]{handler: handler}
}

func (t *TypeTransactionHandler[T]) WithCtx(ctx context.Context, fn func(ctx context.Context) (T, error)) (result T, e error) {
	ctx, tx, e := t.handler.Create(ctx, nil)
	if e != nil {
		return result, e
	}

	result, funcErr := fn(ctx)

	if e = HandleTransaction(tx, funcErr); e != nil {
		return result, e
	}

	return result, nil
}
