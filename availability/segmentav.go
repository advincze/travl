package availability

import (
	"bytes"
	"strconv"
	"time"
)

type SegmentAv struct {
	ID            string
	internalRes   TimeResolution
	segmentSize   int
	bitAvSegments map[string]map[int]*BitSegment
}

func NewSegmentAv(id string, res TimeResolution) *SegmentAv {
	return &SegmentAv{
		ID:            id,
		internalRes:   res,
		segmentSize:   int(Day / res),
		bitAvSegments: make(map[string]map[int]*BitSegment),
	}
}

func (av *SegmentAv) size() int {
	var size int
	segments := av.bitAvSegments[av.ID]
	for _, segment := range segments {
		size += len(segment.Bytes())
	}
	return size
}

func (av *SegmentAv) Set(from, to time.Time, value byte) {
	fromUnit := TimeToUnitFloor(from, av.internalRes)
	toUnit := TimeToUnitFloor(to, av.internalRes)
	av.setUnitInternal(fromUnit, toUnit, value)
}

func (av *SegmentAv) Get(from, to time.Time, res TimeResolution) *BitVector {

	if res > av.internalRes {
		// lower resolution
		fromUnit := TimeToUnitFloor(FloorDate(from, res), av.internalRes)
		toUnit := TimeToUnitFloor(CeilDate(to, res), av.internalRes)
		arr := av.getUnitInternal(fromUnit, toUnit)
		factor := int(res / av.internalRes)
		reducedArr := reduceByFactor(arr, factor, reduceAllOne)
		return NewBitVector(res, av.internalRes, reducedArr, FloorDate(from, res))

	} else if res < av.internalRes {
		// higher resolution
		fromUnitInternalRes := TimeToUnitFloor(from, av.internalRes)
		toUnitInternalRes := TimeToUnitFloor(CeilDate(to, av.internalRes), av.internalRes)
		arr := av.getUnitInternal(fromUnitInternalRes, toUnitInternalRes)
		factor := int(av.internalRes / res)
		arrMultiplied := multiplyByFactor(arr, factor)
		cutoff := TimeToUnitFloor(from, res) - fromUnitInternalRes*factor
		origlen := TimeToUnitFloor(to, res) - TimeToUnitFloor(from, res)
		arrTrimmed := arrMultiplied[cutoff : cutoff+origlen]
		return NewBitVector(res, av.internalRes, arrTrimmed, FloorDate(from, av.internalRes))
	} else {
		// internal resolution
		fromUnit := TimeToUnitFloor(from, res)
		toUnit := TimeToUnitFloor(to, res)
		arr := av.getUnitInternal(fromUnit, toUnit)
		return NewBitVector(res, av.internalRes, arr, FloorDate(from, res))
	}
}

func printarr(arr []byte) string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	buffer.WriteString(strconv.Itoa(len(arr)))
	buffer.WriteString("-")
	for i := 0; i < len(arr); i++ {
		if arr[i] == 1 {
			buffer.WriteString("1, ")
		} else {
			buffer.WriteString("0, ")
		}
	}
	buffer.WriteString("],")

	return buffer.String()
}

func multiplyByFactor(data []byte, factor int) []byte {
	length := len(data) * factor
	var multipliedData []byte = make([]byte, length)
	j := 0
	for _, b := range data {
		for i := 0; i < factor; i++ {
			multipliedData[j] = b
			j++
		}
	}
	return multipliedData
}

func reduceByFactor(data []byte, factor int, reduceFn func([]byte) byte) []byte {
	length := len(data) / factor
	var reducedData []byte = make([]byte, length)
	for i, j := 0, 0; i < length; i++ {
		reducedData[i] = reduceFn(data[j : j+factor])
		j += factor
	}
	return reducedData
}

func reduceAllOne(data []byte) byte {
	for _, b := range data {
		if b != 1 {
			return 0
		}
	}
	return 1
}

func reduceAnyOne(data []byte) byte {
	for _, b := range data {
		if b == 1 {
			return 1
		}
	}
	return 0
}

func reduceMajority(data []byte) byte {
	sizewin := len(data) / 2
	count := 0
	for _, b := range data {
		if b == 1 {
			count++
		}
	}
	if count > sizewin {
		return 1
	}
	return 0
}

func (av *SegmentAv) SetAt(at time.Time, value byte) {
	atUnit := TimeToUnitFloor(at, av.internalRes)
	av.setUnitInternal(atUnit, atUnit+1, value)
}

func (av *SegmentAv) GetAt(at time.Time) byte {
	atUnit := TimeToUnitFloor(at, av.internalRes)
	arr := av.getUnitInternal(atUnit, atUnit+1)
	return byte(arr[0])
}

func (av *SegmentAv) String() string {
	var buffer bytes.Buffer

	segments := av.bitAvSegments[av.ID]
	for _, segment := range segments {
		buffer.WriteString(strconv.Itoa(segment.start))
		buffer.WriteString("->")
		buffer.WriteString(segment.String())
		buffer.WriteRune('\n')
	}

	return buffer.String()
}

func (av *SegmentAv) segmentStart(i int) int {
	return i - i%av.segmentSize
}

func (av *SegmentAv) getOrEmptyBitSegment(startValue int) *BitSegment {
	if segment := av.bitAvSegments[av.ID][startValue]; segment != nil {
		return segment
	}
	return NewBitSegment(av.ID, startValue)
}

func (av *SegmentAv) setUnitInternal(from, to int, value byte) {
	currentBitSegment := av.getOrEmptyBitSegment(av.segmentStart(from))
	for i, j := from, from%av.segmentSize; i < to; i, j = i+1, j+1 {
		if j == av.segmentSize {
			av.Savex(currentBitSegment)
			currentBitSegment = av.getOrEmptyBitSegment(i)
			j = 0
		}
		currentBitSegment.SetBit(&currentBitSegment.Int, j, uint(value))
	}
	av.Savex(currentBitSegment)
}

func (av *SegmentAv) Savex(s *BitSegment) {
	segments, ok := av.bitAvSegments[s.ID]
	if !ok {
		segments = make(map[int]*BitSegment)
		av.bitAvSegments[s.ID] = segments
	}
	segments[s.start] = s

}

func (av *SegmentAv) getUnitInternal(from, to int) []byte {
	length := to - from
	result := make([]byte, length)
	currentBitSegment := av.getOrEmptyBitSegment(av.segmentStart(from))
	for i, j := 0, from%av.segmentSize; i < length; i, j = i+1, j+1 {
		if j == av.segmentSize {
			currentBitSegment = av.getOrEmptyBitSegment(i + from)
			j = 0
		}
		result[i] = byte(currentBitSegment.Bit(j))
	}
	return result
}
