package availability

type TimeResolution int32

const (
	sec      TimeResolution = 1
	Minute   TimeResolution = sec * 60
	Minute5  TimeResolution = Minute * 5
	Minute15 TimeResolution = Minute * 15
	Hour     TimeResolution = Minute * 60
	Day      TimeResolution = Hour * 24
)

var resolutions = map[TimeResolution]string{
	sec:      "sec",
	Minute:   "min",
	Minute5:  "5min",
	Minute15: "15min",
	Hour:     "hour",
	Day:      "day",
}

var resolutionStrings = map[string]TimeResolution{
	"s":     sec,
	"sec":   sec,
	"m":     Minute,
	"min":   Minute,
	"5m":    Minute5,
	"5min":  Minute5,
	"15m":   Minute15,
	"15min": Minute15,
	"h":     Hour,
	"hour":  Hour,
	"d":     Day,
	"day":   Day,
}

func (tr TimeResolution) String() string {
	return resolutions[tr]
}

func ParseTimeResolution(s string) TimeResolution {
	return resolutionStrings[s]
}
