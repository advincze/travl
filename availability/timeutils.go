package availability

import (
	"time"
)

func TimeToUnit(t time.Time, res TimeResolution) int {
	return int(t.Unix() / int64(res))
}

func RoundDown(t time.Time, res TimeResolution) time.Time {
	if tooMuch := t.Unix() % int64(res); tooMuch != 0 {
		return t.Add(time.Duration(-1*tooMuch) * time.Second)
	}
	return t
}

func RoundUp(t time.Time, res TimeResolution) time.Time {
	if tooMuch := t.Unix() % int64(res); tooMuch != 0 {
		return t.Add(time.Duration(-1*tooMuch) * time.Second).Add(time.Duration(int(res)) * time.Second)
	}
	return t
}
