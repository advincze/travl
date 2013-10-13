package availability

import (
	"bytes"
	"strconv"
)

type SegmentedVector struct {
	segmentLength int
	segments      map[int]*BitSegment
}

func NewSegmentedVector(segmentLength int) *SegmentedVector {
	return &SegmentedVector{
		segmentLength: segmentLength,
		segments:      make(map[int]*BitSegment),
	}
}

func (sv *SegmentedVector) setUnit(from, to int, value byte) {
	segmentStart := sv.segmentStart(from)
	segment := sv.getOrEmptyBitSegment(segmentStart)

	for i, j := from, from%sv.segmentLength; i < to; i, j = i+1, j+1 {
		if j == sv.segmentLength {
			sv.storeSegment(segment)
			segment = sv.getOrEmptyBitSegment(i)
			j = 0
		}
		segment.SetUnit(j, value)
	}
	sv.segments[segment.start] = segment
}

func (sv *SegmentedVector) storeSegment(segment *BitSegment) {
	sv.segments[segment.start] = segment
}

func (sv *SegmentedVector) segmentStart(i int) int {
	return i - i%sv.segmentLength
}

func (sv *SegmentedVector) getUnit(from, to int) []byte {
	length := to - from
	result := make([]byte, length)
	currentBitSegment := sv.getOrEmptyBitSegment(sv.segmentStart(from))
	for i, j := 0, from%sv.segmentLength; i < length; i, j = i+1, j+1 {
		if j == sv.segmentLength {
			currentBitSegment = sv.getOrEmptyBitSegment(i + from)
			j = 0
		}
		result[i] = byte(currentBitSegment.Bit(j))
	}
	return result
}

func (sv *SegmentedVector) getOrEmptyBitSegment(startValue int) *BitSegment {
	if segment := sv.segments[startValue]; segment != nil {
		return segment
	}
	return NewBitSegment(startValue)
}

func (sv *SegmentedVector) SizeInBytes() int {
	var sizeInBytes int
	for _, segment := range sv.segments {
		sizeInBytes += len(segment.Bytes())
	}
	return sizeInBytes
}

func (sv *SegmentedVector) String() string {
	var buffer bytes.Buffer
	for _, segment := range sv.segments {
		buffer.WriteString(strconv.Itoa(segment.start))
		buffer.WriteString("->")
		buffer.WriteString(segment.String())
		buffer.WriteRune('\n')
	}
	return buffer.String()
}
