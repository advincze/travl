package availability

import (
	"bytes"
	"math/big"
	"strconv"
)

type BitSegment struct {
	big.Int
	start int
	ID    string
}

func NewBitSegment(id string, start int) *BitSegment {
	return &BitSegment{
		Int:   *big.NewInt(0),
		start: start,
		ID:    id,
	}
}

func (s *BitSegment) String() string {
	var buffer bytes.Buffer

	for i := 0; i < s.BitLen(); i++ {
		buffer.WriteString(strconv.Itoa(int(s.Bit(i))))
	}

	return buffer.String()
}
