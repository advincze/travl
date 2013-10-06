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

type BitSegmentPersistor interface {
	Save(s *BitSegment)
	Find(id string, start int) *BitSegment
	FindAll(id string) map[int]*BitSegment
}

func (s *BitSegment) String() string {
	var buffer bytes.Buffer

	for i := 0; i < s.BitLen(); i++ {
		buffer.WriteString(strconv.Itoa(int(s.Bit(i))))
	}

	return buffer.String()
}

func NewBitSegment(id string, start int) *BitSegment {
	return &BitSegment{
		Int:   *big.NewInt(0),
		start: start,
		ID:    id,
	}
}

type BitSegmentMemPersistor struct {
	bitAvSegments map[string]map[int]*BitSegment
}

func NewBitSegmentMemPersistor() *BitSegmentMemPersistor {
	return &BitSegmentMemPersistor{
		bitAvSegments: make(map[string]map[int]*BitSegment)}
}

func (bsmp *BitSegmentMemPersistor) Save(s *BitSegment) {
	if bsmp.bitAvSegments == nil {
		bsmp.bitAvSegments = make(map[string]map[int]*BitSegment)
	}
	segments, ok := bsmp.bitAvSegments[s.ID]
	if !ok {
		segments = make(map[int]*BitSegment)
		bsmp.bitAvSegments[s.ID] = segments
	}
	segments[s.start] = s

}

func (bsmp *BitSegmentMemPersistor) Find(id string, start int) *BitSegment {
	if bsmp.bitAvSegments == nil {
		return nil
	}
	segments, ok := bsmp.bitAvSegments[id]
	if !ok {
		return nil
	}
	return segments[start]
}

func (bsmp *BitSegmentMemPersistor) FindAll(id string) map[int]*BitSegment {
	if bsmp.bitAvSegments == nil {
		return nil
	}
	return bsmp.bitAvSegments[id]
}
