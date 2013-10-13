package availability

import (
	"time"
)

type Availability struct {
	internalRes TimeResolution
	data        *SegmentedVector
}

func NewAvailability(res TimeResolution) *Availability {
	return &Availability{
		internalRes: res,
		data:        NewSegmentedVector(int(Day / res)),
	}
}

func (av *Availability) Set(from, to time.Time, value byte) {
	fromUnit := TimeToUnit(from, av.internalRes)
	toUnit := TimeToUnit(to, av.internalRes)
	av.data.Set(fromUnit, toUnit, value)
}

func (av *Availability) SetAt(at time.Time, value byte) {
	atUnit := TimeToUnit(at, av.internalRes)
	av.data.Set(atUnit, atUnit+1, value)
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
	arr := av.data.Get(fromUnit, toUnit)
	factor := int(res / av.internalRes)
	reducedArr := reduceByFactor(arr, factor, reduceAllOne)
	return NewAvailabilityResult(res, av.internalRes, reducedArr, RoundDown(from, res))
}

func (av *Availability) getWithHigherResolution(from, to time.Time, res TimeResolution) *AvailabilityResult {
	// higher resolution
	fromUnitInternalRes := TimeToUnit(from, av.internalRes)
	toUnitInternalRes := TimeToUnit(RoundUp(to, av.internalRes), av.internalRes)
	arr := av.data.Get(fromUnitInternalRes, toUnitInternalRes)
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
	arr := av.data.Get(fromUnit, toUnit)
	return NewAvailabilityResult(res, av.internalRes, arr, RoundDown(from, res))
}

func (av *Availability) GetAt(at time.Time) byte {
	fromUnit := TimeToUnit(at, av.internalRes)
	toUnit := fromUnit + 1
	arr := av.data.Get(fromUnit, toUnit)
	return byte(arr[0])
}
