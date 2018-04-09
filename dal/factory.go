package dal

import (
	"fmt"
	"golang.org/x/net/context"
)

func defaultFactory(ctx context.Context, f *Factory) {
}

func (f *Factory) valid() error {
	return nil
}

func (f *Factory) Generate() error {
	if f.Sid == "" {
		f.Sid = fmt.Sprintf("%s-%s", f.Name, Randstr(8))
	}
	f.Token = fmt.Sprintf("%s-%s", f.Sid, Randstr(20))
	return nil
}
