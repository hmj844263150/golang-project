package dal

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"golang.org/x/net/context"
	"strings"
)

func defaultBatch(ctx context.Context, b *Batch) {
	b.Sid = Randstr(10)
}

func (b *Batch) valid() error {
	if b.Sid == "" || b.FactorySid == "" {
		return errors.New("b.Sid is empty or b.FactorySid is empty")
	}
	if b.IsCus {
		if len(b.CusMacFrom) != 12 || len(b.CusMacTo) != 12 {
			return errors.New("bssid length must be 12")
		}
		if (b.CusMacNumTo - b.CusMacNumFrom + 1) != b.Cnt {
			return errors.New("cnt not match")
		}
	}
	if b.Cnt <= 0 || b.Remain <= 0 {
		return errors.New("b.Cnt <= 0 || b.Remain <= 0")
	}
	if b.CusMacNumFrom > b.CusMacNumTo {
		return errors.New("b.CusMacFrom > b.CusMacNumTo")
	}
	return nil
}

func (b *Batch) Adjust() error {
	b.EspMacFrom = normal(b.EspMacFrom)
	b.EspMacTo = normal(b.EspMacTo)
	b.CusMacFrom = normal(b.CusMacFrom)
	b.CusMacTo = normal(b.CusMacTo)
	b.EspMacNumFrom = num(b.EspMacFrom)
	b.EspMacNumTo = num(b.EspMacTo)
	b.CusMacNumFrom = num(b.CusMacFrom)
	b.CusMacNumTo = num(b.CusMacTo)
	if b.CusMacFrom != "" && b.CusMacTo != "" {
		b.IsCus = true
	}
	if b.IsCus {
		b.Cnt = b.CusMacNumTo - b.CusMacNumFrom + 1
	}
	b.Remain = b.Cnt
	return b.Valid()
}

func (b *Batch) NextRemainMac() (int, string, error) {
	if b.Remain == 0 {
		return 0, "", errors.New("there are not space left")
	}
	if !b.IsCus {
		b.Remain = b.Remain - 1
		b.Update(BatchCol.Remain)
		return b.Cnt - b.Remain, "", nil
	}
	current := b.CusMacNumFrom + b.Cnt - b.Remain
	b.Remain = b.Remain - 1
	b.Update(BatchCol.Remain)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(current))
	return b.Cnt - b.Remain, b.CusMacFrom[0:6] + hex.EncodeToString(bs)[2:], nil
}

func normal(mac string) string {
	mac = strings.ToLower(mac)
	mac = strings.Replace(mac, ":", "", -1)
	if len(mac) != 12 {
		return ""
	}
	return mac
}

func num(mac string) int {
	if len(mac) != 12 {
		return 0
	}
	bs := make([]byte, 4)
	macBytes := []byte(mac)
	hex.Decode(bs[1:4], macBytes[6:12])
	return int(binary.BigEndian.Uint32(bs))
}
