package dal

import (
	"golang.org/x/net/context"
)

func defaultModule(ctx context.Context, m *Module) {
}

func (m *Module) valid() error {
	return nil
}
