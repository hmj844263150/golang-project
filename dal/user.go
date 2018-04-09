package dal

import (
	"golang.org/x/net/context"
)

func defaultUser(ctx context.Context, u *User) {
}

func (u *User) valid() error {
	return nil
}
