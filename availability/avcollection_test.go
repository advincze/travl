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
	t.Error("fail2")
}

func TestFindShouldNotRetrieveNilForSavedAV(t *testing.T) {
	avc := NewAvailabilityCollection()
	av := NewAvailability(Minute5)
	avc.SaveAvailability("0", av)
	avFound := avc.FindAvailabilityById("0")
	if avFound == nil {
		t.Error("avc should return the previously saved av")
	}
}

type TestSuite struct {
	param int
}

func (ts *TestSuite) TestWithParam(t *testing.T) {
	if ts.param%2 == 0 {
		t.Errorf("param was %d , fail", ts.param)
	}
}

func TestSuitefunc(t *testing.T) {
	testSuite := &TestSuite{param: 6}
	for i := 0; i < 10; i++ {
		TestNewAvCollectionShouldNotReturnNil(t)
		TestFindShouldReturnNilOnEmptyCollection(t)
		if i == 7 {
			t.Error("fail")
		}
		testSuite.param = i
		testSuite.TestWithParam(t)
	}

}
