package dal

import (
	"golang.org/x/net/context"
)

func defaultTestlog(ctx context.Context, t *Testlog) {
}

func (t *Testlog) valid() error {
	return nil
}
