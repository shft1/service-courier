package service

import "context"

type TxManagerDo interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
