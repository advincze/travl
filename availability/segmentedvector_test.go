package availability

import (
	"testing"
)

func TestNewVectorShouldNotBeNil(t *testing.T) {
	vector := NewSegmentedVector(7)
	if vector == nil {
		t.Error("a new segmented vector hould not be nil")
	}
}
