package availability

import (
	"bytes"
	"strconv"
	"time"
)

type Availability struct {
	id            string
	internalRes   TimeResolution
	segmentLength int
	segments      map[int]*BitSegment
}

func NewAvailability(id string, res TimeResolution) *Availability {
	return &Availability{
		id:            id,
		internalRes:   res,
		segmentLength: int(Day / res),
		segments:      make(map[int]*BitSegment),
	}
}

func (av *Availability) sizeInBytes() int {
	var sizeInBytes int
	for _, segment := range av.segments {
		sizeInBytes += len(segment.Bytes())
	}
	return sizeInBytes
}

func (av *Availability) Set(from, to time.Time, value byte) {
	fromUnit := TimeToUnit(from, av.internalRes)
	toUnit := TimeToUnit(to, av.internalRes)
	av.setUnit(fromUnit, toUnit, value)
}

func (av *Availability) SetAt(at time.Time, value byte) {
	atUnit := TimeToUnit(at, av.internalRes)
	av.setUnit(atUnit, atUnit+1, value)
}

func (av *Availability) setUnit(from, to int, value byte) {
	segmentStart := av.segmentStart(from)
	segment := av.getOrEmptyBitSegment(segmentStart)

	for i, j := from, from%av.segmentLength; i < to; i, j = i+1, j+1 {
		if j == av.segmentLength {
			av.storeSegment(segment)
			segment = av.getOrEmptyBitSegment(i)
			j = 0
		}
		segment.SetUnit(j, value)
	}
	av.segments[segment.start] = segment
}

func (av *Availability) storeSegment(segment *BitSegment) {
	av.segments[segment.start] = segment
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
	fromUnit := TimeToUnit(RoundDown(from, res), av.internalRes)
	toUnit := TimeToUnit(RoundUp(to, res), av.internalRes)
	arr := av.getUnit(fromUnit, toUnit)
	factor := int(res / av.internalRes)
	reducedArr := reduceByFactor(arr, factor, reduceAllOne)
	return NewAvailabilityResult(res, av.internalRes, reducedArr, RoundDown(from, res))
}

func (av *Availability) getWithHigherResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	// higher resolution
	fromUnitInternalRes := TimeToUnit(from, av.internalRes)
	toUnitInternalRes := TimeToUnit(RoundUp(to, av.internalRes), av.internalRes)
	arr := av.getUnit(fromUnitInternalRes, toUnitInternalRes)
	factor := int(av.internalRes / res)
	arrMultiplied := multiplyByFactor(arr, factor)
	cutoff := TimeToUnit(from, res) - fromUnitInternalRes*factor
	origlen := TimeToUnit(to, res) - TimeToUnit(from, res)
	arrTrimmed := arrMultiplied[cutoff : cutoff+origlen]
	return NewAvailabilityResult(res, av.internalRes, arrTrimmed, RoundDown(from, av.internalRes))
}

func (av *Availability) getWithInternalResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	fromUnit := TimeToUnit(from, res)
	toUnit := TimeToUnit(to, res)
	arr := av.getUnit(fromUnit, toUnit)
	return NewAvailabilityResult(res, av.internalRes, arr, RoundDown(from, res))
}

func (av *Availability) GetAt(at time.Time) byte {
	fromUnit := TimeToUnit(at, av.internalRes)
	toUnit := fromUnit + 1
	arr := av.getUnit(fromUnit, toUnit)
	return byte(arr[0])
}

func (av *Availability) getUnit(from, to int) []byte {
	length := to - from
	result := make([]byte, length)
	currentBitSegment := av.getOrEmptyBitSegment(av.segmentStart(from))
	for i, j := 0, from%av.segmentLength; i < length; i, j = i+1, j+1 {
		if j == av.segmentLength {
			currentBitSegment = av.getOrEmptyBitSegment(i + from)
			j = 0
		}
		result[i] = byte(currentBitSegment.Bit(j))
	}
	return result
}

func (av *Availability) getOrEmptyBitSegment(startValue int) *BitSegment {
	if segment := av.segments[startValue]; segment != nil {
		return segment
	}
	return NewBitSegment(startValue)
}

func (av *Availability) segmentStart(i int) int {
	return i - i%av.segmentLength
}

func (av *Availability) String() string {
	var buffer bytes.Buffer
	for _, segment := range av.segments {
		buffer.WriteString(strconv.Itoa(segment.start))
		buffer.WriteString("->")
		buffer.WriteString(segment.String())
		buffer.WriteRune('\n')
	}
	return buffer.String()
}
