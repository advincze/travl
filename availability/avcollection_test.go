package availability

import (
	"testing"
)

func TestNewAvCollectionShouldNotReturnNil(t *testing.T) {
	avc := NewAvailabilityCollection()
	if avc == nil {
		t.Error("the availability collection should not be nil")
	}
}

func TestFindShouldReturnNilOnEmptyCollection(t *testing.T) {
	avc := NewAvailabilityCollection()
	av := avc.FindAvailabilityById("")
	if av != nil {
		t.Error("you should not find an availabilty for this id in an empty collection")
	}
}
