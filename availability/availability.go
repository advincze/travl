package availability

import (
	"time"
)

type Availability interface {
	Set(from, to time.Time, value byte)
	Get(from, to time.Time, resolution TimeResolution) *BitVector
	SetAt(at time.Time, value byte)
	GetAt(at time.Time) byte
}

type AvailabilityCollection interface {
	FindOrCreate(id string) Availability
}
