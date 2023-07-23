package txctx

import "context"

type ResolveCallback struct {
	onCommit   func(ctx context.Context) error
	onRollback func(ctx context.Context) error
}

func (r *ResolveCallback) Commit(ctx context.Context) error {
	if r.onCommit != nil {
		return r.onCommit(ctx)
	}
	return nil
}

func (r *ResolveCallback) Rollback(ctx context.Context) error {
	if r.onRollback != nil {
		return r.onRollback(ctx)
	}
	return nil
}

func OnCommit(ctx context.Context, callback func(ctx context.Context) error) {
	Add(ctx, &ResolveCallback{onCommit: callback})
}

func OnRollback(ctx context.Context, callback func(ctx context.Context) error) {
	Add(ctx, &ResolveCallback{onRollback: callback})
}
