package availability

import (
	"bytes"
	"math/big"
	"strconv"
)

type BitSegment struct {
	big.Int
	start int
}

func (bs *BitSegment) SetUnit(position int, value byte) {
	bs.SetBit(&bs.Int, position, uint(value))
}

func NewBitSegment(start int) *BitSegment {
	return &BitSegment{
		start: start,
	}
}

func (bs *BitSegment) String() string {
	var buffer bytes.Buffer
	for i := 0; i < bs.BitLen(); i++ {
		buffer.WriteString(strconv.Itoa(int(bs.Bit(i))))
	}

	return buffer.String()
}
