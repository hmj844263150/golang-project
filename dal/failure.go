package dal

import (
	"golang.org/x/net/context"
)

func defaultFailure(ctx context.Context, f *Failure) {
}

func (f *Failure) valid() error {
	return nil
}
