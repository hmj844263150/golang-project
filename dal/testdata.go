package dal

import (
	"golang.org/x/net/context"
)

func defaultTestdata(ctx context.Context, t *Testdata) {
	t.Latest = true
}

func (t *Testdata) valid() error {
	return nil
}
