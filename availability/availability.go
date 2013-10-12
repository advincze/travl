package availability

import (
	"bytes"
	"strconv"
	"time"
)

type Availability struct {
	ID          string
	internalRes TimeResolution
	segmentSize int
	Segments    map[int]*BitSegment
}

func NewAvailability(id string, res TimeResolution) *Availability {
	return &Availability{
		ID:          id,
		internalRes: res,
		segmentSize: int(Day / res),
		Segments:    make(map[int]*BitSegment),
	}
}

func (av *Availability) sizeInBytes() int {
	var sizeInBytes int
	for _, segment := range av.Segments {
		sizeInBytes += len(segment.Bytes())
	}
	return sizeInBytes
}

func (av *Availability) Set(from, to time.Time, value byte) {
	fromUnit := TimeToUnitFloor(from, av.internalRes)
	toUnit := TimeToUnitFloor(to, av.internalRes)
	av.setUnitInternal(fromUnit, toUnit, value)
}

func (av *Availability) Get(from, to time.Time, res TimeResolution) *AvailabilityResult {
	if res > av.internalRes {
		return av.getWithLowerResolution(from, to, res)
	} else if res < av.internalRes {
		return av.getWithHigherResolution(from, to, res)
	}
	return av.getWithInternalResolution(from, to, res)
}

func (av *Availability) getWithLowerResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	fromUnit := TimeToUnitFloor(RoundDown(from, res), av.internalRes)
	toUnit := TimeToUnitFloor(RoundUp(to, res), av.internalRes)
	arr := av.getUnitInternal(fromUnit, toUnit)
	factor := int(res / av.internalRes)
	reducedArr := reduceByFactor(arr, factor, reduceAllOne)
	return NewAvailabilityResult(res, av.internalRes, reducedArr, RoundDown(from, res))
}

func (av *Availability) getWithHigherResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	// higher resolution
	fromUnitInternalRes := TimeToUnitFloor(from, av.internalRes)
	toUnitInternalRes := TimeToUnitFloor(RoundUp(to, av.internalRes), av.internalRes)
	arr := av.getUnitInternal(fromUnitInternalRes, toUnitInternalRes)
	factor := int(av.internalRes / res)
	arrMultiplied := multiplyByFactor(arr, factor)
	cutoff := TimeToUnitFloor(from, res) - fromUnitInternalRes*factor
	origlen := TimeToUnitFloor(to, res) - TimeToUnitFloor(from, res)
	arrTrimmed := arrMultiplied[cutoff : cutoff+origlen]
	return NewAvailabilityResult(res, av.internalRes, arrTrimmed, RoundDown(from, av.internalRes))
}

func (av *Availability) getWithInternalResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	fromUnit := TimeToUnitFloor(from, res)
	toUnit := TimeToUnitFloor(to, res)
	arr := av.getUnitInternal(fromUnit, toUnit)
	return NewAvailabilityResult(res, av.internalRes, arr, RoundDown(from, res))
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

func (av *Availability) SetAt(at time.Time, value byte) {
	atUnit := TimeToUnitFloor(at, av.internalRes)
	av.setUnitInternal(atUnit, atUnit+1, value)
}

func (av *Availability) GetAt(at time.Time) byte {
	atUnit := TimeToUnitFloor(at, av.internalRes)
	arr := av.getUnitInternal(atUnit, atUnit+1)
	return byte(arr[0])
}

func (av *Availability) String() string {
	var buffer bytes.Buffer
	for _, segment := range av.Segments {
		buffer.WriteString(strconv.Itoa(segment.start))
		buffer.WriteString("->")
		buffer.WriteString(segment.String())
		buffer.WriteRune('\n')
	}

	return buffer.String()
}

func (av *Availability) segmentStart(i int) int {
	return i - i%av.segmentSize
}

func (av *Availability) getOrEmptyBitSegment(startValue int) *BitSegment {
	if segment := av.Segments[startValue]; segment != nil {
		return segment
	}
	return NewBitSegment(av.ID, startValue)
}

func (av *Availability) setUnitInternal(from, to int, value byte) {
	currentBitSegment := av.getOrEmptyBitSegment(av.segmentStart(from))
	for i, j := from, from%av.segmentSize; i < to; i, j = i+1, j+1 {
		if j == av.segmentSize {
			av.Segments[currentBitSegment.start] = currentBitSegment

			currentBitSegment = av.getOrEmptyBitSegment(i)
			j = 0
		}
		currentBitSegment.SetBit(&currentBitSegment.Int, j, uint(value))
	}
	av.Segments[currentBitSegment.start] = currentBitSegment
}

func (av *Availability) getUnitInternal(from, to int) []byte {
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
