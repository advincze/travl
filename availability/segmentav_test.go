package availability

import (
	"bytes"
	"testing"
	"time"
)

func getAvailability(res TimeResolution) Availability {
	//clear the DB
	return NewSegmentAv("testID", res)
}

func TestNewAvailabilityShouldNotBeNil(t *testing.T) {
	av := getAvailability(Minute5)

	if av == nil {
		t.Errorf("Availability should not be nil")
	}
}

func TestSetAvAtShouldNotPanic(t *testing.T) {
	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)

	defer func() {
		if r := recover(); r != nil {
			t.Error("SetAv should not have caused a panic")
		}
	}()

	av.SetAt(t1, 1)
}

func TestGetAvAtEmpty(t *testing.T) {
	// |0...000000000000000000000...000|
	//          |get
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	av := getAvailability(Minute5)

	//t
	if av.GetAt(t1) == 1 {
		t.Errorf("the bit should not be set")
	}
}

func TestGetAvAtSet(t *testing.T) {
	// |0...0001111111111111111100000...000|
	//          |-get
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	av := getAvailability(Minute5)
	av.SetAt(t1, 1)

	//w
	res := av.GetAt(t1)

	//t
	if res == 0 {
		t.Errorf("the bit should be set")
	}
}

func TestGetAvAtUnset(t *testing.T) {
	// |0...0001111111111111111100000...000|
	//          |-get
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	av := getAvailability(Minute5)
	av.SetAt(t1, 0)

	//w
	res := av.GetAt(t1)

	//t
	if res != 0 {
		t.Errorf("the bit should be unset")
	}
}

func TestSetAvFromToShouldNotPanic(t *testing.T) {
	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(24 * time.Hour)

	defer func() {
		if r := recover(); r != nil {
			t.Error("SetAv should not have caused a panic")
		}
	}()

	av.Set(t1, t2, 1)
}

func TestGetAvNothingFromEmpty(t *testing.T) {
	// |000000...000000000000|
	//       || get
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	av := getAvailability(Minute5)

	bitVector := av.Get(t1, t1, Minute5)

	if bitVector == nil {
		t.Errorf("the bitVector should not be nil")
	}
	if len(bitVector.Data) != 0 {
		t.Errorf("the bitVector should have length zero")
	}

}

func TestGetAvFromEmpty(t *testing.T) {
	// |000000...000000000000|
	//       |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(25 * time.Minute)
	av := getAvailability(Minute5)

	//w
	bitVector := av.Get(t1, t2, Minute5)

	//t
	if len(bitVector.Data) != 5 {
		t.Errorf("the bitVector should have length 5, was %d ", len(bitVector.Data))
	}
	if bitVector.Any() {
		t.Errorf("none of the bits should be set")
	}
}

func TestGetAvFromBeforeSet(t *testing.T) {
	// |000000...000000000011001....01100000|
	//     |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(25 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(75 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t3, t4, 1)

	//w
	bitVector := av.Get(t1, t2, Minute5)

	//t
	if l := len(bitVector.Data); l != 5 {
		t.Errorf("the bitVector bitset should have length 5, was %d", l)
	}
	if bitVector.Any() {
		t.Errorf("none of the bits should be set")
	}
}

func TestGetAvFromAfterSet(t *testing.T) {
	// |0...000111000110111111100000....00000|
	//                          |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t1, t2, 1)

	//w
	bitVector := av.Get(t3, t4, Minute5)

	//t
	if l := len(bitVector.Data); l != 2 {
		t.Errorf("the bitVector bitset should have length 2, was %d", l)
	}
	if bitVector.Any() {
		t.Errorf("none of the bits should be set")
	}
}

func TestGetAvFromInsideSet(t *testing.T) {
	// |0...0001111111111111111100000...000|
	//          |----get-----|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t1, t4, 1)

	//w
	bitVector := av.Get(t2, t3, Minute5)

	//t
	if l := len(bitVector.Data); l != 6 {
		t.Errorf("the bitVector should have length 6, was %v", l)
	}
	if !bitVector.All() {
		t.Errorf("all of the bits should be set")
	}

}

func TestGetAvFromItersectBeforeSet(t *testing.T) {
	// |00..000000001111111111100000...00|
	//         |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t2, t4, 1)

	//w
	bitVector := av.Get(t1, t3, Minute5)

	//t
	if l := len(bitVector.Data); l != 9 {
		t.Errorf("the bitVector should have length 9, was %d \n", l)
	}
	if c := bitVector.Count(); c != 6 {
		t.Errorf("6 of the bits should be set , %d were \n", c)
	}
}

func TestGetAvWithLowerResolution(t *testing.T) {
	// |00..000000001111111111100000...00|
	//         |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t2, t4, 1)

	//w
	bitVector := av.Get(t1, t3, Minute15)

	//t
	if l, exp := len(bitVector.Data), 3; l != exp {
		t.Errorf("the bitVector should have length %d, was %d \n", exp, l)
	}
	if c, exp := bitVector.Count(), 2; c != exp {
		t.Errorf("%d of the bits should be set , %d were \n", exp, c)
	}
}

func TestGetAvWithHigherResolution(t *testing.T) {
	// |00..000000001111111111100000...00|
	//         |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t2, t4, 1)

	//w
	bitVector := av.Get(t1, t3, Minute)

	//t
	if l, exp := len(bitVector.Data), 45; l != exp {
		t.Errorf("the bitVector should have length %d, was %d \n", exp, l)
	}
	if c, exp := bitVector.Count(), 30; c != exp {
		t.Errorf("%d of the bits should be set , %d were \n", exp, c)
	}
}

func TestGetAvFromItersectAfterSet(t *testing.T) {
	// |00..000000001111111111100000...00|
	//         |---get---|
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(15 * time.Minute)
	t3 := t1.Add(45 * time.Minute)
	t4 := t1.Add(55 * time.Minute)
	av := getAvailability(Minute5)
	av.Set(t1, t3, 1)

	//w
	bitVector := av.Get(t2, t4, Minute5)

	//t
	if l := len(bitVector.Data); l != 8 {
		t.Errorf("the bitVector should have length 8, was %d \n", l)
	}
	if c := bitVector.Count(); c != 6 {
		t.Errorf("6 of the bits should be set , %d were \n", c)
	}
}

func TestSetAvTwoYearsWorkingHoursShouldNotPanic(t *testing.T) {

	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 9, 0, 0, 0, time.UTC)

	defer func() {
		if r := recover(); r != nil {
			t.Error("SetAv should not have caused a panic")
		}
	}()

	for i := 0; i < 2*365; i++ {
		av.Set(t1, t1.Add(8*time.Hour), 1)
		av.Set(t1.Add(8*time.Hour), t1.Add(12*time.Hour), 0)
		t1 = t1.Add(24 * time.Hour)
	}
}

func TestSetAvTwoYearsWorkingHoursBackwardsShouldNotPanic(t *testing.T) {

	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 9, 0, 0, 0, time.UTC)

	defer func() {
		if r := recover(); r != nil {
			t.Error("SetAv should not have caused a panic")
		}
	}()

	for i := 0; i < 2*365; i++ {
		av.Set(t1, t1.Add(8*time.Hour), 1)
		av.Set(t1.Add(8*time.Hour), t1.Add(12*time.Hour), 0)
		t1 = t1.Add(-24 * time.Hour)
	}
}

func TestSetAvTwoYearsWorkingHours(t *testing.T) {

	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(1983, 4, 5, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		av.Set(t1, t1.Add(8*time.Hour), 1)
		av.Set(t1.Add(8*time.Hour), t1.Add(12*time.Hour), 0)
		t1 = t1.Add(24 * time.Hour)
	}

	bitVector := av.Get(t2, t2.Add(24*time.Hour), Hour)
	expected := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}
	if d := bitVector.Data; !bytes.Equal(expected, d) {
		t.Errorf("data should be equal to %v, was %v", expected, d)
	}

}

func TestSetAvTwoYearsWorkingHoursBackwards(t *testing.T) {

	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 9, 0, 0, 0, time.UTC)
	t2 := time.Date(1981, 4, 5, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		av.Set(t1, t1.Add(8*time.Hour), 1)
		av.Set(t1.Add(8*time.Hour), t1.Add(12*time.Hour), 0)
		t1 = t1.Add(-24 * time.Hour)
	}

	bitVector := av.Get(t2, t2.Add(24*time.Hour), Hour)
	expected := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}
	if d := bitVector.Data; !bytes.Equal(expected, d) {
		t.Errorf("data should be equal to %v, was %v", expected, d)
	}

}

func BenchmarkSetAvOneDay(b *testing.B) {

	av := getAvailability(Minute5)
	t1 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)

	for i := 0; i < b.N; i++ {
		av.Set(t1, t1.Add(24*time.Hour), 1)
		t1 = t1.Add(48 * time.Hour)
	}
}

func BenchmarkGetAvOneDay(b *testing.B) {

	av := getAvailability(Minute5)
	t0 := time.Date(1982, 2, 7, 0, 0, 0, 0, time.UTC)
	t1 := t0

	for i := 0; i < 2*365; i++ {
		av.Set(t1, t1.Add(8*time.Hour), 1)
		av.Set(t1.Add(8*time.Hour), t1.Add(12*time.Hour), 0)
		t1 = t1.Add(24 * time.Hour)
	}

	t1 = t0.Add(4 * time.Hour)
	for i := 0; i < b.N; i++ {
		av.Get(t1, t1.Add(72*time.Hour), Hour)

		t1 = t1.Add(23 * time.Hour)
		if i%365 == 0 {
			t1 = t0
		}
	}
}
